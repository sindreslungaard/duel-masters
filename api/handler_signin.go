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
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (api *API) signinHandler(w http.ResponseWriter, r *http.Request) {
	var body signinReqBody

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	username, err1 := assert.Is(body.Username).NotEmpty().String()
	password, err2 := assert.Is(body.Password).NotEmpty().String()

	err = assert.First(err1, err2)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user db.User

	if err := db.Users.FindOne(r.Context(), bson.M{"username": primitive.Regex{Pattern: "^" + username + "$", Options: "i"}}).Decode(&user); err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return
	}

	// Check for IP ban
	ip := getIP(r)

	bans, err := db.Users.CountDocuments(r.Context(), bson.M{
		"$or": []bson.M{
			{"type": db.UserBan, "value": user.UID},
			{"type": db.IPBan, "value": ip},
		},
	})

	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if bans > 0 {
		http.Error(w, "Banned", http.StatusForbidden)
		return
	}

	token, err := uuid.NewRandom()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session := db.UserSession{
		Token:   token.String(),
		IP:      ip,
		Expires: int(time.Now().Add(time.Second * 2592000).Unix()),
	}

	db.Users.UpdateOne(r.Context(), bson.M{"uid": user.UID}, bson.M{"$push": bson.M{"sessions": session}})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Json{
		"user":  user,
		"token": session.Token,
	})

	// TODO: Remove expired/unneeded sessions from db
}
