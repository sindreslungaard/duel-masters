package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// GajirabuteVileCenturion ...
func GajirabuteVileCenturion(c *match.Card) {
	c.Name = "Gajirabute, Vile Centurion"
	c.Power = 3000
	c.Civs = []string{civ.Darkness}
	c.Family = []string{family.DemonCommand}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.When(fx.Summoned, fx.PutOpShieldIntoGraveyard))
}
