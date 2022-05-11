package game

import (
	"context"
	"duel-masters/db"
	"duel-masters/game/match"
	"duel-masters/internal"
	"duel-masters/server"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	messageBufferSize int = 100
)

// Lobby struct is used to create a Hub that can parse messages from the websocket server
type Lobby struct {
	pinnedMessages []string
	messages       []server.LobbyChatMessage
	messagesMutex  *sync.Mutex

	subscribers internal.ConcurrentDictionary[server.Socket]

	userCache server.UserListMessage

	matches func() []*match.Match
}

func NewLobby() *Lobby {
	return &Lobby{
		pinnedMessages: []string{},
		messages: append(make([]server.LobbyChatMessage, 0), server.LobbyChatMessage{
			Username:  "[Server]",
			Color:     "#777",
			Message:   fmt.Sprintf("Last server restart: %v. Have fun!", time.Now().Local().UTC().Format("Mon Jan 2 15:04:05 -0700 MST")),
			Timestamp: int(time.Now().Unix()),
		}),
		messagesMutex: &sync.Mutex{},

		subscribers: internal.NewConcurrentDictionary[server.Socket](),

		userCache: server.GetUserList(),

		matches: func() []*match.Match { return []*match.Match{} },
	}
}

// Name just returns "lobby", obligatory for a hub
func (l *Lobby) Name() string {
	return "lobby"
}

func (l *Lobby) SetMatchesFunc(f func() []*match.Match) {
	l.matches = f
}

// Broadcast sends a message to all subscribed sockets
func (l *Lobby) Broadcast(msg interface{}) {
	for _, subscriber := range l.subscribers.Iter() {
		go subscriber.Send(msg)
	}
}

// StartTicker starts the lobby ticker
func (l *Lobby) StartTicker() {

	ticker := time.NewTicker(30 * time.Second) // tick every 30 seconds

	defer ticker.Stop()

	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Recovered from lobby ticker. %v", r)
			debug.PrintStack()
		}
	}()

	for {

		select {
		case <-ticker.C:
			{
				l.UpdateUserCache()
				l.Broadcast(l.userCache)
				l.UpdatePinnedMessages()
			}
		}

	}
}

// UpdateUserCache updates the list of users online
func (l *Lobby) UpdateUserCache() {
	l.userCache = server.GetUserList()
}

func (l *Lobby) UpdatePinnedMessages() {

	l.messagesMutex.Lock()
	defer l.messagesMutex.Unlock()

	l.Broadcast(server.PinnedMessages{
		Header:   "pinned_messages",
		Messages: l.pinnedMessages,
	})

}

func (l *Lobby) PinMessage(message string) {
	l.messagesMutex.Lock()
	defer l.messagesMutex.Unlock()

	l.pinnedMessages = append(l.pinnedMessages, message)
}

// Parse websocket messages
func (l *Lobby) Parse(s *server.Socket, data []byte) {

	defer func() {
		if r := recover(); r != nil {
			logrus.Warnf("Recovered from parsing a message in lobby. %v", r)
			debug.PrintStack()
		}
	}()

	var message server.Message
	if err := json.Unmarshal(data, &message); err != nil {
		return
	}

	switch message.Header {

	case "subscribe":
		{
			for _, subscriber := range l.subscribers.Iter() {
				if subscriber == s {
					return
				}
			}

			_, ok := l.subscribers.Find(s.UID)

			if ok {
				return
			}

			l.subscribers.Add(s.UID, s)

			messages := server.LobbyChatMessages{
				Header:   "chat",
				Messages: []server.LobbyChatMessage{},
			}

			for _, msg := range l.messages {
				if msg.Removed && msg.Username != s.User.Username {
					continue
				}

				messages.Messages = append(messages.Messages, msg)
			}

			s.Send(messages)

			s.Send(server.PinnedMessages{
				Header:   "pinned_messages",
				Messages: l.pinnedMessages,
			})

			// Update and send user list
			l.UpdateUserCache()
			s.Send(l.userCache)

			// Send match list
			l.Broadcast(match.MatchList(l.matches()))

		}

	case "chat":
		{

			var msg struct {
				Message string `json:"message"`
			}

			if err := json.Unmarshal(data, &msg); err != nil {
				return
			}

			if len(msg.Message) < 1 {
				return
			}

			runes := []rune(msg.Message)
			if string(runes[0:1]) == "/" {
				l.handleChatCommand(s, msg.Message)
				return
			}

			l.messagesMutex.Lock()
			defer l.messagesMutex.Unlock()

			if len(l.messages) >= messageBufferSize {
				_, l.messages = l.messages[0], l.messages[1:]
			}

			chatMsg := server.LobbyChatMessage{
				Username:  s.User.Username,
				Color:     s.User.Color,
				Message:   msg.Message,
				Timestamp: int(time.Now().Unix()),
				Removed:   s.User.Chatblocked,
			}

			toBroadcast := server.LobbyChatMessages{
				Header:   "chat",
				Messages: []server.LobbyChatMessage{chatMsg},
			}

			l.messages = append(l.messages, chatMsg)

			if chatMsg.Removed {
				s.Send(toBroadcast)
			} else {
				l.Broadcast(toBroadcast)
			}

		}

	}

}

func chat(s *server.Socket, message string) {
	s.Send(server.LobbyChatMessages{
		Header: "chat",
		Messages: []server.LobbyChatMessage{
			{
				Username:  "[Server -> you]",
				Color:     "#777",
				Message:   message,
				Timestamp: int(time.Now().Unix()),
			},
		},
	})
}

func (l *Lobby) handleChatCommand(s *server.Socket, command string) {

	defer func() {
		if r := recover(); r != nil {
			chat(s, "An error occured while parsing your command")
		}
	}()

	if !strings.HasPrefix(command, "/") {
		return
	}

	parts := strings.Split(command, " ")

	hasRights := false

	for _, permission := range s.User.Permissions {
		if permission == "admin" {
			hasRights = true
		}
	}

	if !hasRights {
		chat(s, "Unknown command and/or missing privileges")
		return
	}

	switch strings.ToLower(strings.TrimPrefix(parts[0], "/")) {
	case "sockets":
		{
			message := ""
			for _, s := range server.Sockets.Iter() {
				if s.Ready() {
					if message != "" {
						message += ", "
					}
					message += s.User.Username
				}
			}
			chat(s, "Sockets: "+message)
		}

	case "matches":
		{
			message := ""
			for _, m := range l.matches() {
				if message != "" {
					message += ", "
				}

				message += m.ID
			}
			chat(s, "Matches: "+message)
		}

	case "shutdown":
		{
			logrus.Info("Shutdown command invoked")
			os.Exit(0)
		}

	case "chatblock":
		{
			if len(parts) < 2 {
				chat(s, "Missing command arguments")
				return
			}

			var user db.User

			if err := db.Collection("users").FindOne(context.Background(), bson.M{"username": primitive.Regex{Pattern: "^" + parts[1] + "$", Options: "i"}}).Decode(&user); err != nil {
				chat(s, fmt.Sprintf("Could not find the user \"%s\"", parts[1]))
				return
			}

			_, err := db.Collection("users").UpdateOne(
				context.Background(),
				bson.M{"uid": user.UID},
				bson.M{"$set": bson.M{"chat_blocked": true}},
			)

			if err != nil {
				chat(s, fmt.Sprintf("Failed to chatblock %s", user.Username))
			}

			l.messagesMutex.Lock()
			defer l.messagesMutex.Unlock()

			for i, msg := range l.messages {
				if msg.Username == user.Username {
					l.messages[i].Removed = true
				}
			}

			blockedSocket, ok := server.FindByUserUID(user.UID)

			if ok {
				blockedSocket.User.Chatblocked = true
			}

			chat(s, fmt.Sprintf("Successfully chatblocked %s", user.Username))
		}

	case "ban":
		{
			if len(parts) < 2 {
				chat(s, "Missing command arguments")
				return
			}

			var user db.User

			if err := db.Collection("users").FindOne(context.Background(), bson.M{"username": primitive.Regex{Pattern: "^" + parts[1] + "$", Options: "i"}}).Decode(&user); err != nil {
				chat(s, fmt.Sprintf("Could not find the user \"%s\"", parts[1]))
				return
			}

			banEntry := db.Ban{
				Type:  db.UserBan,
				Value: user.UID,
			}

			db.Collection("bans").InsertOne(context.Background(), banEntry)

			// clear banned user sessions
			db.Collection("users").UpdateOne(
				context.Background(),
				bson.M{"uid": user.UID},
				bson.M{"$set": bson.M{"sessions": []db.UserSession{}}},
			)

			// disconnect the banned user if online
			bannedSocket, ok := server.FindByUserUID(user.UID)

			if ok {
				bannedSocket.Close()
			}

			chat(s, fmt.Sprintf("Successfully banned %s (%s)", user.Username, user.UID))
		}

	case "ipban":
		{
			if len(parts) < 2 {
				chat(s, "Missing command arguments")
				return
			}

			var user db.User

			if err := db.Collection("users").FindOne(context.Background(), bson.M{"username": primitive.Regex{Pattern: "^" + parts[1] + "$", Options: "i"}}).Decode(&user); err != nil {
				chat(s, fmt.Sprintf("Could not find the user \"%s\"", parts[1]))
				return
			}

			banEntries := []interface{}{}

			banEntries = append(banEntries, db.Ban{
				Type:  db.UserBan,
				Value: user.UID,
			})

			if user.Sessions != nil && len(user.Sessions) > 0 {
				banEntries = append(banEntries, db.Ban{
					Type:  db.IPBan,
					Value: user.Sessions[len(user.Sessions)-1].IP,
				})
			}

			db.Collection("bans").InsertMany(context.Background(), banEntries)

			// clear banned user sessions
			db.Collection("users").UpdateOne(
				context.Background(),
				bson.M{"uid": user.UID},
				bson.M{"$set": bson.M{"sessions": []db.UserSession{}}},
			)

			// disconnect the banned user if online
			bannedSocket, ok := server.FindByUserUID(user.UID)

			if ok {
				bannedSocket.Close()
			}

			if len(banEntries) > 1 {
				chat(s, fmt.Sprintf("Successfully banned %s (%s), and their IP", user.Username, user.UID))
			} else {
				chat(s, fmt.Sprintf("Successfully banned %s (%s), but did not find an IP to ban", user.Username, user.UID))
			}

		}

	case "malloc":
		{
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			chat(s, fmt.Sprintf("Current = %v bytes", (m.Alloc)))
			chat(s, fmt.Sprintf("\tAll time = %v bytes", (m.TotalAlloc)))
			chat(s, fmt.Sprintf("\tReserved = %v bytes", (m.Sys)))
			chat(s, fmt.Sprintf("\tNumGC = %v\n", m.NumGC))
		}

	default:
		{
			chat(s, fmt.Sprintf("%s is not a valid command", command))
		}
	}

}

// OnSocketClose is called when a socket disconnects
func (l *Lobby) OnSocketClose(s *server.Socket) {
	l.subscribers.Remove(s.UID)
}
