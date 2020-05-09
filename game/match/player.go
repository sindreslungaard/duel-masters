package match

import (
	"duel-masters/server"
	"errors"
	"math/rand"
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
	SPELLZONE  = "spellzone"
	HIDDENZONE = "hiddenzone"
)

// PlayerReference ties a player to a websocket connection
type PlayerReference struct {
	Player *Player
	Socket *server.Socket
}

// PlayerAction is the parsed response we retrieve after prompting the client for a selection of cards
type PlayerAction struct {
	Cards  []string `json:"cards"`
	Cancel bool     `json:"cancel"`
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
	Deck       []*Card
	Hand       []*Card
	Shieldzone []*Card
	Manazone   []*Card
	Graveyard  []*Card
	Battlezone []*Card
	Hiddenzone []*Card
	Spellzone  []*Card

	mutex *sync.Mutex

	Action chan PlayerAction

	HasChargedMana bool
	Turn           byte
	Ready          bool
}

// NewPlayer returns a new player
func NewPlayer(turn byte) *Player {

	p := &Player{
		Deck:           make([]*Card, 0),
		Hand:           make([]*Card, 0),
		Shieldzone:     make([]*Card, 0),
		Manazone:       make([]*Card, 0),
		Graveyard:      make([]*Card, 0),
		Battlezone:     make([]*Card, 0),
		Spellzone:      make([]*Card, 0),
		Hiddenzone:     make([]*Card, 0),
		mutex:          &sync.Mutex{},
		Action:         make(chan PlayerAction),
		HasChargedMana: false,
		Turn:           turn,
		Ready:          false,
	}

	return p

}

// Container returns one of the player's card zones based on the specified string
func (p *Player) Container(c string) (*[]*Card, error) {

	switch c {
	case DECK:
		return &p.Deck, nil
	case HAND:
		return &p.Hand, nil
	case SHIELDZONE:
		return &p.Shieldzone, nil
	case MANAZONE:
		return &p.Manazone, nil
	case GRAVEYARD:
		return &p.Graveyard, nil
	case BATTLEZONE:
		return &p.Battlezone, nil
	case SPELLZONE:
		return &p.Spellzone, nil
	case HIDDENZONE:
		return &p.Hiddenzone, nil
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
			ImageID:         card,
			Player:          p,
			Tapped:          false,
			Zone:            HAND,
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

		p.Deck = append(p.Deck, c)

	}

}

// ShuffleDeck randomizes the order of cards in the players deck
func (p *Player) ShuffleDeck() {

	p.mutex.Lock()

	rand.Shuffle(len(p.Deck), func(i, j int) { p.Deck[i], p.Deck[j] = p.Deck[j], p.Deck[i] })

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

	if len(p.Deck) < n {
		n = len(p.Deck)
	}

	for i := 0; i < n; i++ {
		result = append(result, p.Deck[i])
	}

	p.mutex.Unlock()

	return result

}

// DrawCards moves n cards from the players deck to their hand
func (p *Player) DrawCards(n int) {

	toMove := make([]string, 0)

	p.mutex.Lock()

	if len(p.Deck) < n {
		n = len(p.Deck)
	}

	for i := 0; i < n; i++ {
		toMove = append(toMove, p.Deck[i].ID)
	}

	p.mutex.Unlock()

	for _, card := range toMove {
		p.MoveCard(card, DECK, HAND)
	}

}

// HasCard checks if a container has a card
func (p *Player) HasCard(container string, cardID string) bool {

	c, err := p.Container(container)

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
func (p *Player) MoveCard(cardID string, from string, to string) (*Card, error) {

	cFrom, err := p.Container(from)

	if err != nil {
		return nil, err
	}

	if !p.HasCard(from, cardID) {
		return nil, errors.New("Card is not in the specified container")
	}

	cTo, err := p.Container(to)

	if err != nil {
		return nil, err
	}

	p.mutex.Lock()

	temp := make([]*Card, 0)
	var ref *Card

	for _, card := range *cFrom {
		if card.ID != cardID {
			temp = append(temp, card)
		} else {
			ref = card
		}
	}

	*cFrom = temp

	temp2 := append(*cTo, ref)

	*cTo = temp2

	ref.Zone = to

	p.mutex.Unlock()

	return ref, nil

}

// Denormalized returns a server.PlayerState
func (p *Player) Denormalized() *server.PlayerState {

	p.mutex.Lock()

	shields := make([]string, 0)

	for _, card := range p.Shieldzone {
		shields = append(shields, card.ID)
	}

	state := &server.PlayerState{
		Deck:       len(p.Deck),
		Hand:       denormalizeCards(p.Hand),
		Shieldzone: shields,
		Manazone:   denormalizeCards(p.Manazone),
		Graveyard:  denormalizeCards(p.Graveyard),
		Battlezone: denormalizeCards(p.Battlezone),
	}

	p.mutex.Unlock()

	return state

}

// denormalizeCards takes an array of *Card and returns an array of server.CardState
func denormalizeCards(cards []*Card) []server.CardState {

	arr := make([]server.CardState, 0)

	for _, card := range cards {
		cs := server.CardState{
			CardID:      card.ID,
			ImageID:     card.ImageID,
			Name:        card.Name,
			Civ:         card.Civ,
			Tapped:      card.Tapped,
			CanBePlayed: true,
		}
		arr = append(arr, cs)
	}

	return arr

}
