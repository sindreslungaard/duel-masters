package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// ArdentLunatron ...
func ArdentLunatron(c *match.Card) {

	c.Name = "Ardent Lunatron"
	c.Power = 5000
	c.Civ = civ.Water
	c.Family = []string{family.CyberMoon}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Blocker(), fx.CantAttackPlayers, fx.CantAttackCreatures, fx.BlockIfAbleWhenOppAttacks)

}
