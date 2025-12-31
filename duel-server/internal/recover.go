package internal

import (
	"runtime/debug"

	"github.com/sirupsen/logrus"
)

func Recover() {
	if r := recover(); r != nil {
		logrus.Warnf("%v", r)
		debug.PrintStack()
	}
}
