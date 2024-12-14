package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// AquaRanger ...
func AquaRanger(c *match.Card) {

	c.Name = "Aqua Ranger"
	c.Power = 3000
	c.Civ = civ.Water
	c.Family = []string{family.LiquidPeople}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.CantBeBlocked, fx.When(fx.WouldBeDestroyed, fx.ReturnToHand))

}
