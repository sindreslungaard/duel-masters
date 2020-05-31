package game

import "duel-masters/server"

// Lobby struct is used to create a Hub that can parse messages from the websocket server
type Lobby struct{}

// Parse websocket messages
func (l *Lobby) Parse(s *server.Socket, data []byte) {

}
