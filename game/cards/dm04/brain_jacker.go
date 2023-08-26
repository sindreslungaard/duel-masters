package dm04

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// PurplePiercer ...
func PurplePiercer(c *match.Card) {

	c.Name = "Purple Piercer"
	c.Power = 2000
	c.Civ = civ.Darkness
	c.Family = []string{family.BrainJacker}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(
		fx.Creature,
		fx.CantBeBlockedIf(func(blocker *match.Card) bool {
			return blocker.Civ == civ.Light
		}),
		fx.CantBeAttackedIf(func(attacker *match.Card) bool {
			return attacker.Civ == civ.Light
		}),
	)

}
