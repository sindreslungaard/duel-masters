package dm04

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// GulanRiasSpeedGuardian ...
func GulanRiasSpeedGuardian(c *match.Card) {

	c.Name = "Gulan Rias, Speed Guardian"
	c.Power = 2000
	c.Civ = civ.Light
	c.Family = []string{family.Guardian}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light}

	c.Use(
		fx.Creature,
		fx.CantBeBlockedByDarkness,
		fx.CantBeAttackedIf(func(attacker *match.Card) bool {
			return attacker.Civ == civ.Darkness
		}),
	)
}

// MistRiasSonicGuardian ...
func MistRiasSonicGuardian(c *match.Card) {

	c.Name = "Mist Rias, Sonic Guardian"
	c.Power = 2000
	c.Civ = civ.Light
	c.Family = []string{family.Guardian}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.When(fx.InTheBattlezone, func(card *match.Card, ctx *match.Context) {
		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
			if card.Zone != match.BATTLEZONE {
				exit()
				return
			}

			if fx.AnotherCreatureSummoned(card, ctx2) {
				fx.MayDraw1(card, ctx2)
			}
		})
	}))

}
