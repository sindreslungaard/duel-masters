package match

type DuelRecord struct {
	UID             string
	Format          string
	Host            string
	HostDeck        string
	Guest           string
	GuestDeck       string
	Started         int64
	Ended           int64
	Turns           int
	Winner          string
	WonByDisconnect bool
}
