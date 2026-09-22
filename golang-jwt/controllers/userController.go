package controller

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rlg0013/GO-JWT-with-Gin-Gonic/database"
	helper "github.com/rlg0013/GO-JWT-with-Gin-Gonic/helpers"
	"github.com/rlg0013/GO-JWT-with-Gin-Gonic/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

func users() *mongo.Collection { return database.OpenCollection(database.Client, "users") }
func value(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
func stringPtr(s string) *string { return &s }
func publicUser(u models.User) gin.H {
	return gin.H{"id": u.User_id, "first_name": value(u.First_Name), "last_name": value(u.Last_Name), "email": value(u.Email), "phone": value(u.Phone), "user_type": value(u.User_Type), "token": value(u.Token), "refresh_token": value(u.Refresh_Token), "created_at": u.Created_at, "updated_at": u.Updated_at}
}

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var u models.User
		if err := c.ShouldBindJSON(&u); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := validate.Struct(u); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		u.Email = stringPtr(strings.ToLower(strings.TrimSpace(*u.Email)))
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := users().FindOne(ctx, bson.M{"email": *u.Email}).Err(); err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
			return
		} else if err != mongo.ErrNoDocuments {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(*u.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		password := string(hash)
		u.Password = &password
		u.ID = primitive.NewObjectID()
		u.User_id = u.ID.Hex()
		now := time.Now().UTC()
		u.Created_at, u.Updated_at = now, now
		token, refresh, err := helper.GenerateAllTokens(*u.Email, *u.First_Name, value(u.Last_Name), u.User_id, *u.User_Type)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		u.Token, u.Refresh_Token = &token, &refresh
		if _, err = users().InsertOne(ctx, u); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, publicUser(u))
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var in struct {
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required"`
		}
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var u models.User
		if err := users().FindOne(ctx, bson.M{"email": strings.ToLower(strings.TrimSpace(in.Email))}).Decode(&u); err != nil || u.Password == nil || bcrypt.CompareHashAndPassword([]byte(*u.Password), []byte(in.Password)) != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			return
		}
		token, refresh, err := helper.GenerateAllTokens(*u.Email, *u.First_Name, value(u.Last_Name), u.User_id, *u.User_Type)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		u.Token, u.Refresh_Token = &token, &refresh
		u.Updated_at = time.Now().UTC()
		if _, err = users().UpdateByID(ctx, u.ID, bson.M{"$set": bson.M{"token": token, "refresh_token": refresh, "updated_at": u.Updated_at}}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, publicUser(u))
	}
}

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetString("user_type") != "ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		cursor, err := users().Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cursor.Close(ctx)
		var result []models.User
		if err = cursor.All(ctx, &result); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		response := make([]gin.H, 0, len(result))
		for _, user := range result {
			response = append(response, publicUser(user))
		}
		c.JSON(http.StatusOK, response)
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("user_id")
		if err := helper.MatchUserTypeToUid(c, id); err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var u models.User
		if err := users().FindOne(ctx, bson.M{"user_id": id}).Decode(&u); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusOK, publicUser(u))
	}
}
