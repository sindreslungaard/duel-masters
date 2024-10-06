package dm07

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func DreamPirateShadowOfTheft(c *match.Card) {

	c.Name = "Dream Pirate, Shadow of Theft"
	c.Power = 3000
	c.Civ = civ.Darkness
	c.Family = []string{family.Ghost}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.When(fx.WouldBeDestroyed, fx.MayReturnToHandAndDiscardACard))

}
