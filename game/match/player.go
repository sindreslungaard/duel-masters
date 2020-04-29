package match

import (
	"duel-masters/game/cards"
	"duel-masters/server"
	"errors"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/ventu-io/go-shortid"
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

// PlayerReference ties a player to a websocket connection
type PlayerReference struct {
	Player *Player
	Socket *server.Socket
}

// NewPlayerReference returns a new player reference
func NewPlayerReference(p *Player, s *server.Socket) *PlayerReference {

	pr := &PlayerReference{
		Player: p,
		Socket: s,
	}

	return pr

}

// Player holds information about the players state in the match
type Player struct {
	deck       []Card
	hand       []Card
	shieldzone []Card
	manazone   []Card
	graveyard  []Card
	battlezone []Card
	hiddenzone []Card
	mutex      *sync.Mutex

	HasChargedMana bool
	Turn           byte
	Ready          bool
}

// NewPlayer returns a new player
func NewPlayer(turn byte) *Player {

	p := &Player{
		hand:           make([]Card, 0),
		shieldzone:     make([]Card, 0),
		manazone:       make([]Card, 0),
		graveyard:      make([]Card, 0),
		battlezone:     make([]Card, 0),
		hiddenzone:     make([]Card, 0),
		mutex:          &sync.Mutex{},
		HasChargedMana: false,
		Turn:           turn,
		Ready:          false,
	}

	return p

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

// CreateDeck initializes a new deck from a list of card ids
func (p *Player) CreateDeck(deck []string) {

	for _, card := range deck {

		id, err := shortid.Generate()

		if err != nil {
			logrus.Debug("Failed to generate id for card")
			continue
		}

		c := &Card{
			ID:              id,
			Player:          p,
			Tapped:          false,
			Name:            "",
			Civ:             "",
			Family:          "",
			ManaCost:        1,
			ManaRequirement: make([]string, 0),
		}

		cardctor := cards.Cards[card]

		if cardctor == nil {
			logrus.Warnf("Failed to construct card with uid %s", card)
			continue
		}

		cardctor(c)

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
