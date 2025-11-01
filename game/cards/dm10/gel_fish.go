package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// BuoyantBlowfish ...
func BuoyantBlowfish(c *match.Card) {

	c.Name = "Buoyant Blowfish"
	c.Power = 1000
	c.Civ = civ.Water
	c.Family = []string{family.GelFish}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Water}

	c.PowerModifier = func(m *match.Match, attacking bool) int {
		power := 0

		fx.FindFilter(
			m.Opponent(c.Player),
			match.MANAZONE,
			func(x *match.Card) bool {
				return x.Tapped
			},
		).Map(func(x *match.Card) {
			power += 1000
		})

		return power
	}

	c.Use(fx.Creature)

}

// FluorogillManta
func FluorogillManta(c *match.Card) {

	c.Name = "Fluorogill Manta"
	c.Power = 1000
	c.Civ = civ.Water
	c.Family = []string{family.GelFish}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.When(fx.InTheBattlezone, func(card *match.Card, ctx *match.Context) {
		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
			if card.Zone != match.BATTLEZONE {
				fx.FindFilter(
					card.Player,
					match.BATTLEZONE,
					func(x *match.Card) bool {
						return x.HasCondition(cnd.Creature) && (x.Civ == civ.Light || x.Civ == civ.Darkness)
					},
				).Map(func(x *match.Card) {
					x.RemoveConditionBySource(card.ID)
				})

				exit()
				return
			}

			fx.FindFilter(
				card.Player,
				match.BATTLEZONE,
				func(x *match.Card) bool {
					return x.HasCondition(cnd.Creature) && (x.Civ == civ.Light || x.Civ == civ.Darkness)
				},
			).Map(func(x *match.Card) {
				x.AddUniqueSourceCondition(cnd.CantBeBlocked, true, card.ID)
			})
		})
	}))

}
