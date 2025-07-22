package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// AzaghastTyrantOfShadows ...
func AzaghastTyrantOfShadows(c *match.Card) {

	c.Name = "Azaghast, Tyrant of Shadows"
	c.Power = 9000
	c.Civ = civ.Darkness
	c.Family = []string{family.DarkLord}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Evolution, fx.Doublebreaker,
		fx.When(fx.AnotherOwnGhostSummoned, func(card *match.Card, ctx *match.Context) {
			fx.SelectFilter(
				card.Player,
				ctx.Match,
				ctx.Match.Opponent(card.Player),
				match.BATTLEZONE,
				fmt.Sprintf("%s's effect: you may destroy 1 of your opponent's untapped creatures.", card.Name),
				1,
				1,
				true,
				func(x *match.Card) bool {
					return !x.Tapped
				},
				false,
			).Map(func(x *match.Card) {
				ctx.Match.Destroy(x, card, match.DestroyedByMiscAbility)
			})
		}))

}
