package api

import (
	"duel-masters/db"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sindreslungaard/assert"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type signinReqBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (api *API) signinHandler(w http.ResponseWriter, r *http.Request) {
	var body signinReqBody

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		write(w, http.StatusBadRequest, Json{"message": err.Error()})
		return
	}

	username, err1 := assert.Is(body.Username).NotEmpty().String()
	password, err2 := assert.Is(body.Password).NotEmpty().String()

	if err := assert.First(err1, err2); err != nil {
		write(w, http.StatusBadRequest, Json{"message": err.Error()})
		return
	}

	var user db.User

	if err := db.Users().FindOne(r.Context(), bson.M{"username": primitive.Regex{Pattern: "^" + username + "$", Options: "i"}}).Decode(&user); err != nil {
		write(w, http.StatusNotFound, Json{"message": "User not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		write(w, http.StatusUnauthorized, Json{"message": "Incorrect password"})
		return
	}

	// Check for IP ban
	ip := getIP(r)

	bans, err := db.Users().CountDocuments(r.Context(), bson.M{
		"$or": []bson.M{
			{"type": db.UserBan, "value": user.UID},
			{"type": db.IPBan, "value": ip},
		},
	})

	if err != nil {
		logrus.Error(err)
		write(w, http.StatusInternalServerError, Json{"message": err.Error()})
		return
	}

	if bans > 0 {
		write(w, http.StatusInternalServerError, Json{"message": "Banned"})
		http.Error(w, "Banned", http.StatusForbidden)
		return
	}

	token, err := uuid.NewRandom()
	if err != nil {
		write(w, http.StatusInternalServerError, Json{"message": err.Error()})
		return
	}

	session := db.UserSession{
		Token:   token.String(),
		IP:      ip,
		Expires: int(time.Now().Add(time.Second * 2592000).Unix()),
	}

	db.Users().UpdateOne(r.Context(), bson.M{"uid": user.UID}, bson.M{"$push": bson.M{"sessions": session}})

	write(w, http.StatusOK, Json{
		"user":  user,
		"token": session.Token,
	})

	// TODO: Remove expired/unneeded sessions from db
}
