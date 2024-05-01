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
	mux             *http.ServeMux
	middleware      []func(http.Handler) http.Handler
	blockedNetworks IPRange

	lobby       *game.Lobby
	matchSystem *match.MatchSystem
}

func New(lobby *game.Lobby, matchSystem *match.MatchSystem) *API {
	return &API{
		mux:        http.NewServeMux(),
		middleware: []func(http.Handler) http.Handler{},
		blockedNetworks: IPRange{
			cidrs: []*net.IPNet{},
		},

		matchSystem: matchSystem,
		lobby:       lobby,
	}
}

func (api *API) SetBlockedIPs(iprange IPRange) {
	api.blockedNetworks = iprange
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
	addr := "127.0.0.1:" + port

	api.Use(recoverMiddleware)
	api.Use(loggingMiddleware)
	api.Use(corsMiddleware)

	api.HandleFunc("POST /api/auth/signin", api.signinHandler)
	api.HandleFunc("POST /api/auth/signup", api.signupHandler)
	api.HandleFunc("POST /api/auth/recover", api.recoverPasswordHandler)
	api.HandleFunc("POST /api/auth/reset", api.resetPasswordHandler)
	api.HandleFunc("POST /api/auth/reset-password", api.changePasswordHandler)
	api.HandleFunc("GET /api/preferences", api.getPreferencesHandler)

	server := &http.Server{
		Addr:    addr,
		Handler: api.mux,
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
