package api

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strings"

	"duel-masters/db"
	"duel-masters/flags"
	"duel-masters/game"
	"duel-masters/game/match"
	"duel-masters/server"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

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

	m := api.matchSystem.NewMatch(name, user.UID, visible, false)

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

// InviteHandler handles duel invitations
func (api *API) InviteHandler(c *gin.Context) {
	html := fmt.Sprintf(`
	<html>
		<head>
			<title>Redirecting you..</title>
			<meta property="og:type" content="website" />
			<meta name="og:title" property="og:title" content="Duel invite">
			<meta name="og:description" property="og:description" content="You have been invited to a duel">
			<meta name="og:image" property="og:image" content="https://i.imgur.com/8PlN43q.png">
		</head>
		<body style="background: #36393F">
			<p>Please wait while we redirect you.. Make sure javascript is enabled.</p>
			<script>if(!navigator.userAgent.includes("discord")) { window.location.replace("/overview?invite=%s"); }</script>
		</body>
	</html>	
	`, c.Param("id"))
	c.Data(200, "text/html; charset=utf-8", []byte(html))
}
