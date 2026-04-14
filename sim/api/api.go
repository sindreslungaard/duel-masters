package api

import (
	"duel-masters/game"
	"duel-masters/game/match"
	"encoding/json"
	"net"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

type Json map[string]interface{}

type API struct {
	mux        *http.ServeMux
	middleware []func(http.Handler) http.Handler

	lobby       *game.Lobby
	matchSystem *match.MatchSystem
}

func New(lobby *game.Lobby, matchSystem *match.MatchSystem) *API {
	return &API{
		mux:        http.NewServeMux(),
		middleware: []func(http.Handler) http.Handler{},

		lobby:       lobby,
		matchSystem: matchSystem,
	}
}

func (api *API) Use(middleware func(http.Handler) http.Handler) {
	api.middleware = append(api.middleware, middleware)
}

func (api *API) HandleFunc(pattern string, handler http.HandlerFunc) {
	h := http.Handler(handler)
	for i := len(api.middleware) - 1; i >= 0; i-- {
		h = api.middleware[i](h)
	}
	api.mux.Handle(pattern, h)
}

// Start starts the API
func (api *API) Start(port string) {
	/* dir, err := os.Getwd()
	if err != nil {
		panic(err)
	} */

	api.Use(recoverMiddleware)
	api.Use(loggingMiddleware)
	api.Use(corsMiddleware)

	api.HandleFunc("GET /ws/{hub}", api.websocketHandler)
	api.HandleFunc("GET /api/cards", api.getCardsHandler)
	api.HandleFunc("POST /api/match", api.createMatchHandler)

	// webapp
	/* 	dist := http.FileServer(http.Dir(path.Join(dir, "webapp", "dist")))
	   	api.mux.Handle("GET /assets/{path...}", dist)
	   	api.mux.Handle("GET /css/{path...}", dist)
	   	api.mux.Handle("GET /js/{path...}", dist)
	   	api.mux.HandleFunc("GET /favicon.ico", api.staticFileHandler(path.Join(dir, "webapp", "dist", "favicon.ico")))
	   	api.mux.HandleFunc("GET /{path...}", api.staticFileHandler(path.Join(dir, "webapp", "dist", "index.html")))
	*/
	server := &http.Server{
		Addr:    ":" + port,
		Handler: api.mux,
	}

	logrus.Infof("Listening on port %s", port)
	logrus.Fatal(server.ListenAndServe())
}

func getIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	ip = strings.TrimSpace(strings.Split(ip, ",")[0])
	if ip != "" {
		return ip
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}

func write(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
