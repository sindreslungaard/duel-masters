package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// HanusaRadianceElemental ...
func HanusaRadianceElemental(c *match.Card) {

	c.Name = "Hanusa, Radiance Elemental"
	c.Power = 9500
	c.Civ = civ.Light
	c.Family = []string{family.AngelCommand}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Doublebreaker)

}

// UrthPurifyingElemental ...
func UrthPurifyingElemental(c *match.Card) {

	c.Name = "Urth, Purifying Elemental"
	c.Power = 6000
	c.Civ = civ.Light
	c.Family = []string{family.AngelCommand}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Doublebreaker, fx.Untap)

}
