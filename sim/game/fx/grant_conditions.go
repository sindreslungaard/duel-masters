package fx

import (
	"duel-masters/game/match"
	"fmt"
)

// GrantConditionToOwnCreatures gives a condition to every creature the player
// controls that matches the filter.
//
// The grant lasts until the end of the turn without any bookkeeping of its own:
// fx.Creature clears a creature's conditions at EndOfTurnStep, and only the
// intrinsic ones are rebuilt at the following untap step. Each contribution is
// tagged with the source card so two copies of the same effect do not tread on
// each other.
func GrantConditionToOwnCreatures(condition string, val interface{}, filter func(*match.Card) bool, description string) match.HandlerFunc {
	return func(card *match.Card, ctx *match.Context) {
		granted := 0

		FindFilter(card.Player, match.BATTLEZONE, filter).Map(func(creature *match.Card) {
			creature.AddUniqueSourceCondition(condition, val, card.ID)
			granted++
		})

		if granted > 0 {
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s's effect: %d of %s's creatures %s until the end of the turn", card.Name, granted, card.Player.Username(), description))
		}
	}
}

// PutCreatureFromManaIntoBZ puts a creature the player chooses from their own
// mana zone into the battle zone. The choice is mandatory and opens no prompt
// when no card in the mana zone matches.
func PutCreatureFromManaIntoBZ(filter func(*match.Card) bool, description string) match.HandlerFunc {
	return func(card *match.Card, ctx *match.Context) {
		SelectFilter(
			card.Player,
			ctx.Match,
			card.Player,
			match.MANAZONE,
			fmt.Sprintf("%s's effect: Put %s from your mana zone into the battle zone.", card.Name, description),
			1,
			1,
			false,
			filter,
			false,
		).Map(func(creature *match.Card) {
			ForcePutCreatureIntoBZ(ctx, creature, match.MANAZONE, card)
		})
	}
}
