package server

import (
	"net/http"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// Hub is an interface that accepts incoming websocket messages
type Hub interface {
	Parse(s *Socket, data []byte)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func ws(c *gin.Context) {

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		c.Status(500)
		return
	}

	user, err := GetUserForToken(c)
	if err != nil {
		c.Status(401)
		return
	}

	s := newSocket(conn, user)

	// Handle the connection in a new goroutine to free up this memory
	go s.listen()

}

// Start initiates the server
func Start(port string) {

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	r.Use(gin.Recovery())

	// cors
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

	r.GET("/ws", ws)
	r.POST("/api/auth/signin", SigninHandler)
	r.POST("/api/auth/signup", SignupHandler)

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

	logrus.Infof("Listening on port %s", port)

	logrus.Fatal(r.Run(":" + port))

}
