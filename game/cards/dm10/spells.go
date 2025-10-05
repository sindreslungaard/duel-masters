package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// InfernalCommand ...
func InfernalCommand(c *match.Card) {

	c.Name = "Infernal Command"
	c.Civ = civ.Darkness
	c.ManaCost = 1
	c.ManaRequirement = []string{civ.Darkness}

	//@TODO test this with the Nariel fix with Machai and ForceAttack!!
	// What happens if the opponent casts InfernalCommand or Slime Veil?
	// Or any other spell that gives ForceAttack to other creatures that are blocked from attacking by Nariel??

	// Perhaps modify Nariel effect to NOT remove CantAttack conditions on EndOfTurn event
	// so that creatures that are forced to attack by other effects will still be forced to attack
	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
			// remove persistent effect on start of next turn
			if _, ok := ctx2.Event.(*match.StartOfTurnStep); ok && ctx2.Match.IsPlayerTurn(card.Player) {
				exit()
				return
			}

			// on all events, add force attack to opponent's creatures
			fx.Find(
				ctx2.Match.Opponent(card.Player),
				match.BATTLEZONE,
			).Map(func(c *match.Card) {
				if _, ok := ctx2.Event.(*match.EndTurnEvent); ok && c.Zone == match.BATTLEZONE {
					if ctx2.Match.IsPlayerTurn(c.Player) && !fx.HasSummoningSickness(c) && !c.Tapped {
						if c.HasCondition(cnd.CantAttackPlayers) {
							if c.HasCondition(cnd.CantAttackCreatures) {
								return
							}

							attackableCreatures := fx.FindFilter(
								ctx2.Match.Opponent(c.Player),
								match.BATTLEZONE,
								func(x *match.Card) bool { return x.Tapped || c.HasCondition(cnd.AttackUntapped) })

							if len(attackableCreatures) == 0 {
								return
							}
						}

						ctx2.Match.WarnPlayer(c.Player, fmt.Sprintf("%s must attack before you can end your turn", c.Name))
						ctx2.InterruptFlow()
					}
				}
			})
		})
	}))

}
