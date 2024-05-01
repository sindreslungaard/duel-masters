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
	lobby       *game.Lobby
	matchSystem *match.MatchSystem

	blockedNetworks IPRange
}

func New(lobby *game.Lobby, matchSystem *match.MatchSystem) *API {
	return &API{
		matchSystem: matchSystem,
		lobby:       lobby,

		blockedNetworks: IPRange{
			cidrs: []*net.IPNet{},
		},
	}
}

func (api *API) SetBlockedIPs(iprange IPRange) {
	api.blockedNetworks = iprange
}

// Start starts the API
func (api *API) Start(port string) {
	addr := "127.0.0.1:" + port
	mux := http.NewServeMux()

	mux.HandleFunc("POST /auth/signin", api.signinHandler)

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	logrus.Infof("Listening at %s", addr)
	logrus.Fatal(server.ListenAndServe())

	/* dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	r.Use(gin.Recovery())

	// CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Main routes
	r.GET("/ws/:hub", api.WS)
	r.POST("/api/auth/signin", api.SigninHandler)
	r.POST("/api/auth/signup", api.SignupHandler)
	r.POST("/api/auth/recover", api.RecoverPasswordHandler)
	r.POST("/api/auth/reset", api.ResetPasswordHandler)
	r.POST("/api/auth/reset-password", api.ChangePasswordHandler)
	r.GET("/api/preferences", api.GetPreferencesHandler)
	r.PUT("/api/preferences", api.UpdatePreferencesHandler)
	r.GET("/api/match/:id", api.GetMatchHandler)
	r.POST("/api/match", api.MatchHandler)
	r.GET("/api/cards", api.CardsHandler)
	r.GET("/api/deck/:id", api.GetDeckHandler)
	r.GET("/api/decks", api.GetDecksHandler)
	r.POST("/api/decks", api.CreateDeckHandler)
	r.DELETE("/api/deck/:id", api.DeleteDeckHandler)
	r.GET("/invite/:id", api.InviteHandler)

	// Because Gin does not provide an easy way to handle requests where the file does not exist
	// (NoRoute tests on specified routes, not if the file exists) we expose our webapp's folders manually..
	r.Static("/assets", path.Join(dir, "webapp", "dist", "assets"))
	r.Static("/css", path.Join(dir, "webapp", "dist", "css"))
	r.Static("/js", path.Join(dir, "webapp", "dist", "js"))
	r.StaticFile("/favicon.ico", path.Join(dir, "webapp", "dist", "favicon.ico"))

	// route everything else to our SPA
	r.NoRoute(func(c *gin.Context) {
		c.File(path.Join(dir, "webapp", "dist", "index.html"))
	})

	logrus.Infof("Listening on port %s", os.Getenv("port"))

	logrus.Fatal(r.Run(":" + os.Getenv("port"))) */
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

func write(w http.ResponseWriter, status int, data Json) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
