package api

import (
	"context"
	"duel-masters/db"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/sindreslungaard/assert"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type resetPasswordReqBody struct {
	Code     string `json:"code"`
	Password string `json:"password"`
}

func (api *API) resetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var body resetPasswordReqBody

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		write(w, http.StatusBadRequest, Json{"message": err.Error()})
		return
	}

	code, err := assert.Is(body.Code).NotEmpty().MinLen(30).String()

	if err != nil {
		write(w, http.StatusBadRequest, Json{"message": "Invalid code"})
		return
	}

	password, err := assert.Is(body.Password).NotEmpty().MinLen(6).MaxLen(255).String()

	if err != nil {
		write(w, http.StatusBadRequest, Json{"message": "Password must be at least 6 characters"})
		return
	}

	var user db.User

	if err := db.Users().FindOne(r.Context(), bson.M{"recoverycode": code}).Decode(&user); err != nil {
		write(w, http.StatusBadRequest, Json{"message": "Invalid or expired code"})
		return
	}

	ts, err := strconv.Atoi(strings.Split(code, "-")[0])

	if err != nil {
		logrus.Error("Failed to parse recovery code", code)
		write(w, http.StatusBadRequest, Json{"message": "Could not parse recovery code"})
		return
	}

	if int64(ts)+86400 < time.Now().Unix() {
		write(w, http.StatusBadRequest, Json{"message": "Recovery code has expired"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		logrus.Error("Failed to generate password hash during password reset")
		write(w, http.StatusInternalServerError, Json{"message": "Something unexpected happened"})
		return
	}

	db.Users().UpdateOne(context.Background(), bson.M{
		"uid": user.UID,
	}, bson.M{
		"$set": bson.M{
			"password": string(hash),
		},
		"$unset": bson.M{
			"recoverycode": "",
		},
	})

	write(w, http.StatusOK, Json{"message": fmt.Sprintf("Password for the account \"%s\" was successfully changed", user.Username)})
}
