package fx

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"slices"
)

// Blocker adds the Blocker condition at the Untap step and
// Adds the card to the blockers list at the SelectBlockers step if eligible
func Blocker(card *match.Card, ctx *match.Context) {

	addDefaultBlockerCondition(card, ctx, false, nil)

	addDefaultBlocker(card, ctx)

}

// ForceBlocker adds the Blocker condition, no matter the current handled event
// And adds the card to the blockers list at the SelectBlockers step if eligible
func ForceBlocker(card *match.Card, ctx *match.Context, src any) {

	forceAddBlockerCondition(card, src)

	addDefaultBlocker(card, ctx)
}

func FireAndNatureBlocker(card *match.Card, ctx *match.Context) {

	addDefaultBlockerCondition(card, ctx, false, nil)

	addConditionalBlocker(card, ctx, func(attacker *match.Card) bool {
		return attacker.Civ == civ.Fire || attacker.Civ == civ.Nature
	})

}

func DarknessBlocker(card *match.Card, ctx *match.Context) {

	addDefaultBlockerCondition(card, ctx, false, nil)

	addConditionalBlocker(card, ctx, func(attacker *match.Card) bool {
		return attacker.Civ == civ.Darkness
	})

}

func LightBlocker(card *match.Card, ctx *match.Context) {

	addDefaultBlockerCondition(card, ctx, false, nil)

	addConditionalBlocker(card, ctx, func(attacker *match.Card) bool {
		return attacker.Civ == civ.Light
	})

}

func addDefaultBlockerCondition(card *match.Card, ctx *match.Context, uniqueCond bool, src any) {

	if _, ok := ctx.Event.(*match.UntapStep); ok {

		if uniqueCond {
			card.AddUniqueSourceCondition(cnd.Blocker, true, src)
		} else {
			card.AddCondition(cnd.Blocker, true, card.ID)
		}

	}

}

func forceAddBlockerCondition(card *match.Card, src any) {
	card.AddUniqueSourceCondition(cnd.Blocker, true, src)
}

func addDefaultBlocker(card *match.Card, ctx *match.Context) {
	addConditionalBlocker(card, ctx, nil)
}

func addConditionalBlocker(card *match.Card, ctx *match.Context, attackerTest func(*match.Card) bool) {

	// Check to see if I'm eligible to actually block the attack
	// (I'm in the battlezone AND I'm untapped AND opponent doesn't attack me)
	if event, ok := ctx.Event.(*match.SelectBlockers); ok &&
		card.Zone == match.BATTLEZONE &&
		(event.AttackedCard == nil || event.AttackedCard != card) &&
		!card.Tapped &&
		card.HasCondition(cnd.Blocker) {

		// ONLY if the attacker card (event.CardWhoAttacked) passes the test function (if given)
		// and is my opponent who attacks me
		// add myself to the blockers list
		if event.CardWhoAttacked.Zone == match.BATTLEZONE &&
			card.Player != event.CardWhoAttacked.Player &&
			(attackerTest == nil || attackerTest(event.CardWhoAttacked)) &&
			!slices.Contains(event.Blockers, card) {
			event.Blockers = append(event.Blockers, card)
		}

	}

}
