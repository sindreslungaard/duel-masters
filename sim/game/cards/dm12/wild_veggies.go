package dm12

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// UncannyTurnip ...
func UncannyTurnip(c *match.Card) {

	c.Name = "Uncanny Turnip"
	c.Power = 1000
	c.Civs = []string{civ.Nature}
	c.Family = []string{family.WildVeggies}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.WaveStriker, fx.WhileWaveStriker(fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		// The card that goes in is not necessarily the one that comes out: the
		// deck feeds the mana zone first, and then any creature there can be
		// taken to hand.
		fx.Draw1ToMana(card, ctx)
		fx.ReturnCreatureFromManazoneToHand(card, ctx)
	})))

}
