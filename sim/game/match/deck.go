package match

import (
	"slices"
	"sync"
)

var cardUIDs []string
var lookupMutex sync.RWMutex

func InitDecks() {
	lookupMutex.Lock()
	defer lookupMutex.Unlock()

	cardUIDs = make([]string, 0, len(ctors))
	for uid := range ctors {
		cardUIDs = append(cardUIDs, uid)
	}

	slices.Sort(cardUIDs)
}

func GetCardImages() []string {
	lookupMutex.RLock()
	defer lookupMutex.RUnlock()

	return slices.Clone(cardUIDs)
}
