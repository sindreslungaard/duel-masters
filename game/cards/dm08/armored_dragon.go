package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// BruiserDragon ...
func BruiserDragon(c *match.Card) {

	c.Name = "Bruiser Dragon"
	c.Power = 5000
	c.Civ = civ.Fire
	c.Family = []string{family.ArmoredDragon}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.When(fx.WouldBeDestroyed, fx.PutShieldIntoGraveyard))
}
