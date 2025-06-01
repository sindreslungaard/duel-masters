package scenario

import (
	"duel-masters/game/match"
	"duel-masters/server"
	"encoding/json"
)

type MockClient struct {
	socket *server.Socket
	match  *match.Match
}

func NewMockClient(s *server.Socket, m *match.Match) *MockClient {
	return &MockClient{
		socket: s,
		match:  m,
	}
}

func (c *MockClient) SendEndTurnMsg() {
	msg, err := json.Marshal(server.Message{Header: "end_turn"})

	if err != nil {
		panic(err)
	}

	c.match.Parse(c.socket, msg)
}
