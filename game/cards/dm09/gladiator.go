package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// BetraleTheExplorer ...
func BetraleTheExplorer(c *match.Card) {

	c.Name = "Betrale, the Explorer"
	c.Power = 5000
	c.Civ = civ.Light
	c.Family = []string{family.Gladiator}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.CantAttackPlayers, fx.Blocker(),
		fx.When(fx.EndOfMyTurnCreatureBZ, fx.MayUntapSelf))

}
