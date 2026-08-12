package dm11

import (
	"duel-masters/game/civ"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// SolarTrap ...
func SolarTrap(c *match.Card) {

	c.Name = "Solar Trap"
	c.Civs = []string{civ.Light}
	c.ManaCost = 1
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.When(fx.SpellCast, fx.TapOpCreature))

}

// TenTonCrunch ...
func TenTonCrunch(c *match.Card) {

	c.Name = "Ten-Ton Crunch"
	c.Civs = []string{civ.Fire}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, fx.DestroyOpCreatureXPowerOrLess(3000, false, match.DestroyedBySpell)))

}

// MorbidMedicine ...
func MorbidMedicine(c *match.Card) {

	c.Name = "Morbid Medicine"
	c.Civs = []string{civ.Darkness}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Spell, fx.When(fx.SpellCast, fx.ReturnXCreaturesFromGraveToHand(2)))

}

// EmergencyTyphoon ...
func EmergencyTyphoon(c *match.Card) {

	c.Name = "Emergency Typhoon"
	c.Civs = []string{civ.Water}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		fx.DrawUpTo2(card, ctx)

		// The discard is not conditional on the draw: the printed text says to
		// draw up to 2 and then discard, so drawing none still costs a card.
		fx.DiscardOwnXCards(1)(card, ctx)
	}))

}
