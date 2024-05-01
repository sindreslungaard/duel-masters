package api

import (
	"fmt"
	"net/http"

	"duel-masters/server"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// WS handles websocket upgrade
func (api *API) WS(c *gin.Context) {

	hubID := c.Param("hub")

	var hub server.Hub

	if hubID == "lobby" {

		hub = api.lobby

	} else {

		m, ok := api.matchSystem.Matches.Find(hubID)

		if !ok {
			c.Status(404)
			return
		}

		hub = m

	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		c.Status(500)
		return
	}

	s := server.NewSocket(conn, hub)

	// Handle the connection in a new goroutine to free up this memory
	go s.Listen()

}

// InviteHandler handles duel invitations
func (api *API) InviteHandler(c *gin.Context) {
	html := fmt.Sprintf(`
	<html>
		<head>
			<title>Redirecting you..</title>
			<meta property="og:type" content="website" />
			<meta name="og:title" property="og:title" content="Duel invite">
			<meta name="og:description" property="og:description" content="You have been invited to a duel">
			<meta name="og:image" property="og:image" content="https://i.imgur.com/8PlN43q.png">
		</head>
		<body style="background: #36393F">
			<p>Please wait while we redirect you.. Make sure javascript is enabled.</p>
			<script>if(!navigator.userAgent.includes("discord")) { window.location.replace("/overview?invite=%s"); }</script>
		</body>
	</html>	
	`, c.Param("id"))
	c.Data(200, "text/html; charset=utf-8", []byte(html))
}
