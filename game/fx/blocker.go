package fx

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/match"
	"slices"
)

// Regular blocker and conditional blockers
// For usages of adding Blocker inside ApplyPersistentEffects,
// please use ForceBlocker
func Blocker(conditions ...func(event *match.SelectBlockers) bool) func(*match.Card, *match.Context) {
	return func(card *match.Card, ctx *match.Context) {

		// Always add the blocker cnd during the untap step
		if _, ok := ctx.Event.(*match.UntapStep); ok {
			card.AddUniqueSourceCondition(cnd.Blocker, true, card.ID)
			return
		}

		// Add itself to list of blockers if all conditions are met
		addToBlockersListIfConditions(card, ctx, conditions...)

	}
}

// To be used on ApplyPersistentEffects, passing "ctx2" context reference
func ForceBlocker(card *match.Card, ctx *match.Context, src string, conditions ...func(event *match.SelectBlockers) bool) {
	// Always add the blocker cnd, regardless of the current event
	// To be used on ApplyPersistentEffects
	card.AddUniqueSourceCondition(cnd.Blocker, true, src)

	// Add itself to list of blockers if all conditions are met
	addToBlockersListIfConditions(card, ctx, conditions...)
}

func DragonBlocker() func(*match.Card, *match.Context) {
	return Blocker(func(event *match.SelectBlockers) bool {
		return event.Attacker.SharesAFamily(family.Dragons)
	})
}

func DarknessBlocker() func(*match.Card, *match.Context) {
	return Blocker(func(event *match.SelectBlockers) bool {
		return event.Attacker.Civ == civ.Darkness
	})
}

func LightBlocker() func(*match.Card, *match.Context) {
	return Blocker(func(event *match.SelectBlockers) bool {
		return event.Attacker.Civ == civ.Light
	})
}

func FireAndNatureBlocker() func(*match.Card, *match.Context) {
	return Blocker(func(event *match.SelectBlockers) bool {
		return event.Attacker.Civ == civ.Fire ||
			event.Attacker.Civ == civ.Nature
	})
}

func BlockIfAbleWhenOppAttacks(card *match.Card, ctx *match.Context) {
	if event, ok := ctx.Event.(*match.Block); ok {
		if event.Attacker.Player == card.Player ||
			event.Attacker.Zone != match.BATTLEZONE ||
			card.Zone != match.BATTLEZONE {
			return
		}

		if len(event.Blockers) > 0 && event.Attacker.IsBlockable(ctx) {
			for _, b := range event.Blockers {
				if b.ID == card.ID {
					// Force the battle between the attacker and this card
					ctx.InterruptFlow()
					card.Tapped = true
					ctx.Match.Battle(event.Attacker, card, true, len(event.ShieldsAttacked) > 0)
					return
				}
			}
		}
	}
}

// Default helper function that adds the given card to the blockers list
// If all conditions passed are met
func addToBlockersListIfConditions(card *match.Card, ctx *match.Context, conditions ...func(event *match.SelectBlockers) bool) {
	event, ok := ctx.Event.(*match.SelectBlockers)

	if !ok {
		return
	}

	if card.Zone != match.BATTLEZONE {
		return
	}

	if event.AttackedCardID == card.ID {
		return
	}

	if card.Tapped {
		return
	}

	if !card.HasCondition(cnd.Blocker) {
		return
	}

	if event.Attacker.Zone != match.BATTLEZONE {
		return
	}

	if card.Player == event.Attacker.Player {
		return
	}

	if slices.Contains(event.Blockers, card) {
		return
	}

	for _, condition := range conditions {
		if !condition(event) {
			return
		}
	}

	event.Blockers = append(event.Blockers, card)
}
