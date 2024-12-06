package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
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

	c.Use(fx.Creature, fx.Doublebreaker, fx.When(fx.AttackConfirmed, func(card *match.Card, ctx *match.Context) {
		fx.SelectBackside(
			card.Player,
			ctx.Match,
			card.Player,
			match.SHIELDZONE,
			"Select one of your shields and put it into your graveyard",
			1,
			1,
			false,
		).Map(func(x *match.Card) {
			ctx.Match.MoveCard(x, match.GRAVEYARD, card)
		})
	}))
}

// SuperNecrodragonAbzoDolba ...
func SuperNecrodragonAbzoDolba(c *match.Card) {

	c.Name = "Super Necrodragon Abzo Dolba"
	c.Power = 11000
	c.Civ = civ.Darkness
	c.Family = []string{family.ZombieDragon}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Darkness}

	c.PowerModifier = func(m *match.Match, attacking bool) int {

		creaturesInGrave := fx.FindFilter(c.Player, match.GRAVEYARD, func(x *match.Card) bool { return x.HasCondition(cnd.Creature) })

		return len(creaturesInGrave) * 2000

	}

	c.Use(fx.Creature, fx.DragonEvolution, fx.Triplebreaker)

}
