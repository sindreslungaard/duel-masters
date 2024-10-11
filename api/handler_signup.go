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

type signupReqBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (api *API) signupHandler(w http.ResponseWriter, r *http.Request) {
	var body signupReqBody

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		write(w, http.StatusBadRequest, Json{"message": err.Error()})
		return
	}

	// TODO: recaptcha

	username, err1 := assert.Is(body.Username).NotEmpty().MinLen(3).MaxLen(20).AlphaNumeric().String()
	password, err2 := assert.Is(body.Password).NotEmpty().MinLen(6).MaxLen(255).String()
	email, err3 := assert.Is(body.Email).NotEmpty().Email().String()

	if err := assert.First(err1, err2, err3); err != nil {
		write(w, http.StatusBadRequest, Json{"message": "Username, email or password not valid"})
		return
	}

	// Check for IP ban
	ip := getIP(r)

	bans, err := db.Bans.CountDocuments(r.Context(), bson.M{"type": db.IPBan, "value": ip})

	if err != nil {
		logrus.Error(err)
		write(w, http.StatusInternalServerError, Json{"message": err.Error()})
		return
	}

	if bans > 0 {
		write(w, http.StatusForbidden, Json{"message": "Banned"})
		return
	}

	// Check if IP is in list of blocked networks
	if api.blockedNetworks.Contains(ip) {
		write(w, http.StatusForbidden, Json{"message": "Your network has been blocked from creating new accounts. If you believe this is wrong, please contact us on discord: https://discord.gg/FkPTE4p"})
		return
	}

	if err := db.Users.FindOne(r.Context(), bson.M{"username": primitive.Regex{Pattern: "^" + username + "$", Options: "i"}}).Decode(&db.User{}); err == nil {
		write(w, http.StatusBadRequest, Json{"message": "Username has already been taken"})
		return
	}

	if err := db.Users.FindOne(r.Context(), bson.M{"email": primitive.Regex{Pattern: "^" + email + "$", Options: "i"}}).Decode(&db.User{}); err == nil {
		write(w, http.StatusBadRequest, Json{"message": "Email has already been taken"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		write(w, http.StatusInternalServerError, Json{"message": err.Error()})
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

	user := db.User{
		UID:         uuid.New().String(),
		Username:    username,
		Email:       email,
		Password:    string(hash),
		Permissions: []string{},
		Sessions: []db.UserSession{
			session,
		},
	}

	_, err = db.Users.InsertOne(r.Context(), user)

	if err != nil {
		write(w, http.StatusInternalServerError, Json{"message": err.Error()})
		return
	}

	write(w, http.StatusOK, Json{"user": user, "token": session.Token})
}
