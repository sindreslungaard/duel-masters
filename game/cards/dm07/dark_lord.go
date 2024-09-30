package dm07

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func CrathLadeMercilessKing(c *match.Card) {

	c.Name = "Crath Lade, Merciless King"
	c.Power = 4000
	c.Civ = civ.Darkness
	c.Family = []string{family.DarkLord}
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Darkness}
	c.TapAbility = fx.OpponentDiscards2RandomCards

	c.Use(fx.Creature, fx.TapAbility)
}
