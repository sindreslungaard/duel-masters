package match

import (
	"duel-masters/game/cnd"
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/ventu-io/go-shortid"
)

// Condition is used to store turn-specific state to the card such as power amplifiers
type Condition struct {
	ID  string
	Val interface{}
	Src interface{}
}

// Bitmask
type CardFlags uint8

const (
	TappedFlag       CardFlags = 1 << iota // 1
	PlayableFlag                           // 2
	TapAbilityFlag                         // 4
	ShieldFaceUpFlag                       // 8
)

// Card holds information about a specific card
type Card struct {
	ID           string
	ImageID      string
	Player       *Player
	Tapped       bool
	ShieldFaceUp bool
	Zone         string

	Name            string
	Power           int
	Civ             string
	Family          []string
	ManaCost        int
	ManaRequirement []string
	PowerModifier   func(m *Match, attacking bool) int
	TapAbility      func(card *Card, ctx *Context)

	attachedCards []*Card
	conditions    []Condition
	handlers      []HandlerFunc

	localData map[string]string
}

// NewCard returns a new, initialized card
func NewCard(p *Player, image string) (*Card, error) {

	id, err := shortid.Generate()

	if err != nil {
		logrus.Debug("Failed to generate id for card")
		return nil, err
	}

	c := &Card{
		ID:              id,
		ImageID:         image,
		Player:          p,
		Tapped:          false,
		ShieldFaceUp:    false,
		Zone:            DECK,
		Name:            "undefined_card",
		Power:           0,
		Civ:             "undefind_civ",
		Family:          []string{"undefined_family"},
		ManaCost:        1,
		ManaRequirement: make([]string, 0),
		PowerModifier:   func(m *Match, attacking bool) int { return 0 },
		localData:       map[string]string{},
	}

	cardctor, err := CardCtor(image)

	if err != nil {
		logrus.Warn(err)
		return nil, err
	}

	cardctor(c)

	return c, nil

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

// AddUniqueSourceCondition adds a condition only if the specified cnd and source is not already added
func (c *Card) AddUniqueSourceCondition(cnd string, val interface{}, src interface{}) {

	for _, condition := range c.conditions {
		if condition.ID == cnd && condition.Src == src {
			return
		}
	}

	c.conditions = append(c.conditions, Condition{cnd, val, src})
}

// HasCondition returns true or false based on if a given string is added to the cards list of conditions
func (c *Card) HasCondition(cnd string) bool {

	for _, condition := range c.conditions {
		if condition.ID == cnd {
			return true
		}
	}

	return false

}

// GetCondition returns the first condition with the ID given from the card's conditions
func (c *Card) GetCondition(cnd string) (Condition, error) {

	for _, condition := range c.conditions {
		if condition.ID == cnd {
			return condition, nil
		}
	}

	return Condition{}, errors.New("Condition was not found")

}

// RemoveCondition removes all instances of the given string from the cards conditions
func (c *Card) RemoveCondition(cnd string) {

	tmp := make([]Condition, 0)

	for _, condition := range c.conditions {

		if condition.ID != cnd {
			tmp = append(tmp, condition)
		}

	}

	c.conditions = tmp

}

// RemoveConditionBySource removes all instances of conditions with given source
func (c *Card) RemoveConditionBySource(src string) {

	tmp := make([]Condition, 0)

	for _, condition := range c.conditions {

		if condition.Src != src {
			tmp = append(tmp, condition)
		}

	}

	c.conditions = tmp

}

// RemoveConditionBySource removes all instances of specific conditions with given source
func (c *Card) RemoveSpecificConditionBySource(cnd string, src string) {

	tmp := make([]Condition, 0)

	for _, condition := range c.conditions {

		if condition.Src != src || condition.ID != cnd {
			tmp = append(tmp, condition)
		}

	}

	c.conditions = tmp

}

// ClearConditions removes all conditions from the card
func (c *Card) ClearConditions() {

	c.conditions = make([]Condition, 0)

}

// Attach adds a *Card to the card's list of attached cards
func (c *Card) Attach(toAttach ...*Card) {
	c.attachedCards = append(c.attachedCards, toAttach...)
}

// Detach removes a *Card from the card's list of attached cards
func (c *Card) Detach(toDetach *Card) {

	tmp := make([]*Card, 0)

	for _, card := range c.attachedCards {

		if card.ID != toDetach.ID {
			tmp = append(tmp, card)
		}

	}

	c.attachedCards = tmp

}

// Attachments returns a copy of the card's attached cards
func (c *Card) Attachments() []*Card {

	return c.attachedCards

}

// ClearAttachments removes all attached cards
func (c *Card) ClearAttachments() {
	c.attachedCards = make([]*Card, 0)
}

// HasFamily returns true if the card has the input family, false otherwise
func (c *Card) HasFamily(family string) bool {

	families := c.Family

	for _, condition := range c.conditions {
		if condition.ID == cnd.AddFamily {
			if f, ok := condition.Val.(string); ok {
				families = append(families, f)
			}
		}
	}

	for _, f := range families {
		if f == family {
			return true
		}
	}

	return false

}

// SharesAFamily returns true if the card has at least one of the input families, false otherwise
func (c *Card) SharesAFamily(families []string) bool {
	for _, f := range c.Family {
		for _, toCheck := range families {
			if f == toCheck {
				return true
			}
		}
	}

	return false
}

// SharesAllFamilies returns true if the card has all of the input families
func (c *Card) SharesAllFamilies(families []string) bool {
	ok := true

	for _, f := range c.Family {
		for _, toCheck := range families {
			if f != toCheck {
				ok = false
			}
		}
	}

	return ok
}

// StealthActive returns true if the card has stealth
// and the opponent has a card in manazone that matches the stealth civ
func (c *Card) StealthActive(ctx *Context) bool {
	if !c.HasCondition(cnd.Stealth) {
		return false
	}

	for _, cond := range c.Conditions() {
		if cond.ID != cnd.Stealth {
			continue
		}
		if ContainerHas(
			ctx.Match.Opponent(c.Player),
			MANAZONE,
			func(x *Card) bool { return x.Civ == cond.Val },
		) {
			return true
		}
	}

	return false
}

// IsBlockable returns true if the card can be blocked by the opponent's blockers
func (c *Card) IsBlockable(ctx *Context) bool {
	return !c.HasCondition(cnd.CantBeBlocked) && !c.StealthActive(ctx)
}
