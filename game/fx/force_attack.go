package fx

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"fmt"
)

// ForceAttack prevents the user from ending their turn if the card has not attacked this turn
func ForceAttack(card *match.Card, ctx *match.Context) {

	if _, ok := ctx.Event.(*match.EndTurnEvent); ok && card.Zone == match.BATTLEZONE {

		if ctx.Match.IsPlayerTurn(card.Player) && !card.HasCondition(cnd.SummoningSickness) && !card.Tapped {

			if card.HasCondition(cnd.CantAttackPlayers) {

				if card.HasCondition(cnd.CantAttackCreatures) {
					return
				}

				attackableCreatures := FindFilter(
					ctx.Match.Opponent(card.Player),
					match.BATTLEZONE,
					func(c *match.Card) bool { return c.Tapped || card.HasCondition(cnd.AttackUntapped) })

				if len(attackableCreatures) == 0 {
					return
				}

			}

			ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("%s must attack before you can end your turn", card.Name))
			ctx.InterruptFlow()

		}

	}

}
