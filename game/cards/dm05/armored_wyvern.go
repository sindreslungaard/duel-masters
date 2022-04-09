package dm05

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// TwinCannonSkyterror ...
func TwinCannonSkyterror(c *match.Card) {

	c.Name = "Twin-Cannon Skyterror"
	c.Power = 7000
	c.Civ = civ.Fire
	c.Family = family.ArmoredWyvern
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.SpeedAttacker, fx.Doublebreaker)

}
