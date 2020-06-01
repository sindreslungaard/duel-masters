package db

// UserSession struct holds the users session information
type UserSession struct {
	Token   string `json:"token"`
	IP      string `json:"ip"`
	Expires int    `json:"expires"`
}

// User struct holds the users information
type User struct {
	UID         string        `json:"uid"`
	Permissions []string      `json:"permissions"`
	Username    string        `json:"username"`
	Password    string        `json:"-"`
	Email       string        `json:"email"`
	Color       string        `json:"color"`
	Sessions    []UserSession `json:"-"`
}

// Deck struct is a player deck
type Deck struct {
	UID      string   `json:"uid"`
	Owner    string   `json:"owner"`
	Name     string   `json:"name"`
	Public   bool     `json:"public"`
	Standard bool     `json:"standard"`
	Cards    []string `json:"cards"`
}
