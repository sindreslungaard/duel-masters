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
	c.Family = family.MechaThunder
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		if card.Zone != match.BATTLEZONE {
			return
		}

		if event, ok := ctx.Event.(*match.ShieldTriggerEvent); ok {

			if event.Card.Civ == civ.Darkness {
				ctx.InterruptFlow()
			}

		}
	})

}

// ReBilSeekerOfArchery ...
func ReBilSeekerOfArchery(c *match.Card) {

	c.Name = "Re Bil, Seeker of Archery"
	c.Power = 6000
	c.Civ = civ.Light
	c.Family = family.MechaThunder
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Doublebreaker, func(card *match.Card, ctx *match.Context) {

		if card.Zone != match.BATTLEZONE {
			return
		}

		getLightCreatures(card, ctx).Map(func(x *match.Card) {
			x.AddUniqueSourceCondition(cnd.PowerAmplifier, 2000, card.ID)
		})
	})
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
