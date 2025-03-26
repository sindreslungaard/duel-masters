package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/match"
)

// ScreamSlicerShadowOfFear ...
func ScreamSlicerShadowOfFear(c *match.Card) {

	c.Name = "Scream Slicer, Shadow of Fear"
	c.Power = 4000
	c.Civ = civ.Darkness
	c.Family = []string{family.Ghost}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Darkness}

	//TODO:
	// Whenever you put a Dragonoid or a creature that has Dragon in its race
	// into the battlezone, destroy the creature that has the least power in the BZ.
	// If there's a tie, you choose among the tied creatures.
	//TODO 2:
	// Implement fx.AnotherCreatureFilterSummoned,
	// with a filter on SharesAFamily: Dragons + Dragonoid
}
