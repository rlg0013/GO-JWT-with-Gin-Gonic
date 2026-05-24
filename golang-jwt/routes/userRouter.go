package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/rlg0013/GO-JWT-with-Gin-Gonic/controllers"
	"github.com/rlg0013/GO-JWT-with-Gin-Gonic/middleware"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/users", controller.GetUsers())
	incomingRoutes.GET("/users/:user_id", controller.GetUsers())

}
