package dm12

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// BelmolTheExplorer ...
func BelmolTheExplorer(c *match.Card) {

	c.Name = "Belmol, the Explorer"
	c.Power = 3500
	c.Civs = []string{civ.Light}
	c.Family = []string{family.Gladiator}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Blocker(), fx.BlockIfAbleWhenOppAttacks, fx.CantAttackPlayers)

}

// ElectroExplorerSyrion ...
func ElectroExplorerSyrion(c *match.Card) {

	c.Name = "Electro Explorer Syrion"
	c.Power = 4000
	c.Civs = []string{civ.Light, civ.Water}
	c.Family = []string{family.Gladiator, family.CyberLord}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light, civ.Water}

	c.Use(fx.Creature, fx.PutIntoManaZoneTapped)

}

// BingoleTheExplorer ...
func BingoleTheExplorer(c *match.Card) {

	c.Name = "Bingole, the Explorer"
	c.Power = 4000
	c.Civs = []string{civ.Light}
	c.Family = []string{family.Gladiator}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.PutIntoBattleZoneInsteadOfDiscard)

}
