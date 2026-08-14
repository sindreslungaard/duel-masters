package fx

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
)

// CrewBreaker implements "crew breaker", where a creature breaks one more shield
// for each of your other creatures of a given race in the battle zone. Pass a
// count of those other creatures.
//
// The modifier is kept in sync on every event rather than being applied for the
// duration of an attack. HandleFx skips the remaining scheduled callbacks once a
// context is cancelled, so an attack that is called off part way through - by
// cancelling the shield selection, for example - would otherwise leave a stale
// modifier behind for the rest of the turn.
func CrewBreaker(count func(card *match.Card) int) match.HandlerFunc {
	return func(card *match.Card, ctx *match.Context) {
		// GetPower dispatches handlers synchronously; there is nothing to do
		// while a power calculation is already in flight.
		if _, calculatingPower := ctx.Event.(*match.GetPowerEvent); calculatingPower {
			return
		}

		wanted := 0
		if card.Zone == match.BATTLEZONE {
			wanted = count(card)
		}

		current, has := ownShieldBreakModifier(card)
		if has && current == wanted {
			return
		}

		// Only this card's own contribution is touched, so a modifier granted by
		// another card is never removed.
		if has {
			card.RemoveSpecificConditionBySource(cnd.ShieldBreakModifier, card.ID)
		}

		if wanted > 0 {
			card.AddUniqueSourceCondition(cnd.ShieldBreakModifier, wanted, card.ID)
		}
	}
}

// CountOtherOwnCreaturesWithFamily counts the card's controller's other
// creatures in the battle zone that share one of the given families.
func CountOtherOwnCreaturesWithFamily(families []string) func(card *match.Card) int {
	return func(card *match.Card) int {
		return len(FindFilter(
			card.Player,
			match.BATTLEZONE,
			func(x *match.Card) bool {
				return x.ID != card.ID && x.SharesAFamily(families)
			},
		))
	}
}

func ownShieldBreakModifier(card *match.Card) (int, bool) {
	for _, condition := range card.Conditions() {
		if condition.ID != cnd.ShieldBreakModifier || condition.Src != card.ID {
			continue
		}

		if val, ok := condition.Val.(int); ok {
			return val, true
		}
	}

	return 0, false
}
