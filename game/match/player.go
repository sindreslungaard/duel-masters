package match

import (
	"duel-masters/server"
	"errors"
	"fmt"
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
	deck       []*Card
	hand       []*Card
	shieldzone []*Card
	manazone   []*Card
	graveyard  []*Card
	battlezone []*Card
	hiddenzone []*Card
	spellzone  []*Card

	mutex *sync.Mutex

	Action chan PlayerAction

	HasChargedMana bool
	Turn           byte
	Ready          bool

	match *Match
}

// NewPlayer returns a new player
func NewPlayer(match *Match, turn byte) *Player {

	p := &Player{
		deck:           make([]*Card, 0),
		hand:           make([]*Card, 0),
		shieldzone:     make([]*Card, 0),
		manazone:       make([]*Card, 0),
		graveyard:      make([]*Card, 0),
		battlezone:     make([]*Card, 0),
		spellzone:      make([]*Card, 0),
		hiddenzone:     make([]*Card, 0),
		mutex:          &sync.Mutex{},
		Action:         make(chan PlayerAction),
		HasChargedMana: false,
		Turn:           turn,
		Ready:          false,
		match:          match,
	}

	return p

}

// ContainerRef returns a pointer to one of the player's card zones based on the specified string
func (p *Player) ContainerRef(c string) (*[]*Card, error) {

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
	case SPELLZONE:
		return &p.spellzone, nil
	case HIDDENZONE:
		return &p.hiddenzone, nil
	default:
		return nil, errors.New("Invalid container")
	}

}

// Container returns a copy of one of the player's card zones based on the specified string
func (p *Player) Container(c string) ([]*Card, error) {

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
	case SPELLZONE:
		return p.spellzone, nil
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
			ImageID:         card,
			Player:          p,
			Tapped:          false,
			Zone:            HAND,
			Name:            "undefined_card",
			Civ:             "undefind_civ",
			Family:          "undefined_family",
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

	if n > 1 {
		p.match.Chat("Server", fmt.Sprintf("%s drew %v cards", p.match.PlayerRef(p).Socket.User.Username, n))
	} else {
		p.match.Chat("Server", fmt.Sprintf("%s drew %v card", p.match.PlayerRef(p).Socket.User.Username, n))
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

	for _, card := range c {
		if card.ID == cardID {
			return true
		}
	}

	return false

}

// GetCard returns a pointer to a Card by its ID and container
func (p *Player) GetCard(id string, container string) (*Card, error) {

	c, err := p.Container(container)

	if err != nil {
		return nil, err
	}

	p.mutex.Lock()

	defer p.mutex.Unlock()

	for _, card := range c {
		if card.ID == id {
			return card, nil
		}
	}

	return nil, errors.New("Card was not found")

}

// MoveCard tries to move a card from container a to container b
func (p *Player) MoveCard(cardID string, from string, to string) (*Card, error) {

	cFrom, err := p.ContainerRef(from)

	if err != nil {
		return nil, err
	}

	if !p.HasCard(from, cardID) {
		return nil, errors.New("Card is not in the specified container")
	}

	cTo, err := p.ContainerRef(to)

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

	p.match.HandleFx(NewContext(p.match, &CardMoved{
		CardID: ref.ID,
		From:   from,
		To:     to,
	}))

	return ref, nil

}

// CanPlayCard returns true or false based on if the specified card can be played with the specified mana
func (p *Player) CanPlayCard(card *Card, mana []*Card) bool {

	untappedMana := make([]*Card, 0)
	for _, card := range mana {
		if !card.Tapped {
			untappedMana = append(untappedMana, card)
		}
	}

	if card.ManaCost > len(untappedMana) {
		return false
	}

	for _, manaCard := range untappedMana {
		for _, civ := range card.ManaRequirement {
			if manaCard.Civ == civ {
				return true
			}
		}
	}

	return false

}

// Denormalized returns a server.PlayerState
func (p *Player) Denormalized() *server.PlayerState {

	p.mutex.Lock()

	shields := make([]string, 0)

	for _, card := range p.shieldzone {
		shields = append(shields, card.ID)
	}

	state := &server.PlayerState{
		Deck:       len(p.deck),
		Hand:       denormalizeCards(p.hand),
		Shieldzone: shields,
		Manazone:   denormalizeCards(p.manazone),
		Graveyard:  denormalizeCards(p.graveyard),
		Battlezone: denormalizeCards(p.battlezone),
	}

	p.mutex.Unlock()

	return state

}

// denormalizeCards takes an array of *Card and returns an array of server.CardState
func denormalizeCards(cards []*Card) []server.CardState {

	arr := make([]server.CardState, 0)

	for _, card := range cards {

		canBePlayed := true

		mana, err := card.Player.Container(MANAZONE)

		if err == nil {
			canBePlayed = card.Player.CanPlayCard(card, mana)
		}

		cs := server.CardState{
			CardID:      card.ID,
			ImageID:     card.ImageID,
			Name:        card.Name,
			Civ:         card.Civ,
			Tapped:      card.Tapped,
			CanBePlayed: canBePlayed,
		}
		arr = append(arr, cs)
	}

	return arr

}
