package server

type User struct {
	UID         string
	Permissions []string
	Username    string
	Color       string
	Chatblocked bool
}
