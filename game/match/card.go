package match

// Card holds information about a specific card
type Card struct {
	PlayerID string
	Tapped   bool
}

// Use allows different cards to hook into match events
// Can be compared to a typical middleware function
func (card *Card) Use(...HandlerFunc) {

}
