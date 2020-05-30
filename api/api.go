package api

import (
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Start starts the API
func Start(port string) {

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
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Main routes
	r.GET("/ws/:hub", WS)
	r.POST("/api/auth/signin", SigninHandler)
	r.POST("/api/auth/signup", SignupHandler)
	r.POST("/api/match", MatchHandler)
	r.GET("/api/cards", CardsHandler)
	r.GET("/api/decks", GetDecksHandler)
	r.POST("/api/decks", CreateDeckHandler)

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
