package match

import (
	"duel-masters/server"
	"errors"
	"math/rand"
	"sync"
	"time"

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
	deck       []*Card
	hand       []*Card
	shieldzone []*Card
	manazone   []*Card
	graveyard  []*Card
	battlezone []*Card
	hiddenzone []*Card
	mutex      *sync.Mutex

	HasChargedMana bool
	Turn           byte
	Ready          bool
}

// NewPlayer returns a new player
func NewPlayer(turn byte) *Player {

	p := &Player{
		hand:           make([]*Card, 0),
		shieldzone:     make([]*Card, 0),
		manazone:       make([]*Card, 0),
		graveyard:      make([]*Card, 0),
		battlezone:     make([]*Card, 0),
		hiddenzone:     make([]*Card, 0),
		mutex:          &sync.Mutex{},
		HasChargedMana: false,
		Turn:           turn,
		Ready:          false,
	}

	return p

}

func (p *Player) container(c string) ([]*Card, error) {

	switch c {
	case DECK:
		return p.deck, nil
	case HAND:
		return p.hand, nil
	case SHIELDZONE:
		return p.shieldzone, nil
	case MANAZONE:
		return p.manazone, nil
	case GRAVEYARD:
		return p.graveyard, nil
	case BATTLEZONE:
		return p.battlezone, nil
	case HIDDENZONE:
		return p.hiddenzone, nil
	default:
		return nil, errors.New("Invalid container")
	}

}

// CreateDeck initializes a new deck from a list of card ids
func (p *Player) CreateDeck(deck []string) {

	p.mutex.Lock()

	defer p.mutex.Unlock()

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

		cardctor, err := CardCtor(card)

		if err != nil {
			logrus.Warn(err)
			continue
		}

		cardctor(c)

		p.deck = append(p.deck, c)

	}

}

// ShuffleDeck randomizes the order of cards in the players deck
func (p *Player) ShuffleDeck() {

	p.mutex.Lock()

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(p.deck), func(i, j int) { p.deck[i], p.deck[j] = p.deck[j], p.deck[i] })

	p.mutex.Unlock()

}

// InitShieldzone adds 5 cards from the players deck to their shieldzone
func (p *Player) InitShieldzone() {

	cards := p.PeekDeck(5)

	for _, card := range cards {

		p.MoveCard(card.ID, DECK, SHIELDZONE)

	}

}

// PeekDeck returns references to the next n cards in the deck
func (p *Player) PeekDeck(n int) []*Card {

	result := make([]*Card, 0)

	p.mutex.Lock()

	if len(p.deck) < n {
		n = len(p.deck)
	}

	for i := 0; i < n; i++ {
		result = append(result, p.deck[i])
	}

	p.mutex.Unlock()

	return result

}

// DrawCards moves n cards from the players deck to their hand
func (p *Player) DrawCards(n int) {

	toMove := make([]string, 0)

	p.mutex.Lock()

	if len(p.deck) < n {
		n = len(p.deck)
	}

	for i := 0; i < n; i++ {
		toMove = append(toMove, p.deck[i].ID)
	}

	p.mutex.Unlock()

	for _, card := range toMove {
		p.MoveCard(card, DECK, HAND)
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

	for _, card := range c {
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

	temp := make([]*Card, 0)
	var ref *Card

	for _, card := range cFrom {
		if card.ID != cardID {
			temp = append(temp, card)
		}
		ref = card
	}

	cFrom = temp

	cTo = append(cTo, ref)

	p.mutex.Unlock()

	return nil

}
