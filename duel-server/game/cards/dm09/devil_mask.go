package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// AcidRefluxTheFleshboiler ...
func AcidRefluxTheFleshboiler(c *match.Card) {

	c.Name = "Acid Reflux, the Fleshboiler"
	c.Power = 3000
	c.Civ = civ.Darkness
	c.Family = []string{family.DevilMask}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Blocker(), fx.CantAttackCreatures, fx.CantAttackPlayers, fx.Slayer)

}
