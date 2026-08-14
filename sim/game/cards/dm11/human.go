package dm11

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// GankloakRogueCommando ...
func GankloakRogueCommando(c *match.Card) {

	c.Name = "Gankloak, Rogue Commando"
	c.Power = 2000
	c.Civs = []string{civ.Fire}
	c.Family = []string{family.Human}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.SilentSkill(fx.GrantConditionToOwnCreatures(
		cnd.DoubleBreaker,
		true,
		func(x *match.Card) bool { return x.HasCiv(civ.Fire) },
		"got double breaker",
	)))

}
