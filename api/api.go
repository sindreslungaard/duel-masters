package api

import (
	"duel-masters/game"
	"duel-masters/game/match"
	"net"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

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

	dir, err := os.Getwd()
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

	logrus.Fatal(r.Run(":" + os.Getenv("port")))
}
