package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
)

// Authenticated checks for a valid access token in the request
func Authenticated(authToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := request.ParseFromRequest(c.Request, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
			return []byte(authToken), nil
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
}
