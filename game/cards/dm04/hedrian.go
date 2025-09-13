package dm04

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// Locomotiver ...
func Locomotiver(c *match.Card) {

	c.Name = "Locomotiver"
	c.Power = 1000
	c.Civ = civ.Darkness
	c.Family = []string{family.Hedrian}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.ShieldTrigger, fx.When(fx.Summoned, fx.OpponentDiscardsRandomCard))
}

// MongrelMan ...
func MongrelMan(c *match.Card) {

	c.Name = "Mongrel Man"
	c.Power = 2000
	c.Civ = civ.Darkness
	c.Family = []string{family.Hedrian}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.When(fx.InTheBattlezone, func(card *match.Card, ctx *match.Context) {
		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
			if card.Zone != match.BATTLEZONE {
				exit()
				return
			}

			if fx.AnotherCreatureDestroyed(card, ctx2) {
				fx.MayDraw1(card, ctx2)
			}
		})
	}))

}
