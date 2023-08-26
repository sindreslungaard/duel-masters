package dm03

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// BaragaBladeOfGloom ...
func BaragaBladeOfGloom(c *match.Card) {

	c.Name = "Baraga, Blade of Gloom"
	c.Power = 4000
	c.Civ = civ.Darkness
	c.Family = []string{family.DarkLord}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		if match.AmISummoned(card, ctx) {

			fx.SelectBackside(
				card.Player,
				ctx.Match,
				card.Player,
				match.SHIELDZONE,
				"Baraga, Blade of Gloom: Move 1 of your shield to your hand.",
				1,
				1,
				false,
			).Map(func(x *match.Card) {
				ctx.Match.MoveCard(x, match.HAND, card)
			})

		}
	})

}
