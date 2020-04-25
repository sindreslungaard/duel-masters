package server

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sindreslungaard/duel-masters/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// GetUserForToken returns a user from the authorization header or returns an error
func GetUserForToken(c *gin.Context) (db.User, error) {

	collection := db.Collection("users")

	var user db.User

	if err := collection.FindOne(context.TODO(), bson.M{"sessions": bson.M{"$elemMatch": bson.M{"token": c.GetHeader("Authorization")}}}).Decode(&user); err != nil {
		return db.User{}, err
	}

	return user, nil

}

type signinReqBody struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// SigninHandler handles signin requests
func SigninHandler(c *gin.Context) {

	var reqBody signinReqBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.Status(400)
		return
	}

	collection := db.Collection("users")

	var user db.User

	if err := collection.FindOne(context.TODO(), bson.M{"username": primitive.Regex{Pattern: "^" + reqBody.Username + "$", Options: "i"}}).Decode(&user); err != nil {
		c.Status(404)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqBody.Password)); err != nil {
		c.Status(401)
		return
	}

	token, err := uuid.NewRandom()
	if err != nil {
		c.Status(500)
		return
	}

	session := db.UserSession{
		Token:   token.String(),
		IP:      c.ClientIP(),
		Expires: int(time.Now().Add(time.Second * 2592000).Unix()),
	}

	collection.UpdateOne(context.TODO(), bson.M{"uid": user.UID}, bson.M{"$push": bson.M{"sessions": session}})

	c.JSON(200, bson.M{"user": user, "token": session.Token})

	// TODO: Remove expired/unneeded sessions from db

}

type signupReqBody struct {
	Username string `json:"username" binding:"required,alphanum,min=3,max=20"`
	Password string `json:"password" binding:"required,min=6,max=255"`
	Email    string `json:"email" binding:"required,email"`
}

// SignupHandler handles signup requests
func SignupHandler(c *gin.Context) {

	// TODO: recaptcha

	var reqBody signupReqBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.Status(400)
		return
	}

	collection := db.Collection("users")

	if err := collection.FindOne(context.TODO(), bson.M{"username": primitive.Regex{Pattern: "^" + reqBody.Username + "$", Options: "i"}}).Decode(&db.User{}); err == nil {
		c.JSON(400, bson.M{"message": "The username is already taken"})
		return
	}

	if err := collection.FindOne(context.TODO(), bson.M{"email": primitive.Regex{Pattern: "^" + reqBody.Email + "$", Options: "i"}}).Decode(&db.User{}); err == nil {
		c.JSON(400, bson.M{"message": "The email is already taken"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), 10)

	if err != nil {
		c.Status(500)
		return
	}

	token, err := uuid.NewRandom()
	if err != nil {
		c.Status(500)
		return
	}

	session := db.UserSession{
		Token:   token.String(),
		IP:      c.ClientIP(),
		Expires: int(time.Now().Add(time.Second * 2592000).Unix()),
	}

	user := db.User{
		UID:         uuid.New().String(),
		Username:    reqBody.Username,
		Email:       reqBody.Email,
		Password:    string(hash),
		Permissions: []string{},
		Sessions: []db.UserSession{
			session,
		},
	}

	_, err = collection.InsertOne(context.TODO(), user)

	if err != nil {
		c.Status(500)
		return
	}

	c.JSON(200, bson.M{"user": user, "token": session.Token})

}
