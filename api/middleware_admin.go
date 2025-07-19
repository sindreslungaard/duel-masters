package api

import (
	"duel-masters/db"
	"net/http"
)

func adminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := db.GetUserForToken(r.Header.Get("Authorization"))

		if err != nil {
			write(w, http.StatusUnauthorized, Json{"message": "Unauthorized"})
			return
		}

		ok := false
		for _, perm := range user.Permissions {
			if perm == "admin" {
				ok = true
			}
		}

		if !ok {
			write(w, http.StatusForbidden, Json{"message": "Forbidden"})
			return
		}

		next.ServeHTTP(w, r)
	})
}
