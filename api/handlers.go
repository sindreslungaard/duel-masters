package api

import (
	"context"
	"net/http"
	"time"

	"duel-masters/db"
	"duel-masters/game"
	"duel-masters/game/match"
	"duel-masters/server"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
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
		c.Status(401)
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

type matchReqBody struct {
	Name string `json:"name" binding:"required,min=3,max=100"`
}

// MatchHandler handles creation of new mathes
func MatchHandler(c *gin.Context) {

	user, err := db.GetUserForToken(c.GetHeader("Authorization"))
	if err != nil {
		c.Status(401)
		return
	}

	var reqBody matchReqBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.Status(400)
		return
	}

	m := match.New(reqBody.Name, user.UID)

	c.JSON(200, m)

}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// WS handles websocket upgrade
func WS(c *gin.Context) {

	hubID := c.Param("hub")

	var hub server.Hub

	if hubID == "lobby" {

		hub = game.GetLobby()

	} else {

		m, err := match.Find(hubID)

		if err != nil {
			c.Status(404)
			return
		}

		hub = m

	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		c.Status(500)
		return
	}

	s := server.NewSocket(conn, hub)

	// Handle the connection in a new goroutine to free up this memory
	go s.Listen()

}

// CardsHandler returns a list of all the cards in the cache
func CardsHandler(c *gin.Context) {
	c.JSON(200, GetCache())
}

// GetDecksHandler returns an array of the users decks
func GetDecksHandler(c *gin.Context) {

	user, err := db.GetUserForToken(c.GetHeader("Authorization"))
	if err != nil {
		c.Status(401)
		return
	}

	collection := db.Collection("decks")

	cur, err := collection.Find(context.TODO(), bson.M{
		"owner": user.UID,
	})

	if err != nil {
		logrus.Error(err)
		c.Status(500)
		return
	}

	defer cur.Close(context.TODO())

	decks := make([]db.Deck, 0)

	for cur.Next(context.TODO()) {

		var deck db.Deck

		if err := cur.Decode(&deck); err != nil {
			continue
		}

		decks = append(decks, deck)

	}

	c.JSON(200, decks)

}

type createDeckBody struct {
	Name   string   `json:"name" binding:"required,min=1,max=30"`
	Cards  []string `json:"cards" binding:"required"`
	UID    string   `json:"uid"`
	Public bool     `json:"public"`
}

// CreateDeckHandler handles creating/editing decks
func CreateDeckHandler(c *gin.Context) {

	user, err := db.GetUserForToken(c.GetHeader("Authorization"))
	if err != nil {
		c.Status(401)
		return
	}

	var reqBody createDeckBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.Status(400)
		return
	}

	if len(reqBody.Cards) < 40 || len(reqBody.Cards) > 50 {
		c.Status(400)
		return
	}

	for _, cuid := range reqBody.Cards {
		if !CacheHas(cuid) {
			c.Status(400)
			return
		}
	}

	collection := db.Collection("decks")

	if len(reqBody.UID) < 1 {

		// New deck

		decksCount, err := collection.CountDocuments(context.TODO(), bson.M{"owner": user.UID})

		if err != nil {
			logrus.Error(err)
			c.Status(500)
			return
		}

		if decksCount >= 15 {
			c.Status(403)
			return
		}

		deck := db.Deck{
			UID:      uuid.New().String(),
			Owner:    user.UID,
			Name:     reqBody.Name,
			Public:   reqBody.Public,
			Standard: false,
			Cards:    reqBody.Cards,
		}

		_, err = collection.InsertOne(context.TODO(), deck)

		if err != nil {
			c.Status(500)
			return
		}

	} else {

		// Edit deck

		_, err := collection.UpdateOne(
			context.TODO(),
			bson.M{"uid": reqBody.UID, "owner": user.UID},
			bson.M{"$set": bson.M{"name": reqBody.Name, "public": reqBody.Public, "cards": reqBody.Cards}},
		)

		if err != nil {
			logrus.Error(err)
			c.Status(500)
			return
		}

	}

	c.Status(200)

}
