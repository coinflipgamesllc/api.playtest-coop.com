package middleware

import (
	"github.com/coinflipgamesllc/api.playtest-coop.com/domain"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Authenticated checks for a valid access token in the request
func Authenticated(authToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")

		if userID == nil {
			c.AbortWithStatusJSON(401, gin.H{"error": domain.Unauthorized{}.Error()})
			return
		}

		c.Next()
	}
}
