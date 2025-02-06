package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// TerradragonRegarion ...
func TerradragonRegarion(c *match.Card) {

	c.Name = "Terradragon Regarion"
	c.Power = 4000
	c.Civ = civ.Nature
	c.Family = []string{family.EarthDragon}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.Doublebreaker, fx.PowerAttacker3000)
}
