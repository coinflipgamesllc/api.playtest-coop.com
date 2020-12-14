package app

import (
	"github.com/gin-gonic/gin"
)

func (s *Server) routes() {
	v1 := s.router.Group("/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.GET("/user", s.authenticated, s.handleGetUser())
			auth.PUT("/user", s.authenticated, s.handleUpdateUser())

			auth.POST("/signup", s.handleSignup())
			auth.POST("/login", s.handleLogin())
			auth.POST("/token", s.handleRefreshToken())
			auth.GET("/verify-email/:id", s.handleVerifyEmail())
		}

		games := v1.Group("/games")
		{
			games.GET("", s.handleListGames())
			games.POST("", s.authenticated, s.handleCreateGame())
			games.GET("/:id", s.handleGetGame())
			games.PUT("/:id", s.authenticated, s.handleUpdateGame())
		}
	}
}

func serverError(err error) gin.H {
	return gin.H{"error": err.Error()}
}
