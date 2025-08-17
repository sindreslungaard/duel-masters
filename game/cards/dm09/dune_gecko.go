package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// SnaptongueLizard ...
func SnaptongueLizard(c *match.Card) {

	c.Name = "Snaptongue Lizard"
	c.Power = 3000
	c.Civ = civ.Fire
	c.Family = []string{family.DuneGecko}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.PowerAttacker3000, fx.CantBeBlockedWhileAttackingACreature)

}
