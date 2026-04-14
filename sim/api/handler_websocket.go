package api

import (
	"duel-masters/server"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type DuelSession struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (api *API) websocketHandler(w http.ResponseWriter, r *http.Request) {
	hubID := r.PathValue("hub")

	secret := os.Getenv("duel_token_secret")

	if secret == "" {
		write(w, http.StatusInternalServerError, Json{"message": "Server is misconfigured"})
		return
	}

	duelToken := r.URL.Query().Get("duelToken")

	if duelToken == "" {
		write(w, http.StatusUnauthorized, Json{"message": "Missing duelToken"})
		return
	}

	// Decode the JWT token
	token, err := jwt.ParseWithClaims(duelToken, &DuelSession{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret key used for signing
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		write(w, http.StatusUnauthorized, Json{"message": "Invalid duelToken"})
		return
	}

	// Extract the claims
	duelSession, ok := token.Claims.(*DuelSession)
	if !ok {
		write(w, http.StatusUnauthorized, Json{"message": "Invalid token claims"})
		return
	}

	var hub server.Hub

	m, ok := api.matchSystem.Matches.Find(hubID)

	if !ok {
		write(w, http.StatusNotFound, Json{"message": "Match not found"})
		return
	}

	hub = m

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		write(w, http.StatusInternalServerError, Json{"message": "Something went wrong"})
		return
	}

	s := server.NewSocket(conn, hub, duelSession.ID, duelSession.Username)

	// Handle the connection in a new goroutine to free up this memory
	go s.Listen()
}
