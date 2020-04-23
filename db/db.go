package db

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // postgres db driver
	"github.com/sirupsen/logrus"
)

var conn *gorm.DB

// Connect connects to the database
func Connect(host string, port string, user string, password string, dbName string) {

	connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s", host, port, user, dbName, password)

	database, err := gorm.Open("postgres", connectionString)
	if err != nil {
		logrus.Fatal(err)
	}

	conn = database

	logrus.Info("Connected to database")

	// conn.Debug().AutoMigrate(&Account{}, &Contact{}) //Database migration

}
