package api

import (
	"duel-masters/db"
	"duel-masters/game/match"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/sindreslungaard/assert"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

type createDeckBody struct {
	Name   string   `json:"name" binding:"required,min=1,max=30"`
	Cards  []string `json:"cards" binding:"required"`
	UID    string   `json:"uid"`
	Public bool     `json:"public"`
}

func (api *API) createDeckHandler(w http.ResponseWriter, r *http.Request) {
	user, err := db.GetUserForToken(r.Header.Get("Authorization"))
	if err != nil {
		write(w, http.StatusUnauthorized, Json{"message": "Unauthorized"})
		return
	}

	var body createDeckBody

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		write(w, http.StatusBadRequest, Json{"message": err.Error()})
		return
	}

	name, err := assert.Is(body.Name).NotEmpty().MinLen(1).MaxLen(30).String()

	if err != nil {
		write(w, http.StatusBadRequest, Json{"message": err.Error()})
		return
	}

	if body.Cards == nil || len(body.Cards) < 40 || len(body.Cards) > 50 {
		write(w, http.StatusBadRequest, Json{"message": "Deck must consist of 40-50 cards"})
		return
	}

	for _, cuid := range body.Cards {
		if !CacheHas(cuid) {
			write(w, http.StatusBadRequest, Json{"message": "One or more of the cards specified are not valid"})
			return
		}
	}

	if len(body.UID) < 1 {
		// New deck
		decksCount, err := db.Decks.CountDocuments(r.Context(), bson.M{"owner": user.UID})

		if err != nil {
			logrus.Error(err)
			write(w, http.StatusInternalServerError, Json{"message": "Something went wrong"})
			return
		}

		if decksCount >= 200 {
			write(w, http.StatusForbidden, Json{"message": "Max number of decks reached"})
			return
		}

		deck := db.LegacyDeck{
			UID:      uuid.New().String(),
			Owner:    user.UID,
			Name:     name,
			Public:   body.Public,
			Standard: false,
			Cards:    body.Cards,
		}

		_, err = db.Decks.InsertOne(r.Context(), match.ConvertFromLegacyDeck(deck))

		if err != nil {
			write(w, http.StatusInternalServerError, Json{"message": "Something went wrong"})
			return
		}

	} else {

		// Edit deck

		deck := match.ConvertFromLegacyDeck(db.LegacyDeck{
			UID:      body.UID,
			Owner:    user.UID,
			Name:     name,
			Public:   body.Public,
			Standard: false,
			Cards:    body.Cards,
		})

		_, err := db.Decks.UpdateOne(
			r.Context(),
			bson.M{"uid": body.UID, "owner": user.UID},
			bson.M{"$set": bson.M{"name": name, "public": body.Public, "cards": deck.Cards}},
		)

		if err != nil {
			logrus.Error(err)
			write(w, http.StatusInternalServerError, Json{"message": "Something went wrong"})
			return
		}

	}

	write(w, http.StatusOK, Json{"message": "Deck created successfully"})
}
