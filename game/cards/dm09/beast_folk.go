package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// CavernRaider ...
func CavernRaider(c *match.Card) {

	c.Name = "Cavern Raider"
	c.Power = 2000
	c.Civ = civ.Nature
	c.Family = []string{family.BeastFolk}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.When(fx.WheneverThisAttacksPlayerAndIsntBlocked, fx.SearchDeckTake1Creature))

}
