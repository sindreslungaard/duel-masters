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
	Sessions    []UserSession `json:"-"`
}
