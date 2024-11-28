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

	c.Use(fx.Creature, fx.Doublebreaker, fx.When(fx.Attacking, func(card *match.Card, ctx *match.Context) {

		shields := fx.Find(
			card.Player,
			match.SHIELDZONE,
		)

		if len(shields) < 1 {
			return
		}

		ctx.Match.MoveCard(shields[0], match.GRAVEYARD, card)

	}))

}
