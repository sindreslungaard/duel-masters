package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// SpinalParasite ...
func SpinalParasite(c *match.Card) {

	c.Name = "Spinal Parasite"
	c.Power = 2000
	c.Civ = civ.Darkness
	c.Family = []string{family.BrainJacker}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {
		if _, ok := ctx.Event.(*match.StartOfTurnStep); ok && card.Zone == match.BATTLEZONE && !ctx.Match.IsPlayerTurn(card.Player) {
			fx.SelectFilter(
				ctx.Match.Opponent(card.Player),
				ctx.Match,
				ctx.Match.Opponent(card.Player),
				match.BATTLEZONE,
				fmt.Sprintf("%s's effect: Choose one of your creatures in the battle zone that can attack. That creature attacks this turn if able.", card.Name),
				1,
				1,
				false,
				func(x *match.Card) bool {
					return fx.CanAttack(x)
				},
				false,
			).Map(func(x *match.Card) {
				ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
					if x.Zone != match.BATTLEZONE {
						exit()
						return
					}

					if _, ok := ctx2.Event.(*match.EndOfTurnStep); ok {
						exit()
						return
					}

					fx.ForceAttack(x, ctx2)
				})
			})
		}
	})

}
