package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// BakkraHornTheSilent ...
func BakkraHornTheSilent(c *match.Card) {

	c.Name = "Bakkra Horn, the Silent"
	c.Power = 2000
	c.Civ = civ.Nature
	c.Family = []string{family.HornedBeast}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature,
		fx.When(
			fx.AnotherOwnDragonoidOrDragonSummoned,
			fx.Draw1ToMana,
		))
}
