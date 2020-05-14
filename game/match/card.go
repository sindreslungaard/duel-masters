package match

// Card holds information about a specific card
type Card struct {
	ID      string
	ImageID string
	Player  *Player
	Tapped  bool
	Zone    string

	Name            string
	Civ             string
	Family          string
	ManaCost        int
	ManaRequirement []string

	conditions []string
	handlers   []HandlerFunc
}

// Use allows different cards to hook into match events
// Can be compared to a typical middleware function
func (c *Card) Use(handlers ...HandlerFunc) {
	c.handlers = append(c.handlers, handlers...)
}

// AddCondition stores a string to the state of the card that will stay there until removed
func (c *Card) AddCondition(cnd string) {
	c.conditions = append(c.conditions, cnd)
}

// HasCondition returns true or false based on if a given string is added to the cards list of conditions
func (c *Card) HasCondition(cnd string) bool {

	for _, condition := range c.conditions {
		if condition == cnd {
			return true
		}
	}

	return false

}

// RemoveCondition removes all instances of the given string from the cards conditions
func (c *Card) RemoveCondition(cnd string) {

	tmp := make([]string, 0)

	for _, condition := range c.conditions {

		if condition != cnd {
			tmp = append(tmp, condition)
		}

	}

	c.conditions = tmp

}
