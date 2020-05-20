package fx

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
)

// Blocker adds the card to a list of blockers when a creature/player is attacked
func Blocker(card *match.Card, ctx *match.Context) {

	if _, ok := ctx.Event.(*match.UntapStep); ok {

		card.AddCondition(cnd.Blocker, true, card.ID)

	}

	if event, ok := ctx.Event.(*match.AttackPlayer); ok {

		// Only add to list of blockers if it is our player that is being attacked, i.e. not our players turn
		if !ctx.Match.IsPlayerTurn(card.Player) && !card.Tapped && card.Player.HasCard(match.BATTLEZONE, card.ID) {
			event.Blockers = append(event.Blockers, card)
		}

	}

	if event, ok := ctx.Event.(*match.AttackCreature); ok {

		// Only add to list of blockers if it is our creature that is being attacked, i.e. not our players turn
		if !ctx.Match.IsPlayerTurn(card.Player) && !card.Tapped && card.Player.HasCard(match.BATTLEZONE, card.ID) {
			event.Blockers = append(event.Blockers, card)
		}

	}

}
