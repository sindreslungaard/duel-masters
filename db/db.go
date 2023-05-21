package db

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var conn *gorm.DB

// Connect connects to the database
func Connect(user string, password string, host string, port string, database string) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	conn = db

	logrus.Info("Connected to database")

}

func Migrate() {

	if conn == nil {
		panic("Can't migrate database before connecting")
	}

	conn.AutoMigrate(
		&User{},
		&UserSession{},
		&UserPermission{},
		&Deck{},
		&DeckCards{},
		&Card{},
	)
}

// GetUserForToken returns a user from the authorization header or returns an error
func GetUserForToken(token string) (User, error) {

	var session UserSession
	tx := conn.First(&session, "token = ?", token)

	if tx.Error != nil {
		return User{}, tx.Error
	}

	var user User
	tx = conn.First(&user, session.UserID)

	if tx.Error != nil {
		return user, tx.Error
	}

	return user, nil

}
