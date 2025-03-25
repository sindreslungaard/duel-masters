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

// AquaGrappler ...
func AquaGrappler(c *match.Card) {

	c.Name = "Aqua Grappler"
	c.Power = 3000
	c.Civ = civ.Water
	c.Family = []string{family.LiquidPeople}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.When(fx.AttackConfirmed, func(c *match.Card, ctx *match.Context) {
		//TODO draw a card for each other tapped creature you have in the BZ
	}))

}
