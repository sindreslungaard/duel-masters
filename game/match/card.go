package match

// Card holds information about a specific card
type Card struct {
	ID     string
	Player *Player
	Tapped bool

	Name            string
	Civ             string
	Family          string
	ManaCost        int
	ManaRequirement []string
}

// Use allows different cards to hook into match events
// Can be compared to a typical middleware function
func (card *Card) Use(...HandlerFunc) {

}
