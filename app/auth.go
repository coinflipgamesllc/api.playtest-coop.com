package app

import (
	"errors"
	"time"

	"github.com/coinflipgamesllc/api.playtest-coop.com/domain"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
)

func (s *Server) createTokensForUser(user *domain.User) (string, string, error) {
	accessToken := jwt.New(jwt.GetSigningMethod("HS256"))
	accessToken.Claims = jwt.MapClaims{
		"sub":  user.ID,
		"name": user.Name,
		"exp":  time.Now().Add(time.Minute * 15).Unix(),
	}

	at, err := accessToken.SignedString([]byte(s.authToken))
	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.New(jwt.GetSigningMethod("HS256"))
	refreshToken.Claims = jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	rt, err := refreshToken.SignedString([]byte(s.authToken))
	if err != nil {
		return "", "", err
	}

	return at, rt, nil
}

func (s *Server) handleGetUser() gin.HandlerFunc {
	type response struct {
		User *domain.User `json:"user"`
	}

	return func(c *gin.Context) {
		user := s.user(c)

		c.JSON(200, response{User: user})
	}
}

func (s *Server) handleUpdateUser() gin.HandlerFunc {
	type request struct {
		Name        string `json:"name" binding:"omitempty,min=2"`
		Email       string `json:"email" binding:"omitempty,email"`
		NewPassword string `json:"new_password" binding:"omitempty,nefield=OldPassword,min=10"`
		OldPassword string `json:"old_password" binding:"omitempty"`
		Pronouns    string `json:"pronouns" binding:"omitempty,contains=/"`
	}

	type response struct {
		User *domain.User `json:"user"`
	}

	return func(c *gin.Context) {
		user := s.user(c)

		// Validate request
		var req request
		if err := c.ShouldBind(&req); err != nil {
			c.AbortWithStatusJSON(400, serverError(err))
			return
		}

		// Update the user
		if req.Name != "" {
			user.Rename(req.Name)
		}

		if req.Email != "" {
			user.ChangeEmail(req.Email)
		}

		if req.NewPassword != "" && req.OldPassword != "" {
			err := user.ChangePassword(req.NewPassword, req.OldPassword)
			if err != nil {
				c.AbortWithStatusJSON(500, serverError(err))
				return
			}
		}

		if req.Pronouns != "" {
			user.SetPronouns(req.Pronouns)
		}

		// Save changes
		err := s.userRepository.Save(user)
		if err != nil {
			c.AbortWithStatusJSON(500, serverError(err))
			return
		}

		c.JSON(200, response{User: user})
	}
}

func (s *Server) handleSignup() gin.HandlerFunc {
	type request struct {
		Name     string `json:"name" binding:"required,min=2"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=10"`
	}

	type response struct {
		User         *domain.User `json:"user"`
		AccessToken  string       `json:"access_token"`
		RefreshToken string       `json:"refresh_token"`
	}

	return func(c *gin.Context) {
		// Validate request
		var req request
		if err := c.ShouldBind(&req); err != nil {
			c.AbortWithStatusJSON(400, serverError(err))
			return
		}

		// Create user
		user, err := domain.NewUser(req.Name, req.Email, req.Password)
		if err != nil {
			c.AbortWithStatusJSON(500, serverError(err))
			return
		}

		// Save user
		err = s.userRepository.Save(user)
		if err != nil {
			c.AbortWithStatusJSON(500, serverError(err))
			return
		}

		// Generate tokens for future requests
		at, rt, err := s.createTokensForUser(user)
		if err != nil {
			c.AbortWithStatusJSON(500, serverError(err))
			return
		}

		c.JSON(200, response{User: user, AccessToken: at, RefreshToken: rt})
	}
}

func (s *Server) handleLogin() gin.HandlerFunc {
	type request struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=10"`
	}

	type response struct {
		User         *domain.User `json:"user"`
		AccessToken  string       `json:"access_token"`
		RefreshToken string       `json:"refresh_token"`
	}

	return func(c *gin.Context) {
		// Validate request
		var req request
		if err := c.ShouldBind(&req); err != nil {
			c.AbortWithStatusJSON(400, serverError(err))
			return
		}

		// Retrieve user
		user, err := s.userRepository.UserOfEmail(req.Email)
		if err != nil {
			c.AbortWithStatusJSON(500, serverError(err))
			return
		}

		if user == nil {
			c.AbortWithStatusJSON(404, serverError(errors.New("no account found with that email and password")))
			return
		}

		// Verify password
		if ok, err := user.ValidPassword(req.Password); !ok || err != nil {
			c.AbortWithStatusJSON(401, serverError(errors.New("unauthorized")))
			return
		}

		// Generate tokens for future requests
		at, rt, err := s.createTokensForUser(user)
		if err != nil {
			c.AbortWithStatusJSON(500, serverError(err))
			return
		}

		c.JSON(200, response{User: user, AccessToken: at, RefreshToken: rt})
	}
}

func (s *Server) handleRefreshToken() gin.HandlerFunc {
	type request struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	type response struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	return func(c *gin.Context) {
		// Validate request
		var req request
		if err := c.ShouldBind(&req); err != nil {
			c.AbortWithStatusJSON(400, serverError(err))
			return
		}

		// Validate token
		token, err := jwt.Parse(req.RefreshToken, func(token *jwt.Token) (interface{}, error) {
			return []byte(s.authToken), nil
		})

		if err != nil {
			c.AbortWithStatusJSON(400, serverError(err))
			return
		}

		// Extract and validate that the user account is still active
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			id := claims["sub"]

			user, err := s.userRepository.UserOfID(uint(id.(float64)))
			if err != nil {
				c.AbortWithStatusJSON(400, serverError(err))
				return
			}

			if user == nil {
				c.AbortWithStatusJSON(401, serverError(errors.New("unauthorized")))
				return
			}

			// Generate a new token pair
			at, rt, err := s.createTokensForUser(user)
			if err != nil {
				c.AbortWithStatusJSON(500, serverError(err))
				return
			}

			c.JSON(200, response{AccessToken: at, RefreshToken: rt})
			return
		}

		c.JSON(401, serverError(errors.New("unauthorized")))
	}
}

func (s *Server) handleVerifyEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Fetch user by ID
		id := c.Param("id")
		user, err := s.userRepository.UserOfVerificationID(id)
		if err != nil {
			c.HTML(500, "500.html", gin.H{"error": err.Error()})
			return
		}

		if user == nil {
			// User can't be found so either they don't exist or the validation ID was already used.
			// Just redirect to home and hope for the best.
			c.Redirect(307, "https://playtest-coop.com")
			return
		}

		// Mark verified and save
		user.VerifyEmail()
		if err := s.userRepository.Save(user); err != nil {
			c.HTML(500, "500.html", gin.H{"error": err.Error()})
			return
		}

		// Send em home
		c.Redirect(307, "https://playtest-coop.com")
	}
}

func (s *Server) authenticated(c *gin.Context) {
	token, err := request.ParseFromRequest(c.Request, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.authToken), nil
	})

	if err != nil {
		c.AbortWithError(401, err)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		c.Set("user_id", claims["sub"])
	} else {
		c.AbortWithStatus(401)
	}
}

func (s *Server) user(c *gin.Context) *domain.User {
	// Retrieve the user ID from the context
	id, ok := c.Get("user_id")
	if !ok {
		c.AbortWithStatusJSON(401, serverError(errors.New("unauthorized")))
		return nil
	}

	// Fetch the user
	user, err := s.userRepository.UserOfID(uint(id.(float64)))
	if err != nil {
		c.AbortWithStatusJSON(500, serverError(err))
		return nil
	}

	return user
}
