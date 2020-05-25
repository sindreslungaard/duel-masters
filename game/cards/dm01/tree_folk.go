package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// CoilingVines ...
func CoilingVines(c *match.Card) {

	c.Name = "Coiling Vines"
	c.Power = 3000
	c.Civ = civ.Nature
	c.Family = family.TreeFolk
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.ReturnToMana)

}

// PoisonousDahlia ...
func PoisonousDahlia(c *match.Card) {

	c.Name = "Poisonous Dahlia"
	c.Power = 5000
	c.Civ = civ.Nature
	c.Family = family.TreeFolk
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.CantAttackPlayers)

}
