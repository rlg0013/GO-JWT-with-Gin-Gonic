package middleware

import (
	"github.com/gin-gonic/gin"
	helper "github.com/rlg0013/GO-JWT-with-Gin-Gonic/helpers"
	"net/http"
	"strings"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		parts := strings.Fields(c.GetHeader("Authorization"))
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header must be Bearer <token>"})
			return
		}
		claims, err := helper.ValidateToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}
		c.Set("uid", claims.UID)
		c.Set("user_type", claims.UserType)
		c.Next()
	}
}
