package dm11

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// BairaTheHiddenLunatic ...
func BairaTheHiddenLunatic(c *match.Card) {

	c.Name = "Baira, the Hidden Lunatic"
	c.Power = 5000
	c.Civs = []string{civ.Darkness}
	c.Family = []string{family.PandorasBox}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Darkness}

	// It can only ever battle by blocking, and destroying itself afterwards is
	// what stops a 5000 power blocker holding the board on its own.
	c.Use(fx.Creature, fx.Blocker(), fx.CantAttackCreatures, fx.CantAttackPlayers, fx.Suicide)

}

// BeratchaTheHiddenGlutton ...
func BeratchaTheHiddenGlutton(c *match.Card) {

	c.Name = "Beratcha, the Hidden Glutton"
	c.Power = 3000
	c.Civs = []string{civ.Darkness}
	c.Family = []string{family.PandorasBox}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Slayer)

}
