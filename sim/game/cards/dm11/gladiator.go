package dm11

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// WarlordAilzonius ...
func WarlordAilzonius(c *match.Card) {

	c.Name = "Warlord Ailzonius"
	c.Power = 8000
	c.Civs = []string{civ.Light}
	c.Family = []string{family.Gladiator}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Evolution, fx.Doublebreaker, fx.CantBeSelectedByOpp)

}
