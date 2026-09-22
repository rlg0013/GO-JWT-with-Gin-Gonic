package helper

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func MatchUserTypeToUid(c *gin.Context, userID string) error {
	if c.GetString("user_type") == "ADMIN" || c.GetString("uid") == userID {
		return nil
	}
	return fmt.Errorf("you are not authorized to access this user")
}
