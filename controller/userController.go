package controller

import (
	"context"
	"fmt"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wahri/Technical_Test_MNC_2/db"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

type UserController struct{}

type User struct {
	UserID      uuid.UUID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	FirstName   string    `json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName    string    `json:"last_name,omitempty" bson:"last_name,omitempty"`
	PhoneNumber string    `json:"phone_number,omitempty" bson:"phone_number,omitempty"`
	Address     string    `json:"address,omitempty" bson:"address,omitempty"`
	PIN         string    `json:"pin,omitempty" bson:"pin,omitempty"`
	Balance     string    `json:"balance,omitempty" bson:"balance,omitempty"`
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func (u *UserController) RegisterHandler(c *gin.Context) {
	var request User
	if err := c.ShouldBindJSON(&request); err != nil {
		data := gin.H{
			"message": "field required",
		}
		c.JSON(420, data)
		return
	}

	client, _ := db.ConnectMongoDB()
	collection := client.Database("example_data").Collection("users")

	filter := bson.M{"phone_number": request.PhoneNumber}
	var existingUser User
	err := collection.FindOne(context.Background(), filter).Decode(&existingUser)
	if err == nil {
		data := gin.H{
			"message": "Phone number already taken",
		}
		c.JSON(400, data)
		return
	}

	hashedPin, err := Hash(request.PIN)
	if err != nil {
		data := gin.H{
			"message": "Internal Server Error",
		}
		c.JSON(500, data)
		return
	}
	request.PIN = string(hashedPin)

	request.UserID = uuid.New()
	request.Balance = "0"

	fmt.Println(request.Balance)
	// Insert new user into database
	collection.InsertOne(context.Background(), request)

	c.JSON(200, gin.H{
		"status": "SUCCESS",
		"result": request,
	})
}

type TopUpController struct{}

type Claims struct {
	PhoneNumber string `json:"phoneNumber"`
	jwt.StandardClaims
}

type TopUpRequest struct {
	Amount string `json:"amount"`
}

func (t *TopUpController) TopUpHandler(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(401, gin.H{"message": "Authorization header required"})
		return
	}

	tokenString := authHeader[len("Bearer "):]
	claims := &Claims{}

	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("your-secret"), nil // TODO: Change this to your own secret key
	})

	if err != nil {
		c.JSON(401, gin.H{"message": "Invalid token"})
		return
	}

	phoneNumber := claims.PhoneNumber

	client, _ := db.ConnectMongoDB()
	collection := client.Database("example_data").Collection("users")

	var user User
	var request TopUpRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		data := gin.H{
			"message": "field required",
		}
		c.JSON(420, data)
		return
	}

	err = collection.FindOne(context.Background(), bson.M{"phone_number": phoneNumber}).Decode(&user)
	if err != nil {
		c.JSON(404, gin.H{"message": "User not found"})
		return
	}

	amount, err := strconv.ParseFloat(request.Amount, 64)
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid amount"})
		return
	}

	if amount < 0 {
		c.JSON(400, gin.H{"message": "Amount must be positive"})
		return
	}

	// user.Balance = string(float64(user.Balance) + amount)
	currentBalance, _ := strconv.ParseFloat(user.Balance, 64)
	newBalance := currentBalance + amount

	_, err = collection.UpdateOne(context.Background(), bson.M{"phone_number": phoneNumber}, bson.M{"$set": bson.M{"balance": strconv.FormatFloat(newBalance, 'f', 0, 64)}})
	if err != nil {
		c.JSON(500, gin.H{"message": "Internal server error"})
		return
	}

	c.JSON(200, gin.H{
		"status":  "SUCCESS",
		"balance": newBalance,
	})
}
