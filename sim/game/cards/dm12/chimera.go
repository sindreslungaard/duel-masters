package dm12

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// Gigaslug ...
func Gigaslug(c *match.Card) {

	c.Name = "Gigaslug"
	c.Power = 1000
	c.Civs = []string{civ.Darkness}
	c.Family = []string{family.Chimera}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Blocker(), fx.Slayer, fx.CantAttackCreatures, fx.CantAttackPlayers)

}

// Gigabalza ...
func Gigabalza(c *match.Card) {

	c.Name = "Gigabalza"
	c.Power = 1000
	c.Civs = []string{civ.Darkness}
	c.Family = []string{family.Chimera}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.ShieldTrigger, fx.When(fx.Summoned, fx.OpponentDiscardsRandomCard))

}

// GigappiPonto ...
func GigappiPonto(c *match.Card) {

	c.Name = "Gigappi Ponto"
	c.Power = 4000
	c.Civs = []string{civ.Darkness, civ.Fire}
	c.Family = []string{family.Chimera, family.FireBird}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Darkness, civ.Fire}

	c.Use(fx.Creature, fx.PutIntoManaZoneTapped)

}
