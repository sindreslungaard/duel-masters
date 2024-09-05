package dm04

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// FuReilSeekerOfStorms ...
func FuReilSeekerOfStorms(c *match.Card) {

	c.Name = "Fu Reil, Seeker of Storms"
	c.Power = 5000
	c.Civ = civ.Light
	c.Family = []string{family.MechaThunder}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		if card.Zone != match.BATTLEZONE {
			return
		}

		fx.FilterShieldTriggers(ctx, func(x *match.Card) bool { return x.Civ != civ.Darkness })
	})

}

// ReBilSeekerOfArchery ...
func ReBilSeekerOfArchery(c *match.Card) {

	c.Name = "Re Bil, Seeker of Archery"
	c.Power = 6000
	c.Civ = civ.Light
	c.Family = []string{family.MechaThunder}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Doublebreaker, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {

			getLightCreatures(card, ctx2).Map(func(x *match.Card) {
				x.AddUniqueSourceCondition(cnd.PowerAmplifier, 2000, card.ID)
			})

			if card.Zone != match.BATTLEZONE {
				getLightCreatures(card, ctx2).Map(func(x *match.Card) {
					x.RemoveConditionBySource(card.ID)
				})
				exit()
			}

		})
	}))
}

func getLightCreatures(card *match.Card, ctx *match.Context) fx.CardCollection {

	lightCreatures := fx.FindFilter(
		card.Player,
		match.BATTLEZONE,
		func(x *match.Card) bool { return x.Civ == civ.Light && x.ID != card.ID },
	)

	lightCreatures = append(lightCreatures,

		fx.FindFilter(
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			func(x *match.Card) bool { return x.Civ == civ.Light && x.ID != card.ID },
		)...,
	)

	return lightCreatures
}
