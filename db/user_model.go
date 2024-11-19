package db

type UserSession struct {
	Token   string `json:"token"`
	IP      string `json:"ip"`
	Expires int    `json:"expires"`
}

type User struct {
	UID              string        `json:"uid"`
	Permissions      []string      `json:"permissions"`
	Username         string        `json:"username"`
	Password         string        `json:"-"`
	Email            string        `json:"email"`
	Color            string        `json:"color"`
	Playmat          string        `json:"playmat"`
	Sessions         []UserSession `json:"-"`
	MutedUsers       []string      `json:"muted_users"`
	Chatblocked      bool          `json:"-" bson:"chat_blocked"`
	TotalGamesPlayed int           `json:"total_games_played" bson:"total_games_played"`
	GamesWon         int           `json:"games_won" bson:"games_won"`
	GamesLost        int           `json:"games_lost" bson:"games_lost"`
	WinRate          float64       `json:"win_rate" bson:"win_rate"`
}

var Users = conn().Collection("users")
