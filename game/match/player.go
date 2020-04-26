package match

import (
	"duel-masters/server"
	"errors"
	"sync"
)

// Card containers
const (
	DECK       = "deck"
	HAND       = "hand"
	SHIELDZONE = "shieldzone"
	MANAZONE   = "manazone"
	GRAVEYARD  = "graveyard"
	BATTLEZONE = "battlezone"
	HIDDENZONE = "hiddenzone"
)

// Player holds information about the players state in the match
type Player struct {
	socket *server.Socket

	deck       []Card
	hand       []Card
	shieldzone []Card
	manazone   []Card
	graveyard  []Card
	battlezone []Card
	hiddenzone []Card
	mutex      *sync.Mutex

	HasChargedMana bool
}

func (p *Player) container(c string) (*[]Card, error) {

	switch c {
	case DECK:
		return &p.deck, nil
	case HAND:
		return &p.hand, nil
	case SHIELDZONE:
		return &p.shieldzone, nil
	case MANAZONE:
		return &p.manazone, nil
	case GRAVEYARD:
		return &p.graveyard, nil
	case BATTLEZONE:
		return &p.battlezone, nil
	case HIDDENZONE:
		return &p.hiddenzone, nil
	default:
		return nil, errors.New("Invalid container")
	}

}

// HasCard checks if a container has a card
func (p *Player) HasCard(container string, cardID string) bool {

	c, err := p.container(container)

	if err != nil {
		return false
	}

	p.mutex.Lock()

	defer p.mutex.Unlock()

	for _, card := range *c {
		if card.ID == cardID {
			return true
		}
	}

	return false

}

// MoveCard tries to move a card from container a to container b
func (p *Player) MoveCard(cardID string, from string, to string) error {

	cFrom, err := p.container(from)

	if err != nil {
		return err
	}

	if !p.HasCard(from, cardID) {
		return errors.New("Card is not in the specified container")
	}

	cTo, err := p.container(to)

	if err != nil {
		return err
	}

	p.mutex.Lock()

	temp := make([]Card, 0)
	var ref *Card

	for _, card := range *cFrom {
		if card.ID != cardID {
			temp = append(temp, card)
		}
		ref = &card
	}

	*cFrom = temp

	*cTo = append(*cTo, *ref)

	p.mutex.Unlock()

	return nil

}
