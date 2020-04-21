package main

import (
	"log"

	"github.com/sindreslungaard/duel-masters/server"
)

func main() {
	log.Println("Starting..")
	server.Start("80")
}
