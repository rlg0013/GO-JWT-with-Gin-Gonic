package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/rlg0013/GO-JWT-with-Gin-Gonic/controllers"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("users/signup", controller.Signup())
	incomingRoutes.POST("users/login", controller.Login())
}
