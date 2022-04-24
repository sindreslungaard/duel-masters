package match

import (
	"duel-masters/internal"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var system = MatchSystem{
	matches: internal.NewConcurrentDictionary[Match](),
}

type MatchSystem struct {
	sync.RWMutex

	matches internal.ConcurrentDictionary[Match]
}

func StartTicker(s MatchSystem) {

	defer internal.Recover()

	ticker := time.NewTicker(10 * time.Second) // tick every 10 seconds

	defer ticker.Stop()

	for {

		select {
		case <-ticker.C:
			{
				for _, m := range s.matches.Iter() {
					ProcessMatch(m)
				}
			}
		}

	}

}

func ProcessMatch(m *Match) {

	defer internal.Recover()

	// Close the match if it was not started within 10 minutes of creation
	if !m.Started && m.created < time.Now().Unix()-60*10 {
		logrus.Debugf("Closing match %s", m.ID)
		m.Dispose()
	}
}
