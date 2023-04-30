package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/wahri/Technical_Test_MNC_2/db"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

const (
	tokenDuration = time.Hour * 24
	tokenIssuer   = "your-issuer"
	tokenSecret   = "your-secret"
)

type AuthController struct{}

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	PIN         string `json:"pin"`
}

func (a *AuthController) LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	// Find the user with the given phone number
	client, _ := db.ConnectMongoDB()
	collection := client.Database("example_data").Collection("users")
	filter := bson.M{"phone_number": req.PhoneNumber}
	var user User
	if err := collection.FindOne(context.Background(), filter).Decode(&user); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid credentials",
		})
		return
	}

	// Check if the PIN is correct
	if err := bcrypt.CompareHashAndPassword([]byte(user.PIN), []byte(req.PIN)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid credentials",
		})
		return
	}

	// Generate a JWT token for the user
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"phoneNumber": user.PhoneNumber,
		"firstName":   user.FirstName,
		"lastName":    user.LastName,
		"exp":         time.Now().Add(tokenDuration).Unix(),
		"iat":         time.Now().Unix(),
		"iss":         tokenIssuer,
	})
	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": signedToken,
	})
}
