package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// DarkRavenShadowOfGrief ...
func DarkRavenShadowOfGrief(c *match.Card) {

	c.Name = "Dark Raven, Shadow of Grief"
	c.Power = 1000
	c.Civ = civ.Darkness
	c.Family = family.Ghost
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Blocker)

}
