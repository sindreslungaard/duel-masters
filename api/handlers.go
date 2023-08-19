package api

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"duel-masters/db"
	"duel-masters/flags"
	"duel-masters/game"
	"duel-masters/game/match"
	"duel-masters/internal"
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

	var user db.User

	if err := db.Users.FindOne(context.TODO(), bson.M{"username": primitive.Regex{Pattern: "^" + reqBody.Username + "$", Options: "i"}}).Decode(&user); err != nil {
		c.Status(401)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqBody.Password)); err != nil {
		c.Status(401)
		return
	}

	// Check for IP ban
	ip := c.ClientIP()

	bans, err := db.Users.CountDocuments(context.Background(), bson.M{
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

	db.Users.UpdateOne(context.TODO(), bson.M{"uid": user.UID}, bson.M{"$push": bson.M{"sessions": session}})

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
	ip := c.ClientIP()

	bans, err := db.Bans.CountDocuments(context.Background(), bson.M{"type": db.IPBan, "value": ip})

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

	if err := db.Users.FindOne(context.TODO(), bson.M{"username": primitive.Regex{Pattern: "^" + reqBody.Username + "$", Options: "i"}}).Decode(&db.User{}); err == nil {
		c.JSON(400, bson.M{"message": "The username is already taken"})
		return
	}

	if err := db.Users.FindOne(context.TODO(), bson.M{"email": primitive.Regex{Pattern: "^" + reqBody.Email + "$", Options: "i"}}).Decode(&db.User{}); err == nil {
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

	_, err = db.Users.InsertOne(context.TODO(), user)

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

// MatchHandler handles creation of new mathes
func (api *API) MatchHandler(c *gin.Context) {

	if !flags.NewMatchesEnabled {
		c.Status(403)
		return
	}

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
		name = game.DefaultMatchNames[rand.Intn(len(game.DefaultMatchNames))]
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

	err := db.Decks.FindOne(
		context.Background(),
		bson.M{"uid": deckUID, "public": true},
	).Decode(&deck)

	if err != nil {
		c.Status(404)
		return
	}

	var user db.User

	err = db.Users.FindOne(
		context.Background(),
		bson.M{"uid": deck.Owner},
	).Decode(&user)

	if err != nil {
		c.Status(404)
		return
	}

	deck.Owner = user.Username

	d, err := match.ConvertToLegacyDeck(deck)

	if err != nil {
		c.Status(404)
	}

	c.JSON(200, d)

}

// GetDecksHandler returns an array of the users decks
func (api *API) GetDecksHandler(c *gin.Context) {

	user, err := db.GetUserForToken(c.GetHeader("Authorization"))
	if err != nil {
		c.Status(401)
		return
	}

	cur, err := db.Decks.Find(context.TODO(), bson.M{
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

	legacyDecks := []db.LegacyDeck{}

	for _, deck := range decks {
		legacyDeck, err := match.ConvertToLegacyDeck(deck)
		if err != nil {
			continue
		}
		legacyDecks = append(legacyDecks, legacyDeck)
	}

	c.JSON(200, legacyDecks)

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

	if len(reqBody.UID) < 1 {

		// New deck

		decksCount, err := db.Decks.CountDocuments(context.TODO(), bson.M{"owner": user.UID})

		if err != nil {
			logrus.Error(err)
			c.Status(500)
			return
		}

		if decksCount >= 200 {
			c.Status(403)
			return
		}

		deck := db.LegacyDeck{
			UID:      uuid.New().String(),
			Owner:    user.UID,
			Name:     reqBody.Name,
			Public:   reqBody.Public,
			Standard: false,
			Cards:    reqBody.Cards,
		}

		_, err = db.Decks.InsertOne(context.TODO(), match.ConvertFromLegacyDeck(deck))

		if err != nil {
			c.Status(500)
			return
		}

	} else {

		// Edit deck

		deck := match.ConvertFromLegacyDeck(db.LegacyDeck{
			UID:      reqBody.UID,
			Owner:    user.UID,
			Name:     reqBody.Name,
			Public:   reqBody.Public,
			Standard: false,
			Cards:    reqBody.Cards,
		})

		_, err := db.Decks.UpdateOne(
			context.TODO(),
			bson.M{"uid": reqBody.UID, "owner": user.UID},
			bson.M{"$set": bson.M{"name": reqBody.Name, "public": reqBody.Public, "cards": deck.Cards}},
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

	result, err := db.Decks.DeleteOne(
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

type changePasswordReqBody struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6"`
}

func (api *API) ChangePasswordHandler(c *gin.Context) {

	user, err := db.GetUserForToken(c.GetHeader("Authorization"))
	if err != nil {
		c.Status(401)
		return
	}

	var reqBody changePasswordReqBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(400, bson.M{"error": "New password must be at least 6 characters long"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqBody.OldPassword)); err != nil {
		c.JSON(401, bson.M{"error": "Old password is incorrect"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(reqBody.NewPassword), 10)

	if err != nil {
		c.Status(500)
		return
	}

	db.Users.UpdateOne(context.TODO(), bson.M{"uid": user.UID}, bson.M{"$set": bson.M{"password": hash}})

	c.JSON(200, bson.M{"message": "Successfully changed your password"})
}

func (api *API) GetPreferencesHandler(c *gin.Context) {
	user, err := db.GetUserForToken(c.GetHeader("Authorization"))
	if err != nil {
		c.Status(401)
		return
	}

	c.JSON(200, bson.M{
		"playmat": user.Playmat,
	})
}

type preferencesReqBody struct {
	Playmat string `json:"playmat"`
}

func (api *API) UpdatePreferencesHandler(c *gin.Context) {

	user, err := db.GetUserForToken(c.GetHeader("Authorization"))
	if err != nil {
		c.Status(401)
		return
	}

	var reqBody preferencesReqBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(400, bson.M{"error": "New password must be at least 6 characters long"})
		return
	}

	if reqBody.Playmat != "" && !strings.HasPrefix(reqBody.Playmat, "https://i.imgur.com/") {
		c.JSON(400, bson.M{"error": "Playmat images must be uploaded to imgur and the url must start with https://i.imgur.com/. Make sure the link includes the file extension (.png, .jpg)"})
		return
	}

	db.Users.UpdateOne(context.Background(), bson.M{
		"uid": user.UID,
	}, bson.M{"$set": bson.M{
		"playmat": reqBody.Playmat,
	}})

	c.JSON(200, bson.M{"message": "Successfully saved your preferences"})

}

type recoverPasswordReqBody struct {
	Email string `json:"email"`
}

func (api *API) RecoverPasswordHandler(c *gin.Context) {

	if internal.RateLimited(fmt.Sprintf("%s/recoverpw", c.ClientIP()), 3, 1000*60*15) {
		c.JSON(400, bson.M{"error": "Please wait a while before requesting to recover password again"})
		return
	}

	var reqBody recoverPasswordReqBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(400, bson.M{"error": "Please provide a valid email"})
		return
	}

	genericResponse := "If the email you specified matches any registered users you will soon receive a mail with a link to reset your password"

	var user db.User

	if err := db.Users.FindOne(context.TODO(), bson.M{"email": primitive.Regex{Pattern: "^" + reqBody.Email + "$", Options: "i"}}).Decode(&user); err != nil {
		logrus.Debug("Attempt at recovering password with email that does not belong to any users ", reqBody.Email)
		c.JSON(200, bson.M{"message": genericResponse})
		return
	}

	code, err := internal.RandomString(50)
	code = fmt.Sprintf("%v-%s", time.Now().Unix(), code)

	if err != nil {
		logrus.Error("Error generating password recovery code", err)
		c.JSON(500, bson.M{"error": "Something went wrong"})
		return
	}

	db.Users.UpdateOne(context.Background(), bson.M{
		"uid": user.UID,
	}, bson.M{"$set": bson.M{
		"recoverycode": code,
	}})

	err = internal.SendMail(user.Email, "Recover your password", fmt.Sprintf(`
	Use the link below to recover the password for your account <b>%s</b>
	<br><br>
	https://shobu.io/recover-password/%s
	<br><br>
	If you did not request to reset your password, please disregard this email
	`, user.Username, code))

	if err != nil {
		logrus.Error("Failed to send email", err)
		c.JSON(500, bson.M{"error": "Something went wrong"})
		return
	}

	c.JSON(200, bson.M{"message": genericResponse})

}

type resetPasswordReqBody struct {
	Code     string `json:"code"`
	Password string `json:"password"`
}

func (api *API) ResetPasswordHandler(c *gin.Context) {

	var reqBody resetPasswordReqBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(400, bson.M{"error": "Invalid payload"})
		return
	}

	if len(reqBody.Password) < 6 {
		c.JSON(400, bson.M{"error": "Password must be at least 6 characters long"})
		return
	}

	if len(reqBody.Code) < 30 {
		c.JSON(400, bson.M{"error": "Invalid code"})
		return
	}

	var user db.User

	if err := db.Users.FindOne(context.TODO(), bson.M{"recoverycode": reqBody.Code}).Decode(&user); err != nil {
		c.JSON(400, bson.M{"error": "Invalid or expired code"})
		return
	}

	ts, err := strconv.Atoi(strings.Split(reqBody.Code, "-")[0])

	if err != nil {
		logrus.Error("Failed to parse recovery code", reqBody.Code)
		c.JSON(400, bson.M{"error": "Could not parse recovery code"})
		return
	}

	if int64(ts)+86400 < time.Now().Unix() {
		c.JSON(400, bson.M{"error": "Recovery code has expired"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), 10)

	if err != nil {
		logrus.Error("Failed to generate password hash during password reset")
		c.JSON(500, bson.M{"error": "Something unexpected happened"})
		return
	}

	db.Users.UpdateOne(context.Background(), bson.M{
		"uid": user.UID,
	}, bson.M{
		"$set": bson.M{
			"password": string(hash),
		},
		"$unset": bson.M{
			"recoverycode": "",
		},
	})

	c.JSON(200, bson.M{"message": fmt.Sprintf("Password for the account \"%s\" was successfully changed", user.Username)})

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
