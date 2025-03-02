package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// UberdragonBajula ...
func UberdragonBajula(c *match.Card) {

	c.Name = "\u00c3\u0153berdragon Bajula"
	c.Power = 13000
	c.Civ = civ.Fire
	c.Family = []string{family.ArmoredDragon}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.DragonEvolution, fx.Triplebreaker, fx.When(fx.AttackConfirmed, fx.ManaBurnX(2)))

}
