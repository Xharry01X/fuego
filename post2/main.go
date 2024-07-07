package main

import (
	"context"
	"time"

	"github.com/go-fuego/fuego"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
)

type UserInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserOutput struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

var validate *validator.Validate
var userCollection *mongo.Collection
var jwtKey = []byte("your_secret_key")

func main() {
	validate = validator.New()

	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}

	userCollection = client.Database("Fuego").Collection("users")

	s := fuego.NewServer(
		fuego.WithAddr(":4000"),
	)

	fuego.Post(s,"/user",createUserHandler)
	fuego.Get(s,"/",getAllUser)

	s.Run()
}

func createUserHandler(c *fuego.ContextWithBody[UserInput]) (UserOutput, error) {
	body, err := c.Body()
	if err != nil {
		return UserOutput{}, err
	}

	// Validate the body
	err = validate.Struct(body)
	if err != nil {
		return UserOutput{}, err
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return UserOutput{}, err
	}

	// Save user to MongoDB
	user := bson.D{
		{Key: "username", Value: body.Username},
		{Key: "password", Value: string(hashedPassword)},
	}

	_, err = userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		return UserOutput{}, err
	}

	// Create JWT token
	token, err := generateJWT(body.Username)
	if err != nil {
		return UserOutput{}, err
	}

	return UserOutput{
		Username: body.Username,
		Token:    token,
	}, nil
}
func generateJWT(username string) (string, error) {
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		Issuer:    username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}


func getAllUser(c *fuego.ContextNoBody) ([]UserOutput, error) {
	var users []UserOutput
	cursor, err := userCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var user bson.D
		err := cursor.Decode(&user)
		if err != nil {
			return nil, err
		}
		username, _ := user.Map()["username"].(string)
		token, err := generateJWT(username)
		if err != nil {
			return nil, err
		}
		users = append(users, UserOutput{Username: username, Token: token})
	}
	return users, nil
}

