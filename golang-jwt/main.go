package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rlg0013/GO-JWT-with-Gin-Gonic/database"
	routes "github.com/rlg0013/GO-JWT-with-Gin-Gonic/routes"
)

func main() {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Fatalf("load .env: %v", err)
	}
	if err := database.Connect(os.Getenv("MONGODB_URL")); err != nil {
		log.Fatalf("connect to MongoDB: %v", err)
	}
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())

	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	router.GET("/api-1", func(c *gin.Context) {
		c.JSON(200, gin.H{"succeess": "access granted for api-1"})
	})

	router.GET("/api-2", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "access granted for api-2"})
	})

	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}

}
