package api

import (
	"crypto/subtle"
	"net/http"
	"os"
	"strings"
)

func authorizeServerRequest(w http.ResponseWriter, r *http.Request) bool {
	secret := os.Getenv("secret")
	if secret == "" {
		write(w, http.StatusInternalServerError, Json{"message": "Server is misconfigured"})
		return false
	}

	authorization := r.Header.Get("Authorization")
	provided, ok := strings.CutPrefix(authorization, "Bearer ")
	if !ok || subtle.ConstantTimeCompare([]byte(provided), []byte(secret)) != 1 {
		write(w, http.StatusUnauthorized, Json{"message": "Unauthorized"})
		return false
	}

	return true
}
