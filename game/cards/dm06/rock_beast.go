package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func RumblesaurQ(c *match.Card) {
	c.Name = "Rumblesaur Q"
	c.Power = 3000
	c.Civ = civ.Fire
	c.Family = []string{family.RockBeast, family.Survivor}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.SpeedAttacker, fx.Survivor)
	c.Use(fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		fx.FindFilter(
			card.Player,
			match.BATTLEZONE,
			func(x *match.Card) bool { return x.HasCondition(cnd.Survivor) },
		).Map(func(x *match.Card) {
			x.RemoveCondition(cnd.SummoningSickness)
		})
	}))

	c.Use(func(card *match.Card, ctx *match.Context) {
		if card.Zone != match.BATTLEZONE {
			return
		}

		if event, ok := ctx.Event.(*match.CardMoved); ok && event.To == match.BATTLEZONE {

			fx.FindFilter(
				card.Player,
				match.BATTLEZONE,
				func(x *match.Card) bool { return x.HasCondition(cnd.Survivor) },
			).Map(func(x *match.Card) {
				x.RemoveCondition(cnd.SummoningSickness)
			})

		}
	})
}
