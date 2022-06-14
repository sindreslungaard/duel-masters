package api

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"duel-masters/db"
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
func (api *API) SigninHandler(c *gin.Context) {

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

	// Check for IP ban
	bansCollection := db.Collection("bans")
	ip := c.ClientIP()

	bans, err := bansCollection.CountDocuments(context.Background(), bson.M{
		"$or": []bson.M{
			{"type": db.UserBan, "value": user.UID},
			{"type": db.IPBan, "value": ip},
		},
	})

	if err != nil {
		logrus.Error(err)
		c.JSON(500, bson.M{"message": "An internal error occured"})
		return
	}

	if bans > 0 {
		c.JSON(403, bson.M{"message": "You have been banned"})
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
func (api *API) SignupHandler(c *gin.Context) {

	// TODO: recaptcha

	var reqBody signupReqBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(400, bson.M{"message": "Make sure your username only consist of a-Z and 0-9 (3-20 characters long). Password minimum 6 characters."})
		return
	}

	// Check for IP ban
	collection := db.Collection("bans")
	ip := c.ClientIP()

	bans, err := collection.CountDocuments(context.Background(), bson.M{"type": db.IPBan, "value": ip})

	if err != nil {
		logrus.Error(err)
		c.JSON(500, bson.M{"message": "An internal error occured"})
		return
	}

	if bans > 0 {
		c.JSON(403, bson.M{"message": "You have been banned"})
		return
	}

	// Check if IP is in list of blocked networks
	if api.blockedNetworks.Contains(ip) {
		c.JSON(403, bson.M{"message": "Your network has been blocked from creating new accounts. If you believe this is wrong, please contact us on discord: https://discord.gg/FkPTE4p"})
		return
	}

	collection = db.Collection("users")

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
	Name       string `json:"name" binding:"max=50"`
	Visibility string `json:"visibility" binding:"required"`
}

var defaultMatchNames = []string{
	"Kettou Da!",
	"I challenge you!",
	"Ikuzo!",
	"I'm ready!",
	"Koi!",
	"Bring it on!",
}

// MatchHandler handles creation of new mathes
func (api *API) MatchHandler(c *gin.Context) {

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

	visible := true
	if reqBody.Visibility == "private" {
		visible = false
	}

	name := reqBody.Name

	if name == "" {
		name = defaultMatchNames[rand.Intn(len(defaultMatchNames))]
	}

	m := api.matchSystem.NewMatch(name, user.UID, visible)

	c.JSON(200, m)

}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// WS handles websocket upgrade
func (api *API) WS(c *gin.Context) {

	hubID := c.Param("hub")

	var hub server.Hub

	if hubID == "lobby" {

		hub = api.lobby

	} else {

		m, ok := api.matchSystem.Matches.Find(hubID)

		if !ok {
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
func (api *API) CardsHandler(c *gin.Context) {
	c.JSON(200, GetCache())
}

// GetDeckHandler returns a single deck, if public
func (api *API) GetDeckHandler(c *gin.Context) {

	deckUID := c.Param("id")

	var deck db.Deck

	err := db.Collection("decks").FindOne(
		context.Background(),
		bson.M{"uid": deckUID, "public": true},
	).Decode(&deck)

	if err != nil {
		c.Status(404)
		return
	}

	var user db.User

	err = db.Collection("users").FindOne(
		context.Background(),
		bson.M{"uid": deck.Owner},
	).Decode(&user)

	if err != nil {
		c.Status(404)
		return
	}

	deck.Owner = user.Username

	c.JSON(200, deck)

}

// GetDecksHandler returns an array of the users decks
func (api *API) GetDecksHandler(c *gin.Context) {

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
func (api *API) CreateDeckHandler(c *gin.Context) {

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

		if decksCount >= 50 {
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

// DeleteDeckHandler deletes the specified deck
func (api *API) DeleteDeckHandler(c *gin.Context) {

	user, err := db.GetUserForToken(c.GetHeader("Authorization"))
	if err != nil {
		c.Status(401)
		return
	}

	deckUID := c.Param("id")

	result, err := db.Collection("decks").DeleteOne(
		context.Background(),
		bson.M{"uid": deckUID, "owner": user.UID},
	)

	if err != nil {
		c.Status(401)
		return
	}

	if result.DeletedCount < 1 {
		c.Status(401)
		return
	}

	c.Status(200)

}

func (api *API) GetMatchHandler(c *gin.Context) {

	m, ok := api.matchSystem.Matches.Find(c.Param("id"))

	if !ok {
		c.Status(404)
		return
	}

	c.JSON(200, bson.M{"name": m.MatchName, "host": m.HostID, "started": m.Started})

}

// InviteHandler handles duel invitations
func (api *API) InviteHandler(c *gin.Context) {

	var res string

	match, ok := api.matchSystem.Matches.Find(c.Param("id"))

	if !ok {
		res = fmt.Sprintf(`<!DOCTYPE html>
<html>
	<head>
		<title>Redirecting you..</title>
		<meta property="og:type" content="website" />
		<meta name="og:title" property="og:title" content="Invitation expired!">
		<meta name="og:description" property="og:description" content="This duel is no longer available">
		<meta name="og:image" property="og:image" content="https://i.imgur.com/g4I6jEL.png">
		<meta name="og:url" property="og:url" content="https://shobu.io/invite/%s" />
	</head>
	<body style="background: #36393F">
		<p>Please wait while we redirect you.. Make sure javascript is enabled.</p>
		<script>if(!navigator.userAgent.includes("discord")) { window.location.replace("/overview"); }</script>
	</body>
</html>
		`, c.Param("id"))
	} else if match.Started {
		res = fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
			<head>
				<title>Redirecting you..</title>
				<meta property="og:type" content="website" />
				<meta name="og:title" property="og:title" content="Invitation expired! The duel has already begun.">
				<meta name="og:description" property="og:description" content="%s is duelling %s!">
				<meta name="og:image" property="og:image" content="https://i.imgur.com/qdOnH8k.png">
				
			</head>
			<body>
				<p>Please wait while we redirect you.. Make sure javascript is enabled.</p>
				<script>if(!navigator.userAgent.includes("discord")) { window.location.replace("/overview"); }</script>
			</body>
		</html>
		`, match.Player1.Socket.User.Username, match.Player2.Socket.User.Username)
	} else if match.Player1 != nil {
		res = fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
			<head>
				<title>Redirecting you..</title>
				<meta property="og:type" content="website" />
				<meta name="og:title" property="og:title" content="%s invited you to a duel!">
				<meta name="og:image" property="og:image" content="https://i.imgur.com/8PlN43q.png">
				
			</head>
			<body>
				<p>Please wait while we redirect you.. Make sure javascript is enabled.</p>
				<script>if(!navigator.userAgent.includes("discord")) { window.location.replace("/duel/%s"); }</script>
			</body>
		</html>
		`, match.Player1.Socket.User.Username, c.Param("id"))
	} else {
		res = `
		<!DOCTYPE html>
		<html>
			<head>
				<title>Redirecting you..</title>
				<meta property="og:type" content="website" />
				<meta name="og:title" property="og:title" content="Invitation is loading..">
				<meta name="og:description" property="og:description" content="This duel is in the progress of being created">
				<meta name="og:image" property="og:image" content="https://i.imgur.com/FEiBdKe.png">
				
			</head>
			<body>
				<p>Please wait while we redirect you.. Make sure javascript is enabled.</p>
				<script>if(!navigator.userAgent.includes("discord")) { window.location.replace("/overview"); }</script>
			</body>
		</html>
		`
	}

	c.Data(200, "text/html; charset=utf-8", []byte(res))
}
