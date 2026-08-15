package fx

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"fmt"
)

// canBeForcedToAttack reports whether card's controller has a legal attack
// available with it right now, on their own turn, which is what "attacks if
// able" turns on.
func canBeForcedToAttack(card *match.Card, m *match.Match) bool {

	if !m.IsPlayerTurn(card.Player) || HasSummoningSickness(card) || card.Tapped {
		return false
	}

	if card.HasCondition(cnd.CantAttackPlayers) {

		if card.HasCondition(cnd.CantAttackCreatures) {
			return false
		}

		attackableCreatures := FindFilter(
			m.Opponent(card.Player),
			match.BATTLEZONE,
			func(c *match.Card) bool { return c.Tapped || card.HasCondition(cnd.AttackUntapped) })

		if len(attackableCreatures) == 0 {
			return false
		}

	}

	return true

}

// ForceAttack prevents the user from ending their turn if the card has not
// attacked this turn, and prevents using a tap ability instead of attacking
// while an attack is legally available.
//
// Per the official Slime Veil ruling: "On your turn, when you have the option
// either to attack with your creature or to use its tap ability, [this
// effect] says that you must attack with it."
func ForceAttack(card *match.Card, ctx *match.Context) {

	if card.Zone != match.BATTLEZONE {
		return
	}

	if event, ok := ctx.Event.(*match.TapAbility); ok {

		if event.CardID != card.ID || !canBeForcedToAttack(card, ctx.Match) {
			return
		}

		ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("%s must attack instead of using its tap ability", card.Name))
		ctx.InterruptFlow()

		return

	}

	if _, ok := ctx.Event.(*match.EndTurnEvent); ok && canBeForcedToAttack(card, ctx.Match) {

		ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("%s must attack before you can end your turn", card.Name))
		ctx.InterruptFlow()

	}

}
