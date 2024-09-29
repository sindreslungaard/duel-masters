package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"duel-masters/api"
	"duel-masters/db"
	"duel-masters/db/migrations"
	"duel-masters/game"
	"duel-masters/game/cards"
	"duel-masters/game/match"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		logrus.Warn("Failed to load .env file")
	}

	logrus.SetFormatter(&logrus.TextFormatter{})

	loglevel := strings.ToLower(os.Getenv("log_level"))

	switch loglevel {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		loglevel = "info"
		logrus.SetLevel(logrus.InfoLevel)
	}

	logrus.Infof("Using log level %s", loglevel)

	rand.Seed(time.Now().UnixNano())
}

func main() {
	for _, set := range cards.Sets {
		for uid, ctor := range *set {
			if ctor == nil {
				continue
			}
			match.AddCard(uid, ctor)
		}
	}

	lobby := game.NewLobby()
	go lobby.StartTicker()

	matchSystem := match.NewSystem(lobby.Broadcast)
	go matchSystem.StartTicker()

	game.Matchmaker.Initialize(lobby.Broadcast, matchSystem)

	lobby.SetMatchesFunc(func() []*match.Match { return matchSystem.Matches.Iter() })

	// Setup API
	API := api.New(lobby, matchSystem)

	api.CreateCardCache()

	blockedIps := os.Getenv("blocked_networks")

	if blockedIps != "" {
		iprange, err := api.IPRangeFromExternalSrc(blockedIps)

		if err != nil {
			logrus.Error(err)
		}

		API.SetBlockedIPs(iprange)

		logrus.Infof("Blocked %v networks from using certain API features", iprange.Size())

	}

	migrations.Migrate(db.Connection())

	go checkForAutoRestart(lobby)

	API.Start(os.Getenv("port"))

}

func checkForAutoRestart(lobby *game.Lobby) {

	if os.Getenv("restart_after") == "" {
		logrus.Debug("No autorestart policy found")
		return
	}

	n, err := strconv.Atoi(os.Getenv("restart_after"))

	if err != nil {
		panic(err)
	}

	d := time.Now().Add(time.Second * time.Duration(n))

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	notified := false

	logrus.Info(fmt.Sprintf("Scheduled to shutdown %s", d.Format("2 Jan 2006 15:04")))

	for range ticker.C {

		if time.Now().After(d) {
			logrus.Info("Performing scheduled shutdown")
			os.Exit(0)
		}

		// less than 2 hours until restart and have not yet notified
		if time.Now().Add(2*time.Hour).After(d) && !notified {
			notified = true

			lobby.PinMessage(fmt.Sprintf("Scheduled restart in time:%v", d.Unix()))
		}

	}

}
