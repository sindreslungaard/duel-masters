package server

import (
	"strconv"

	"github.com/sirupsen/logrus"
)

// Parse handles a websocket message
func Parse(s *Socket, data []byte) {

	runes := []rune(string(data))

	header, err := strconv.Atoi(string(runes[0:4]))

	if err != nil {
		logrus.Debug("Received message in incorrect format %s", string(data))
		return
	}

	logrus.Debugf("Received message with header %v", header)

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
