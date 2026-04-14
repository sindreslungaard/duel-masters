package api

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.Debugf("Request: %s %s %s", r.Method, r.URL.Path, getIP(r))
		next.ServeHTTP(w, r)
	})
}
