package fx

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"fmt"
)

// SilentSkill implements "Silent skill (At the start of each of your turns, if
// this creature is tapped, you may keep it tapped and use its ability.)"
//
// The printed timing straddles the untap, so the mechanic is split across two
// steps. BeginTurnStep runs before anything untaps, and is the only place the
// creature's tapped state at the start of the turn can still be observed.
// StartOfTurnStep is where the ability actually resolves, because UntapStep is
// what rebuilds every creature's conditions: an ability resolved any earlier
// would look at a battle zone where nothing has "blocker", no power modifier is
// in place, and even cnd.Creature is missing, all of which were cleared at the
// previous EndOfTurnStep.
//
// Keeping the creature tapped is therefore a re-tap after the fact rather than
// a prevented untap. That is invisible in play — nothing happens between the
// two steps — and it keeps the mechanic out of fx.Creature's untap path, which
// every creature in the game runs.
func SilentSkill(ability match.HandlerFunc) match.HandlerFunc {
	return func(card *match.Card, ctx *match.Context) {

		// Identity for the cards that care which creatures have silent skill.
		// Added in every zone and rebuilt each untap, the same way blocker and
		// the other printed keywords are, so an effect can ask about a creature
		// that is still in hand.
		if _, ok := ctx.Event.(*match.UntapStep); ok {
			card.AddUniqueSourceCondition(cnd.SilentSkill, true, card.ID)
			return
		}

		if _, ok := ctx.Event.(*match.BeginTurnStep); ok {
			if card.Zone != match.BATTLEZONE || !ctx.Match.IsPlayerTurn(card.Player) || !card.Tapped {
				return
			}

			card.AddUniqueSourceCondition(cnd.SilentSkillReady, true, card.ID)
			return
		}

		if _, ok := ctx.Event.(*match.StartOfTurnStep); ok {
			if !card.HasCondition(cnd.SilentSkillReady) {
				return
			}

			// Resolve after the rest of the step so the ability sees the board
			// every other start-of-turn effect has already left behind.
			ctx.ScheduleAfter(func() {
				card.RemoveSpecificConditionBySource(cnd.SilentSkillReady, card.ID)

				// Something may have removed the creature between the two
				// steps, and an ability of a creature that is no longer in the
				// battle zone must not resolve.
				if card.Zone != match.BATTLEZONE {
					return
				}

				if !BinaryQuestion(card.Player, ctx.Match, fmt.Sprintf("Do you want to keep %s tapped and use its silent skill?", card.Name)) {
					return
				}

				card.Tapped = true

				ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s used %s's silent skill", card.Player.Username(), card.Name))

				ability(card, ctx)
			})
		}

	}
}
