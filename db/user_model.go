package db

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID          uint   `json:"uid"`
	Username    string `json:"username"`
	Password    string `json:"-"`
	Email       string `json:"email"`
	Color       string `json:"color"`
	Playmat     string `json:"playmat"`
	Chatblocked bool   `json:"-"`
}

type UserSession struct {
	gorm.Model
	ID      uint
	UserID  uint
	Token   string
	IP      string
	Created uint
	Expires uint
}

type UserPermission struct {
	gorm.Model
	ID         uint
	UserID     uint
	Permission string
}
