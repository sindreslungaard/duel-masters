package match

// Match struct
type Match struct {
	MatchName    string
	Player1ID    string
	Player2ID    string
	InviteID     string
	PlayerTurnID string
}

// New returns a new match object
func New(matchName string, player1 string) *Match {

	m := &Match{
		MatchName: matchName,
		Player1ID: player1,
		Player2ID: "",
		//InviteID:  uuid.New().String(),
	}

	return m

}
