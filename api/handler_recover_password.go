package api

import (
	"context"
	"duel-masters/db"
	"duel-masters/internal"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sindreslungaard/assert"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type recoverPasswordReqBody struct {
	Email string `json:"email"`
}

func (api *API) recoverPasswordHandler(w http.ResponseWriter, r *http.Request) {
	ip := getIP(r)
	if internal.RateLimited(fmt.Sprintf("%s/recoverpw", ip), 3, 1000*60*15) {
		write(w, http.StatusTooManyRequests, Json{"message": "Please wait a while before requesting to recover password again"})
		return
	}

	var body recoverPasswordReqBody

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		write(w, http.StatusBadRequest, Json{"message": err.Error()})
		return
	}

	email, err := assert.Is(body.Email).NotEmpty().Email().String()

	if err != nil {
		write(w, http.StatusBadRequest, Json{"message": "Invalid email"})
		return
	}

	genericResponse := "If the email you specified matches any registered users you will soon receive a mail with a link to reset your password"

	var user db.User

	if err := db.Users.FindOne(context.TODO(), bson.M{"email": primitive.Regex{Pattern: "^" + email + "$", Options: "i"}}).Decode(&user); err != nil {
		logrus.Debug("Attempt at recovering password with email that does not belong to any users ", email)
		write(w, http.StatusOK, Json{"message": genericResponse})
		return
	}

	code, err := internal.RandomString(50)
	code = fmt.Sprintf("%v-%s", time.Now().Unix(), code)

	if err != nil {
		logrus.Error("Error generating password recovery code", err)
		write(w, http.StatusInternalServerError, Json{"message": "Something went wrong"})
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
		write(w, http.StatusInternalServerError, Json{"message": "Something went wrong"})
		return
	}

	write(w, http.StatusOK, Json{"message": genericResponse})
}
