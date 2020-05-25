package match

// Condition is used to store turn-specific state to the card such as power amplifiers
type Condition struct {
	id  string
	val interface{}
	src interface{}
}

// Card holds information about a specific card
type Card struct {
	ID      string
	ImageID string
	Player  *Player
	Tapped  bool
	Zone    string

	Name            string
	Power           int
	Civ             string
	Family          string
	ManaCost        int
	ManaRequirement []string
	PowerModifier   func(m *Match, attacking bool) int

	conditions []Condition
	handlers   []HandlerFunc
}

// Use allows different cards to hook into match events
// Can be compared to a typical middleware function
func (c *Card) Use(handlers ...HandlerFunc) {
	c.handlers = append(c.handlers, handlers...)
}

// Conditions returns a slice with the cards conditions
func (c *Card) Conditions() []Condition {
	return c.conditions
}

// AddCondition stores a string to the state of the card that will stay there until removed
func (c *Card) AddCondition(cnd string, val interface{}, src interface{}) {
	c.conditions = append(c.conditions, Condition{cnd, val, src})
}

// HasCondition returns true or false based on if a given string is added to the cards list of conditions
func (c *Card) HasCondition(cnd string) bool {

	for _, condition := range c.conditions {
		if condition.id == cnd {
			return true
		}
	}

	return false

}

// RemoveCondition removes all instances of the given string from the cards conditions
func (c *Card) RemoveCondition(cnd string) {

	tmp := make([]Condition, 0)

	for _, condition := range c.conditions {

		if condition.id != cnd {
			tmp = append(tmp, condition)
		}

	}

	c.conditions = tmp

}

// RemoveConditionBySource removes all instances of conditions with given source
func (c *Card) RemoveConditionBySource(src string) {

	tmp := make([]Condition, 0)

	for _, condition := range c.conditions {

		if condition.src != src {
			tmp = append(tmp, condition)
		}

	}

	c.conditions = tmp

}

// ClearConditions removes all conditions from the card
func (c *Card) ClearConditions() {

	c.conditions = make([]Condition, 0)

}
