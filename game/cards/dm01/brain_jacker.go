package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// BloodySquito ...
func BloodySquito(c *match.Card) {

	c.Name = "Bloody Squito"
	c.Power = 4000
	c.Civ = civ.Darkness
	c.Family = []string{family.BrainJacker}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Blocker, fx.CantAttackPlayers, fx.CantAttackCreatures, fx.Suicide)

}

// DarkClown ...
func DarkClown(c *match.Card) {

	c.Name = "Dark Clown"
	c.Power = 6000
	c.Civ = civ.Darkness
	c.Family = []string{family.BrainJacker}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Blocker, fx.CantAttackPlayers, fx.CantAttackCreatures, fx.Suicide)

}
