package api

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				logrus.Error(err)
				write(w, http.StatusInternalServerError, Json{"message": "Something unexpected happened"})
			}

		}()

		next.ServeHTTP(w, r)
	})
}
