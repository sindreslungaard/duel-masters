package fx

import (
	"duel-masters/game/match"
)

// allOwnedZones is every zone a card of one player can sit in. An aura sweeps
// all of them when it is withdrawn, because a creature that left the battle
// zone keeps its conditions until the end of the turn and must not carry a
// stale bonus back into play.
var allOwnedZones = []string{
	match.BATTLEZONE,
	match.HAND,
	match.GRAVEYARD,
	match.MANAZONE,
	match.SHIELDZONE,
	match.HIDDENZONE,
}

// AuraForOwnCreatures installs a continuous grant that follows the battle zone.
//
// While the source is in play every creature of its controller matching the
// filter carries the condition; anything that stops matching, or leaves the
// battle zone, loses it again. The grant is re-evaluated on every event rather
// than applied once, because the set of creatures it covers changes as they are
// summoned and destroyed, and persistent effects run before any card handler so
// the condition is always in place before anything reads it.
//
// The filter is given the source as well as the candidate, so "each of your
// other X" can exclude the source itself.
func AuraForOwnCreatures(condition string, val interface{}, filter func(source *match.Card, candidate *match.Card) bool) match.HandlerFunc {
	return func(card *match.Card, ctx *match.Context) {
		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
			live := card.Zone == match.BATTLEZONE

			FindMultiple(card.Player, allOwnedZones).Map(func(x *match.Card) {
				if live && x.Zone == match.BATTLEZONE && filter(card, x) {
					x.AddUniqueSourceCondition(condition, val, card.ID)
					return
				}

				x.RemoveSpecificConditionBySource(condition, card.ID)
			})

			if !live {
				exit()
			}
		})
	}
}

// OtherCreaturesSharingAFamily filters to the source's kin, excluding itself.
func OtherCreaturesSharingAFamily(families []string) func(*match.Card, *match.Card) bool {
	return func(source *match.Card, candidate *match.Card) bool {
		return candidate.ID != source.ID && candidate.SharesAFamily(families)
	}
}

// CreaturesSharingAFamily filters to the source's kin, the source included.
func CreaturesSharingAFamily(families []string) func(*match.Card, *match.Card) bool {
	return func(source *match.Card, candidate *match.Card) bool {
		return candidate.SharesAFamily(families)
	}
}
