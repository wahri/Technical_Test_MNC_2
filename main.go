package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/wahri/Technical_Test_MNC_2/controller"
	"github.com/wahri/Technical_Test_MNC_2/db"
	"golang.org/x/crypto/bcrypt"
)

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func main() {
	// Connect to MongoDB
	client, err := db.ConnectMongoDB()
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(context.Background())

	router := gin.Default()

	UserController := new(controller.UserController)
	router.POST("/register", UserController.RegisterHandler)

	AuthController := new(controller.AuthController)
	router.POST("/login", AuthController.LoginHandler)

	TopUpController := new(controller.TopUpController)
	router.POST("/topup", TopUpController.TopUpHandler)

	// Start the server and listen for incoming requests
	router.Run(":8080")
}
