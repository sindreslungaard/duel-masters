package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// CorpseCharger ...
func CorpseCharger(c *match.Card) {
	c.Name = "Corpse Charger"
	c.Civ = civ.Darkness
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Spell, fx.Charger, fx.When(fx.SpellCast, fx.ReturnXCreaturesFromGraveToHand(1)))
}

// CraniumClamp ...
func CraniumClamp(c *match.Card) {
	c.Name = "Cranium Clamp"
	c.Civ = civ.Darkness
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Spell, fx.When(fx.SpellCast, fx.OpDiscardsXCards(2)))
}
