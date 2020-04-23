package server

import (
	"strconv"
)

// Parse handles a websocket message
func Parse(data string) {

	runes := []rune(data)

	header, err := strconv.Atoi(string(runes[0:4]))

	if err != nil {
		return
	}

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
