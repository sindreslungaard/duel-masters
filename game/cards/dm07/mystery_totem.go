package dm07

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func CrypticTotem(c *match.Card) {

	c.Name = "Cryptic Totem"
	c.Power = 6000
	c.Civ = civ.Nature
	c.Family = []string{family.MysteryTotem}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.Doublebreaker, fx.When(fx.Tapped, func(card *match.Card, ctx *match.Context) {
		fx.FilterShieldTriggers(ctx, func(x *match.Card) bool { return card.Player == x.Player })
	}))
}
