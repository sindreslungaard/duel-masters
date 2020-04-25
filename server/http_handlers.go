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
