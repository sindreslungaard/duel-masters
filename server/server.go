package server

import (
	"net/http"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func ws(c *gin.Context) {

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		// TODO: handle
		return
	}

	s := newSocket(conn, 0)

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

	r.GET("/ws", ws)

	// route everything else to our SPA
	r.NoRoute(func(c *gin.Context) {
		c.File(path.Join(dir, "webapp", "dist", "index.html"))
	})

	logrus.Infof("Listening on port %s", port)

	logrus.Fatal(r.Run(":" + port))

}
