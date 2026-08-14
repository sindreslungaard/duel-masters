package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// SporeblastErengi ...
func SporeblastErengi(c *match.Card) {

	c.Name = "Sporeblast Erengi"
	c.Power = 4000
	c.Civs = []string{civ.Nature}
	c.Family = []string{family.BalloonMushroom}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.SilentSkill(fx.SearchDeckTake1Creature))

}
