package api

import (
	"duel-masters/db"
	"encoding/json"
	"net/http"

	"github.com/sindreslungaard/assert"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type changePasswordReqBody struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

func (api *API) changePasswordHandler(w http.ResponseWriter, r *http.Request) {
	user, err := db.GetUserForToken(r.Header.Get("Authorization"))
	if err != nil {
		write(w, http.StatusUnauthorized, Json{"message": "Unauthorized"})
		return
	}

	var body changePasswordReqBody

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		write(w, http.StatusBadRequest, Json{"message": err.Error()})
		return
	}

	oldPassword, err1 := assert.Is(body.OldPassword).NotEmpty().String()
	newPassword, err2 := assert.Is(body.NewPassword).NotEmpty().MinLen(6).MaxLen(255).String()

	if err := assert.First(err1, err2); err != nil {
		write(w, http.StatusBadRequest, Json{"message": "New password must be at least 6 characters long"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		write(w, http.StatusForbidden, Json{"message": "Old password is incorrect"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), 10)

	if err != nil {
		write(w, http.StatusInternalServerError, Json{"message": "Something went wrong"})
		return
	}

	db.Users.UpdateOne(r.Context(), bson.M{"uid": user.UID}, bson.M{"$set": bson.M{"password": hash}})

	write(w, http.StatusOK, Json{"message": "Successfully changed your password"})
}
