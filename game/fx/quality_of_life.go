package fx

import (
	"duel-masters/game/family"
	"duel-masters/game/match"
)

// CardCollection is a slice of cards with a mapping function
type CardCollection []*match.Card

// Map iterates through cards in the collection and executes the function on them
func (c CardCollection) Map(h func(x *match.Card)) CardCollection {
	for _, card := range c {
		h(card)
	}

	return c
}

func (c CardCollection) Or(h func()) {
	if len(c) > 0 {
		return
	}

	h()
}

func FilterCardList(cards []*match.Card, filter func(*match.Card) bool) (CardCollection, CardCollection) {
	accepted := make([]*match.Card, 0)
	rejected := make([]*match.Card, 0)

	for _, mCard := range cards {
		if filter(mCard) {
			accepted = append(accepted, mCard)
		} else {
			rejected = append(rejected, mCard)
		}
	}

	return accepted, rejected
}

// FindFilter returns a CardCollection matching the filter
func FindFilter(p *match.Player, collection string, h func(card *match.Card) bool) CardCollection {

	result := CardCollection{}

	container, err := p.Container(collection)

	if err != nil {
		return result
	}

	for _, card := range container {
		if h(card) {
			result = append(result, card)
		}
	}

	return result

}

// Find returns a CardCollection for the specified container
func Find(p *match.Player, collection string) CardCollection {
	return FindFilter(p, collection, func(x *match.Card) bool { return true })
}

// When performs the specified function if the test is successful
func When(test func(*match.Card, *match.Context) bool, h func(*match.Card, *match.Context)) func(*match.Card, *match.Context) {

	return func(card *match.Card, ctx *match.Context) {
		if test(card, ctx) {
			h(card, ctx)
		}
	}

}

// Select prompts the user to select n cards from the specified container
func Select(p *match.Player, m *match.Match, containerOwner *match.Player, containerName string, text string, min int, max int, cancellable bool) CardCollection {
	return SelectFilterSelectablesOnly(p, m, containerOwner, containerName, text, min, max, cancellable, func(x *match.Card) bool { return true })
}

// SelectCount prompts to user to select a number in an interval
func SelectCount(p *match.Player, m *match.Match, text string, min int, max int) int {
	result := 0

	m.NewCountAction(p, text, min, max)

	defer m.CloseAction(p)

	if !m.IsPlayerTurn(p) {
		m.Wait(m.Opponent(p), "Waiting for your opponent to make an action")
		defer m.EndWait(m.Opponent(p))
	}

	for {

		action := <-p.Action

		if action.Count < min || action.Count > max {
			m.ActionWarning(p, "The amount selected does not match the requirements")
			continue
		}

		result = action.Count

		break

	}

	return result
}

func BinaryQuestion(p *match.Player, m *match.Match, text string) bool {
	m.NewQuestionAction(p, text)

	defer m.CloseAction(p)

	if !m.IsPlayerTurn(p) {
		m.Wait(m.Opponent(p), "Waiting for your opponent to make an action")
		defer m.EndWait(m.Opponent(p))
	}

	action := <-p.Action

	return !action.Cancel
}

// SelectFilterSelectablesOnly prompts the user to select n cards from the specified container that matches the given filter
//
// Deprecated: New cards should use `fx.SelectFilter`
func SelectFilterSelectablesOnly(p *match.Player, m *match.Match, containerOwner *match.Player, containerName string, text string, min int, max int, cancellable bool, filter func(*match.Card) bool) CardCollection {
	return SelectFilter(p, m, containerOwner, containerName, text, min, max, cancellable, filter, false)
}

// SelectFilter prompts the user to select n cards from the specified container that matches the given filter.
// It also allows to show all the other cards from the container that are unselectable
func SelectFilter(p *match.Player, m *match.Match, containerOwner *match.Player, containerName string, text string, min int, max int, cancellable bool, filter func(*match.Card) bool, showUnselectables bool) CardCollection {

	result := make([]*match.Card, 0)

	cards, err := containerOwner.Container(containerName)

	if err != nil || len(cards) < 1 {
		return result
	}

	filtered, unselectables := FilterCardList(cards, filter)
	if !showUnselectables {
		unselectables = nil
	}

	if len(filtered) < 1 {
		return result
	}

	m.NewActionFullList(p, filtered, min, max, text, cancellable, unselectables)

	defer m.CloseAction(p)

	for {

		action := <-p.Action

		if cancellable && action.Cancel {
			break
		}

		if len(action.Cards) < min || len(action.Cards) > max || !match.AssertCardsIn(filtered, action.Cards...) {
			m.ActionWarning(p, "The cards you selected does not meet the requirements")
			continue
		}

		for _, c := range action.Cards {

			selectedCard, err := containerOwner.GetCard(c, containerName)

			if err != nil {
				continue
			}

			result = append(result, selectedCard)

		}

		break

	}

	return result

}

// SelectMultipart prompts the user to select n cards from the specified list of cards
func SelectMultipart(p *match.Player, m *match.Match, cards map[string][]*match.Card, text string, min int, max int, cancellable bool) CardCollection {
	return selectMultipartBase(p, m, cards, text, min, max, cancellable, false)
}

// SelectMultipart prompts the user to select n cards from the specified list of cards
func SelectMultipartBackside(p *match.Player, m *match.Match, cards map[string][]*match.Card, text string, min int, max int, cancellable bool) CardCollection {
	return selectMultipartBase(p, m, cards, text, min, max, cancellable, true)
}

func selectMultipartBase(p *match.Player, m *match.Match, cards map[string][]*match.Card, text string, min int, max int, cancellable bool, backsideOnly bool) CardCollection {

	result := make([]*match.Card, 0)

	notEmpty := false

	for _, cardList := range cards {
		if len(cardList) > 0 {
			notEmpty = true
		}
	}

	if !notEmpty {
		return result
	}

	if backsideOnly {
		m.NewMultipartActionBackside(p, cards, min, max, text, cancellable)
	} else {
		m.NewMultipartAction(p, cards, min, max, text, cancellable)
	}

	if !m.IsPlayerTurn(p) {
		m.Wait(m.Opponent(p), "Waiting for your opponent to make an action")
		defer m.EndWait(m.Opponent(p))
	}

	defer m.CloseAction(p)

	for {

		action := <-p.Action

		if cancellable && action.Cancel {
			break
		}

		if len(action.Cards) < min || len(action.Cards) > max {
			m.ActionWarning(p, "The cards you selected does not meet the requirements")
			continue
		}

		for _, vid := range action.Cards {

			for _, cardList := range cards {
				for _, card := range cardList {
					if card.ID == vid {
						result = append(result, card)
					}
				}
			}

		}

		break

	}

	return result

}

// SelectBackside prompts the user to select n cards from the specified container
func SelectBackside(p *match.Player, m *match.Match, containerOwner *match.Player, containerName string, text string, min int, max int, cancellable bool) CardCollection {
	return SelectBacksideFilter(p, m, containerOwner, containerName, text, min, max, cancellable, func(x *match.Card) bool { return true })
}

// SelectBacksideFilter prompts the user to select n cards from the specified container that matches the given filter
func SelectBacksideFilter(p *match.Player, m *match.Match, containerOwner *match.Player, containerName string, text string, min int, max int, cancellable bool, filter func(*match.Card) bool) CardCollection {

	result := make([]*match.Card, 0)

	cards, err := containerOwner.Container(containerName)

	if err != nil || len(cards) < 1 {
		return result
	}

	filtered := make([]*match.Card, 0)

	for _, mCard := range cards {
		if filter(mCard) {
			filtered = append(filtered, mCard)
		}
	}

	if len(filtered) < 1 {
		return result
	}

	m.NewBacksideAction(p, filtered, min, max, text, cancellable)

	defer m.CloseAction(p)

	for {

		action := <-p.Action

		if cancellable && action.Cancel {
			break
		}

		if len(action.Cards) < min || len(action.Cards) > max || !match.AssertCardsIn(filtered, action.Cards...) {
			m.ActionWarning(p, "The cards you selected does not meet the requirements")
			continue
		}

		for _, c := range action.Cards {

			selectedCard, err := containerOwner.GetCard(c, containerName)

			if err != nil {
				continue
			}

			result = append(result, selectedCard)

		}

		break

	}

	return result

}

// Hooks below:
// hooks are shorthands for checking if the context matches a certain condition

// Summoned returns true if the card was summoned
//
// Does not activate if the card was under an Evolution card and becomes visible again.
func Summoned(card *match.Card, ctx *match.Context) bool {
	event, ok := ctx.Event.(*match.CardMoved)
	if !ok {
		return false
	}

	return CreatureSummoned(card, ctx) && event.CardID == card.ID
}

// InTheBattlezone returns true if the card arrives in the battlezone.
//
// Similar to summon but activates also if the card was under an Evolution card and becomes visible again.
// Used for cards that have continuous effects while in the battlezone.
func InTheBattlezone(card *match.Card, ctx *match.Context) bool {
	if event, ok := ctx.Event.(*match.CardMoved); ok {

		if event.CardID == card.ID && event.To == match.BATTLEZONE {
			return true
		}

	}

	return false
}

// SpellCast returns true if the spell was cast
func SpellCast(card *match.Card, ctx *match.Context) bool {

	if event, ok := ctx.Event.(*match.SpellCast); ok {

		if event.CardID == card.ID {
			return true
		}

	}

	return false

}

// Attacking returns true if the card is attacking a player or creature
func Attacking(card *match.Card, ctx *match.Context) bool {

	if event, ok := ctx.Event.(*match.AttackCreature); ok {
		if event.CardID == card.ID {
			return true
		}
	}

	if event, ok := ctx.Event.(*match.AttackPlayer); ok {
		if event.CardID == card.ID {
			return true
		}
	}

	return false

}

// AttackConfirmed returns true if the card is attacking and it cannot be cancelled
func AttackConfirmed(card *match.Card, ctx *match.Context) bool {

	if event, ok := ctx.Event.(*match.AttackConfirmed); ok {
		if event.CardID == card.ID {
			return true
		}
	}

	return false
}

// AttackingPlayer returns true if the card is attacking a player
func AttackingPlayer(card *match.Card, ctx *match.Context) bool {

	if event, ok := ctx.Event.(*match.AttackPlayer); ok {
		if event.CardID == card.ID {
			return true
		}
	}

	return false

}

// AttackingCreature returns true if the card is attacking a Creature
func AttackingCreature(card *match.Card, ctx *match.Context) bool {

	if event, ok := ctx.Event.(*match.AttackCreature); ok {
		if event.CardID == card.ID {
			return true
		}
	}

	return false

}

// Destroyed returns true if the card was destroyed
func Destroyed(card *match.Card, ctx *match.Context) bool {

	if event, ok := ctx.Event.(*match.CardMoved); ok && event.From == match.BATTLEZONE && event.To == match.GRAVEYARD {
		if event.CardID == card.ID {
			return true
		}
	}

	return false

}

// EndOfTurn returns true if the turn is ending, pre end of turn triggers
func EndOfTurn(card *match.Card, ctx *match.Context) bool {
	_, ok := ctx.Event.(*match.EndStep)

	if !ok {
		return false
	}

	return ok
}

// EndOfTurn returns true if the turn is ending, pre end of turn triggers
func EndOfMyTurn(card *match.Card, ctx *match.Context) bool {
	_, ok := ctx.Event.(*match.EndStep)

	if !ok {
		return false
	}

	if ctx.Match.IsPlayerTurn(card.Player) {
		return true
	}

	return false
}

// BreakShield returns true if a shield is about to be broken
func BreakShield(card *match.Card, ctx *match.Context) bool {

	if card.Zone != match.BATTLEZONE {
		return false
	}

	_, ok := ctx.Event.(*match.BreakShieldEvent)
	return ok

}

// ShieldBroken returns true if a shield has been broken
func ShieldBroken(card *match.Card, ctx *match.Context) bool {

	_, ok := ctx.Event.(*match.BrokenShieldEvent)
	return ok

}

// CreatureSummoned returns true if a card was summoned
//
// Does not activate if a card that was under an Evolution card becomes visible again.
func CreatureSummoned(card *match.Card, ctx *match.Context) bool {
	if card.Zone != match.BATTLEZONE {
		return false
	}

	if event, ok := ctx.Event.(*match.CardMoved); ok {

		if event.To == match.BATTLEZONE && event.From != match.HIDDENZONE {
			return true
		}

	}

	return false
}

// MySurvivorSummoned returns true if one of my survivors is summoned
//
// Does not activate if a card that was under an Evolution card becomes visible again.
func MySurvivorSummoned(card *match.Card, ctx *match.Context) bool {

	if !CreatureSummoned(card, ctx) {
		return false
	}

	event, ok := ctx.Event.(*match.CardMoved)
	if !ok {
		return false
	}

	creature, err := card.Player.GetCard(event.CardID, match.BATTLEZONE)
	if err != nil {
		return false
	}

	if !creature.HasFamily(family.Survivor) {
		return false
	}
	return true
}

// AnotherCreatureSummoned returns true if another card was summoned
//
// Does not activate if this current card is summoned.
// Does not activate if a card that was under an Evolution card becomes visible again.
func AnotherCreatureSummoned(card *match.Card, ctx *match.Context) bool {
	event, ok := ctx.Event.(*match.CardMoved)
	if !ok {
		return false
	}

	return CreatureSummoned(card, ctx) && event.CardID != card.ID
}

func AnotherCreatureDestroyed(card *match.Card, ctx *match.Context) bool {
	if card.Zone != match.BATTLEZONE {
		return false
	}

	if event, ok := ctx.Event.(*match.CardMoved); ok &&
		event.From == match.BATTLEZONE &&
		event.To == match.GRAVEYARD &&
		event.CardID != card.ID {
		return true
	}

	return false
}
