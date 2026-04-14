package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// TrixoWickedDoll ...
func TrixoWickedDoll(c *match.Card) {

	c.Name = "Trixo, Wicked Doll"
	c.Power = 2000
	c.Civ = civ.Darkness
	c.Family = []string{family.DeathPuppet}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.When(fx.WheneverThisAttacksPlayerAndIsntBlocked, fx.OpponentChoosesAndDestroysCreature))

}
