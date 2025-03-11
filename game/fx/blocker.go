package fx

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/match"
)

// Blocker adds the Blocker condition at the Untap step and
// Adds the card to the blockers list at the SelectBlockers step if eligible
func Blocker(card *match.Card, ctx *match.Context) {

	addBlockerCondition(card, ctx)

	defaultBlocker(card, ctx)

}

func FireAndNatureBlocker(card *match.Card, ctx *match.Context) {

	addBlockerCondition(card, ctx)

	conditionalBlocker(card, ctx, func(attacker *match.Card) bool {
		return attacker.Civ == civ.Fire || attacker.Civ == civ.Nature
	})

}

func DarknessBlocker(card *match.Card, ctx *match.Context) {

	addBlockerCondition(card, ctx)

	conditionalBlocker(card, ctx, func(attacker *match.Card) bool {
		return attacker.Civ == civ.Darkness
	})

}

func LightBlocker(card *match.Card, ctx *match.Context) {

	addBlockerCondition(card, ctx)

	conditionalBlocker(card, ctx, func(attacker *match.Card) bool {
		return attacker.Civ == civ.Light
	})

}

func addBlockerCondition(card *match.Card, ctx *match.Context) {

	if _, ok := ctx.Event.(*match.UntapStep); ok {

		card.AddCondition(cnd.Blocker, true, card.ID)

	}

}

func defaultBlocker(card *match.Card, ctx *match.Context) {
	conditionalBlocker(card, ctx, nil)
}

func conditionalBlocker(card *match.Card, ctx *match.Context, attackerTest func(*match.Card) bool) {

	// Check to see if I'm eligible to actually block the attack
	// (I'm in the battlezone AND I'm untapped AND opponent doesn't attack me)
	if event, ok := ctx.Event.(*match.SelectBlockers); ok &&
		card.Zone == match.BATTLEZONE &&
		(event.AttackedCard == nil || event.AttackedCard != card) &&
		!card.Tapped {

		// ONLY if the attacker card (event.CardWhoAttacked) passes the test function (if given)
		// and is my opponent who attacks me
		// add myself to the blockers list
		if event.CardWhoAttacked.Zone == match.BATTLEZONE &&
			card.Player != event.CardWhoAttacked.Player &&
			(attackerTest == nil || attackerTest(event.CardWhoAttacked)) {
			event.Blockers = append(event.Blockers, card)
		}

	}

}
