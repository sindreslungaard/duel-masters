package internal

import (
	"sync"
	"time"
)

type entry struct {
	Attempts     int
	FirstAttempt time.Time
}

var store map[string]*entry = map[string]*entry{}
var mutex = &sync.Mutex{}

func RateLimited(identifier string, requestsAllowed int, perMs int64) (rateLimited bool) {
	mutex.Lock()
	defer mutex.Unlock()

	e, ok := store[identifier]

	if !ok {
		store[identifier] = &entry{
			Attempts:     1,
			FirstAttempt: time.Now(),
		}
		return false
	}

	if e.FirstAttempt.Add(time.Duration(perMs) * time.Millisecond).Before(time.Now()) {
		e.FirstAttempt = time.Now()
		e.Attempts = 0
	}

	if e.Attempts >= requestsAllowed {
		return true
	}

	e.Attempts += 1
	return false

}
