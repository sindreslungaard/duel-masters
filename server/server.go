package server

// Hub is an interface that accepts incoming websocket messages
type Hub interface {
	Parse(s *Socket, data []byte)
	Name() string
	OnSocketClose(s *Socket)
}
