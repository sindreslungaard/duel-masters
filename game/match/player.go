package match

import (
	"duel-masters/game/cnd"
	"duel-masters/server"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
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
	UID      string
	Username string
	DeckStr  string
	Color    string
	Player   *Player
	Socket   *server.Socket
	LastPong int64
}

type Spectator struct {
	UID      string
	Username string
	Color    string
	Socket   *server.Socket
	LastPong int64
}

// PlayerAction is the parsed response we retrieve after prompting the client for a selection of cards
type PlayerAction struct {
	Cards  []string `json:"cards"`
	Count  int      `json:"count"`
	Cancel bool     `json:"cancel"`
}

// PlayerActionState is used to store the last sent action message to the client.
// When a player reconnects, if the last action is still active, we send it to the player
type PlayerActionState struct {
	resolved bool
	data     interface{}
}

// NewPlayerReference returns a new player reference
func NewPlayerReference(p *Player, s *server.Socket) *PlayerReference {

	pr := &PlayerReference{
		UID:      s.User.UID,
		Username: s.User.Username,
		Color:    s.User.Color,
		Player:   p,
		Socket:   s,
		LastPong: time.Now().Unix(),
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

	ShieldCounter int
	ShieldMap     map[string]int

	mutex *sync.Mutex

	ActionState PlayerActionState
	Action      chan PlayerAction

	HasChargedMana bool
	CanChargeMana  bool
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
		ActionState:    PlayerActionState{resolved: true},
		Action:         make(chan PlayerAction),
		HasChargedMana: false,
		CanChargeMana:  true,
		Turn:           turn,
		Ready:          false,
		match:          match,

		ShieldCounter: 0,
		ShieldMap:     make(map[string]int),
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

func (p *Player) Cards() []*Card {
	cards := []*Card{}

	concat := func(srcs ...[]*Card) {
		for _, src := range srcs {
			cards = append(cards, src...)
		}
	}

	concat(
		p.deck,
		p.hand,
		p.shieldzone,
		p.manazone,
		p.graveyard,
		p.battlezone,
		p.spellzone,
		p.hiddenzone,
	)

	return cards
}

// CreateDeck initializes a new deck from a list of card ids
func (p *Player) CreateDeck(deck []string) {

	p.mutex.Lock()

	defer p.mutex.Unlock()

	for _, card := range deck {

		c, err := NewCard(p, card)

		if err != nil {
			logrus.Warnf("Failed to create card with id %s", card)
			continue
		}

		p.deck = append(p.deck, c)

	}

}

func (p *Player) CreateRandomDeck() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	cards := GetCardImages()
	cardsAddedToDeck := make(map[string]int)

	recursions := 0
	var addRandomCardToDeck func()
	addRandomCardToDeck = func() {
		recursions++

		// sanity check
		if recursions > 1000 {
			logrus.Error("Max recursions reached in CreateRandomDeck")
			return
		}

		cardToAdd := cards[rand.Intn(len(cards))]

		n, ok := cardsAddedToDeck[cardToAdd]

		if !ok {
			cardsAddedToDeck[cardToAdd] = 1
		}

		// max limit of this card reached
		if n >= 4 {
			addRandomCardToDeck()
			return
		}

		cardsAddedToDeck[cardToAdd] = n + 1

		c, err := NewCard(p, cardToAdd)

		if err != nil {
			logrus.Warnf("Failed to create card with id %s", cardToAdd)
			return
		}

		p.deck = append(p.deck, c)
	}

	for i := 0; i < 40; i++ {
		addRandomCardToDeck()
	}
}

// Empties the players deck
func (p *Player) DestroyDeck() {
	p.deck = nil
}

// SpawnCard creates a new card from an id and adds it to the players' given zone
// used for debugging and development
func (p *Player) SpawnCard(id string, zone string) (*Card, error) {

	p.mutex.Lock()

	defer p.mutex.Unlock()

	c, err := NewCard(p, id)

	if err != nil {
		logrus.Warnf("Failed to create card with id %s", id)
		return nil, err
	}

	c.Zone = zone

	switch zone {
	case HAND:
		p.hand = append(p.hand, c)
	case MANAZONE:
		p.manazone = append(p.manazone, c)
	case SHIELDZONE:
		p.shieldzone = append(p.shieldzone, c)
		p.ShieldCounter = p.ShieldCounter + 1
		p.ShieldMap[c.ID] = p.ShieldCounter
	case DECK:
		p.deck = append(p.deck, c)
	case GRAVEYARD:
		p.graveyard = append(p.graveyard, c)
	case BATTLEZONE:
		p.battlezone = append(p.battlezone, c)
	default:
		logrus.Warnf("Failed to create card with id %s - invalid zone", id)
		return nil, errors.New("invalid zone")
	}

	return c, nil
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

		p.MoveCard(card.ID, DECK, SHIELDZONE, "init_shieldzone")

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
// Returns a reference to the drawn cards array
func (p *Player) DrawCards(n int) []*Card {

	toMove := make([]*Card, 0)

	p.mutex.Lock()

	if len(p.deck) < n {
		n = len(p.deck)
	}

	for i := range n {
		toMove = append(toMove, p.deck[i])
	}

	p.mutex.Unlock()

	for _, card := range toMove {
		p.MoveCard(card.ID, DECK, HAND, "draw")
	}

	if n == 1 {
		p.match.ReportActionInChat(p, fmt.Sprintf("%s drew %v card", p.match.PlayerRef(p).Socket.User.Username, n))
	} else {
		p.match.ReportActionInChat(p, fmt.Sprintf("%s drew %v cards", p.match.PlayerRef(p).Socket.User.Username, n))
	}

	if len(p.deck) <= 0 {
		// deck out
		p.match.End(p.match.Opponent(p), fmt.Sprintf("%s won by deck out!", p.match.Opponent(p).Username()))
	}

	return toMove
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
func (p *Player) MoveCard(cardID string, from string, to string, source string) (*Card, error) {
	c, err := p.GetCard(cardID, from)

	if err != nil {
		return nil, err
	}

	ctx := NewContext(p.match, &MoveCard{
		CardID: cardID,
		From:   from,
		To:     to,
		Source: source,
	})

	p.match.HandleFx(ctx)

	if ctx.Cancelled() {
		return c, nil
	}

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
	ref.Tapped = false

	if to == SHIELDZONE {
		p.ShieldCounter = p.ShieldCounter + 1
		p.ShieldMap[ref.ID] = p.ShieldCounter
	}

	p.mutex.Unlock()

	p.match.HandleFx(NewContext(p.match, &CardMoved{
		CardID:        ref.ID,
		From:          from,
		To:            to,
		Source:        source,
		MatchPlayerID: p.match.getPlayerMatchId(ref.Player),
	}))

	return ref, nil
}

// MoveCard tries to move a card from container a to the front of container b
func (p *Player) MoveCardToFront(cardID string, from string, to string) (*Card, error) {

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

	temp2 := append([]*Card{ref}, *cTo...)

	*cTo = temp2

	ref.Zone = to
	ref.Tapped = false

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

	manaCost := card.ManaCost
	for _, condition := range card.Conditions() {
		if condition.ID == cnd.ReducedCost {
			manaCost -= condition.Val.(int)
			if manaCost < 1 {
				manaCost = 1
			}
		}

		if condition.ID == cnd.IncreasedCost {
			manaCost += condition.Val.(int)
		}
	}

	if manaCost > len(untappedMana) {
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
		HandCount:  len(p.hand),
		Hand:       denormalizeCards(p.hand, false),
		Shieldzone: shields,
		ShieldMap:  p.ShieldMap,
		Manazone:   denormalizeCards(p.manazone, false),
		Graveyard:  denormalizeCards(p.graveyard, false),
		Battlezone: denormalizeCards(p.battlezone, false),
	}

	p.mutex.Unlock()

	return state

}

// denormalizeCards takes an array of *Card and returns an array of server.CardState
// if partial is true, the cards' name and image will not be included
func denormalizeCards(cards []*Card, partial bool) []server.CardState {

	arr := make([]server.CardState, 0)

	for _, card := range cards {

		var flags CardFlags

		mana, err := card.Player.Container(MANAZONE)

		if err == nil && card.Player.CanPlayCard(card, mana) {
			flags |= PlayableFlag
		}

		if card.Tapped {
			flags |= TappedFlag
		}

		if card.HasCondition(cnd.TapAbility) {
			flags |= TapAbilityFlag
		}

		cs := server.CardState{
			CardID:  card.ID,
			ImageID: card.ImageID,
			Name:    card.Name,
			Civ:     card.Civ,
			Flags:   uint8(flags),
		}

		if partial {
			cs.ImageID = "backside"
			cs.Name = ""
			cs.Civ = "water" // blue highlight color when selected in actions

			var f CardFlags
			cs.Flags = uint8(f)
		}

		arr = append(arr, cs)
	}

	return arr

}

// Username returns the username of the player
func (p *Player) Username() string {
	return p.match.PlayerRef(p).Socket.User.Username
}

// Dispose clears out references in the player object
func (p *Player) Dispose() {

	p.mutex.Lock()

	defer p.mutex.Unlock()

	close(p.Action)

	for _, c := range p.deck {
		c.Player = nil
	}

	for _, c := range p.hand {
		c.Player = nil
	}

	for _, c := range p.shieldzone {
		c.Player = nil
	}

	for _, c := range p.manazone {
		c.Player = nil
	}

	for _, c := range p.graveyard {
		c.Player = nil
	}

	for _, c := range p.battlezone {
		c.Player = nil
	}

	for _, c := range p.hiddenzone {
		c.Player = nil
	}

	for _, c := range p.spellzone {
		c.Player = nil
	}

}
