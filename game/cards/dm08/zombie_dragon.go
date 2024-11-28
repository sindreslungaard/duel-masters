package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// NecrodragonGiland ...
func NecrodragonGiland(c *match.Card) {

	c.Name = "Necrodragon Giland"
	c.Power = 6000
	c.Civ = civ.Darkness
	c.Family = []string{family.ZombieDragon}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Doublebreaker, fx.Suicide)

}

// NecrodragonGalbazeek ...
func NecrodragonGalbazeek(c *match.Card) {

	c.Name = "Necrodragon Galbazeek"
	c.Power = 9000
	c.Civ = civ.Darkness
	c.Family = []string{family.ZombieDragon}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Doublebreaker, fx.WheneverThisAttacks(func(card *match.Card, ctx *match.Context) {

		fx.SelectBackside(
			card.Player,
			ctx.Match, card.Player,
			match.SHIELDZONE,
			"Select one shield and send it to graveyard",
			1,
			1,
			false,
		).Map(func(x *match.Card) {
			ctx.Match.MoveCard(x, match.GRAVEYARD, card)
		})

	}))

}
