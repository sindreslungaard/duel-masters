package fx

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
)

// PowerBreakerTiers keeps "double breaker" and "triple breaker" in sync with a
// creature's current effective power. It implements cards printed as "while
// this creature has power X or more, it has double breaker" together with the
// optional "while this creature has power Y or more, it has triple breaker
// instead of double breaker".
//
// Pass a non-positive tripleAt for a card that only has the double breaker tier.
// Both conditions are tagged with the card's own ID so a breaker granted by
// another card is never removed, and both are dropped while the card is outside
// the battle zone.
func PowerBreakerTiers(doubleAt int, tripleAt int) match.HandlerFunc {
	return func(card *match.Card, ctx *match.Context) {
		// GetPower dispatches handlers synchronously. Do not recursively ask for
		// power while already handling that calculation.
		if _, calculatingPower := ctx.Event.(*match.GetPowerEvent); calculatingPower {
			return
		}

		hasDoubleBreaker := hasSelfSourcedCondition(card, cnd.DoubleBreaker)
		hasTripleBreaker := hasSelfSourcedCondition(card, cnd.TripleBreaker)

		if card.Zone != match.BATTLEZONE {
			if hasDoubleBreaker {
				card.RemoveSpecificConditionBySource(cnd.DoubleBreaker, card.ID)
			}
			if hasTripleBreaker {
				card.RemoveSpecificConditionBySource(cnd.TripleBreaker, card.ID)
			}
			return
		}

		attacking := false
		switch event := ctx.Event.(type) {
		case *match.AttackPlayer:
			attacking = event.CardID == card.ID
		case *match.AttackCreature:
			attacking = event.CardID == card.ID
		case *match.SelectShields:
			attacking = event.Attacker == card
		}

		power := ctx.Match.GetPower(card, attacking)

		// "Triple breaker" replaces "double breaker" rather than stacking with it.
		wantsTripleBreaker := tripleAt > 0 && power >= tripleAt
		wantsDoubleBreaker := power >= doubleAt && !wantsTripleBreaker

		if wantsDoubleBreaker && !hasDoubleBreaker {
			card.AddUniqueSourceCondition(cnd.DoubleBreaker, true, card.ID)
		} else if !wantsDoubleBreaker && hasDoubleBreaker {
			card.RemoveSpecificConditionBySource(cnd.DoubleBreaker, card.ID)
		}

		if wantsTripleBreaker && !hasTripleBreaker {
			card.AddUniqueSourceCondition(cnd.TripleBreaker, true, card.ID)
		} else if !wantsTripleBreaker && hasTripleBreaker {
			card.RemoveSpecificConditionBySource(cnd.TripleBreaker, card.ID)
		}
	}
}

func hasSelfSourcedCondition(card *match.Card, id string) bool {
	for _, condition := range card.Conditions() {
		if condition.ID == id && condition.Src == card.ID {
			return true
		}
	}

	return false
}
