package fx

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"fmt"
	"math/rand"
)

func GiveTapAbilityToAllies(card *match.Card, ctx *match.Context, alliesFilter func(card *match.Card) bool, tapAbility func(card *match.Card, ctx *match.Context)) {
	// This is added for the case where the card is added to the field. There is another creature
	// that doesn't initially have a tap abbility but should receive one. The change doesn't propagate fast
	// enough to the FE and that creature doesn't get tap ability until another action takes places.
	// This is an ugly workaround.
	FindFilter(
		card.Player,
		match.BATTLEZONE,
		alliesFilter,
	).Map(func(x *match.Card) {
		x.AddUniqueSourceCondition(cnd.TapAbility, tapAbility, card.ID)
	})

	ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
		if card.Zone != match.BATTLEZONE {
			Find(
				card.Player,
				match.BATTLEZONE,
			).Map(func(x *match.Card) {
				x.RemoveConditionBySource(card.ID)
			})

			exit()
			return
		}

		FindFilter(
			card.Player,
			match.BATTLEZONE,
			alliesFilter,
		).Map(func(x *match.Card) {
			x.AddUniqueSourceCondition(cnd.TapAbility, tapAbility, card.ID)
		})
	})
}

func FilterShieldTriggers(ctx *match.Context, filter func(*match.Card) bool) {

	if event, ok := ctx.Event.(*match.ShieldTriggerEvent); ok {
		validCards, invalidCards := FilterCardList(event.Cards, filter)
		event.Cards = validCards
		event.UnplayableCards = append(event.UnplayableCards, invalidCards...)
	}

}

func OpponentDiscardsRandomCard(card *match.Card, ctx *match.Context) {

	hand, err := ctx.Match.Opponent(card.Player).Container(match.HAND)

	if err != nil || len(hand) < 1 {
		return
	}

	discardedCard, err := ctx.Match.Opponent(card.Player).MoveCard(hand[rand.Intn(len(hand))].ID, match.HAND, match.GRAVEYARD, card.ID)
	if err == nil {
		ctx.Match.ReportActionInChat(ctx.Match.Opponent(card.Player), fmt.Sprintf("%s was discarded from %s's hand by %s", discardedCard.Name, discardedCard.Player.Username(), card.Name))
	}

}

func OpponentDiscards2RandomCards(card *match.Card, ctx *match.Context) {
	OpponentDiscardsRandomCard(card, ctx)
	OpponentDiscardsRandomCard(card, ctx)
}

// To be used as part of a card effect, not for initial shuffle
func ShuffleDeck(card *match.Card, ctx *match.Context, forOpponent bool) {
	if !forOpponent {
		card.Player.ShuffleDeck()
		ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s shuffled their deck", card.Player.Username()))
	} else {
		opponent := ctx.Match.Opponent(card.Player)
		opponent.ShuffleDeck()
		ctx.Match.ReportActionInChat(opponent, fmt.Sprintf("%s deck shuffled by %s effect", opponent.Username(), card.Name))
	}

}

func BlockerWhenNoShields(card *match.Card, ctx *match.Context) {
	condition := &match.Condition{ID: cnd.Blocker, Val: true, Src: card.ID}
	HaveSelfConditionsWhenNoShields(card, ctx, []*match.Condition{condition})
}

func HaveSelfConditionsWhenNoShields(card *match.Card, ctx *match.Context, conditions []*match.Condition) {
	ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {

		notInTheBZ := card.Zone != match.BATTLEZONE
		if notInTheBZ || IHaveShields(card) {
			for _, cond := range conditions {
				card.RemoveSpecificConditionBySource(cond.ID, card.ID)
			}
		}

		if notInTheBZ {
			exit()
			return
		}

		if IDontHaveShields(card, ctx2) {
			for _, cond := range conditions {
				if cond.ID == cnd.Blocker {
					ForceBlocker(card, ctx2, card.ID)
				} else {
					card.AddUniqueSourceCondition(cond.ID, cond.Val, card.ID)
				}
			}
		}

	})
}

func RotateShields(card *match.Card, ctx *match.Context, max int) {

	nrShields, err := card.Player.Container(match.SHIELDZONE)
	if err != nil {
		return
	}

	if len(nrShields) < 1 {
		return
	}

	toShield := Select(
		card.Player,
		ctx.Match,
		card.Player,
		match.HAND,
		fmt.Sprintf("%s: You may select up to %d card(s) from your hand and put it into the shield zone", card.Name, max),
		0, max, true,
	)

	cardsMoved := len(toShield)
	if cardsMoved < 1 {
		return
	}

	for _, c := range toShield {
		c.Player.MoveCard(c.ID, match.HAND, match.SHIELDZONE, card.ID)
	}

	toHand := SelectBackside(
		card.Player,
		ctx.Match,
		card.Player,
		match.SHIELDZONE,
		fmt.Sprintf("%s: Select %d of your shields that will be moved to your hand", card.Name, cardsMoved),
		cardsMoved,
		cardsMoved,
		false,
	)

	for _, c := range toHand {
		c.Player.MoveCard(c.ID, match.SHIELDZONE, match.HAND, card.ID)
	}

}

func DestoryOpShield(card *match.Card, ctx *match.Context) {
	opponent := ctx.Match.Opponent(card.Player)

	ctx.Match.BreakShields(SelectBackside(
		card.Player,
		ctx.Match,
		opponent,
		match.SHIELDZONE,
		fmt.Sprintf("%s effect: select shield to break", card.Name),
		1,
		1,
		false,
	), card)

	ctx.Match.ReportActionInChat(ctx.Match.Opponent(card.Player),
		fmt.Sprintf("%s effect broke one of %s's shields", card.Name, opponent.Username()))

}

func OpDiscardsXCards(x int) func(*match.Card, *match.Context) {
	return func(card *match.Card, ctx *match.Context) {

		min := 0
		handCount := ctx.Match.Opponent(card.Player).Denormalized().HandCount

		if x > handCount {
			min = handCount
		} else {
			min = x
		}

		Select(
			ctx.Match.Opponent(card.Player),
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.HAND,
			fmt.Sprintf("%s: Select %d card(s) from your hand that will be sent to your graveyard", card.Name, x),
			min,
			x,
			false,
		).Map(func(x *match.Card) {
			x.Player.MoveCard(x.ID, match.HAND, match.GRAVEYARD, card.ID)
			ctx.Match.ReportActionInChat(ctx.Match.Opponent(card.Player), fmt.Sprintf("%s was moved from %s's hand to his graveyard by %s", x.Name, x.Player.Username(), card.Name))
		})
	}
}
