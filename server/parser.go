package server

import (
	"strconv"

	"github.com/sirupsen/logrus"
)

// Parse handles a websocket message
func Parse(data string) {

	runes := []rune(data)

	header, err := strconv.Atoi(string(runes[0:4]))

	if err != nil {
		return
	}

	logrus.Debugf("Received message with header %s", header)

	switch header {

	// Login
	case 1000:
		return

	// Register
	case 1001:
		return

	default:
		return

	}

}
