package dm05

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// SmashHornQ ...
func SmashHornQ(c *match.Card) {

	c.Name = "Smash Horn Q"
	c.Power = 2000
	c.Civ = civ.Nature
	c.Family = []string{family.HornedBeast}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.Survivor, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {

			if card.Zone != match.BATTLEZONE {

				fx.FindFilter(
					card.Player,
					match.BATTLEZONE,
					func(x *match.Card) bool { return x.HasCondition(cnd.Survivor) },
				).Map(func(x *match.Card) {
					x.RemoveConditionBySource(card.ID)
				})

				exit()
				return

			}

			fx.FindFilter(
				card.Player,
				match.BATTLEZONE,
				func(x *match.Card) bool { return x.HasCondition(cnd.Survivor) },
			).Map(func(x *match.Card) {
				x.AddUniqueSourceCondition(cnd.PowerAmplifier, 1000, card.ID)
			})

		})

	}))

}

// MoonHorn ...
func MoonHorn(c *match.Card) {

	c.Name = "Moon Horn"
	c.Power = 6000
	c.Civ = civ.Nature
	c.Family = []string{family.HornedBeast}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Nature}

	c.PowerModifier = func(m *match.Match, attacking bool) int {

		cards := fx.FindFilter(m.Opponent(c.Player), match.BATTLEZONE, func(x *match.Card) bool { return x.Civ == civ.Water || x.Civ == civ.Darkness })

		return len(cards) * 1000

	}

	c.Use(fx.Creature, fx.Doublebreaker)

}
