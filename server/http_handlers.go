package server

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

func requesterUserID(r *http.Request) (int, error) {

	authorization := r.Header.Get("Authorization")

	if len(authorization) < 1 {
		return 0, errors.New("Missing authorization header")
	}

}

type AuthenticationRequestBody struct {
	username string
	password string
}

func AuthenticationHandler(w http.ResponseWriter, r *http.Request) {

	rawBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(400)
		return
	}

	var body AuthenticationRequestBody

	err = json.Unmarshal(rawBody, &body)

	if err != nil {
		w.WriteHeader(400)
		return
	}

}
