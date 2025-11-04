package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// ExplosiveTrooperZalmez ...
func ExplosiveTrooperZalmez(c *match.Card) {

	c.Name = "Explosive Trooper Zalmez"
	c.Power = 2000
	c.Civ = civ.Fire
	c.Family = []string{family.Armorloid}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		oppShields, err := ctx.Match.Opponent(card.Player).Container(match.SHIELDZONE)

		if err == nil && len(oppShields) <= 2 {
			fx.SelectFilter(
				card.Player,
				ctx.Match,
				ctx.Match.Opponent(card.Player),
				match.BATTLEZONE,
				fmt.Sprintf("%s's effect: You may destroy one of your opponent's creatures that has power 3000 or less.", card.Name),
				1,
				1,
				true,
				func(x *match.Card) bool {
					return ctx.Match.GetPower(x, false) <= 3000
				},
				false,
			).Map(func(x *match.Card) {
				ctx.Match.Destroy(x, card, match.DestroyedByMiscAbility)
			})
		}
	}))

}
