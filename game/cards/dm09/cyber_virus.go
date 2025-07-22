package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// MarchingMotherboard ...
func MarchingMotherboard(c *match.Card) {

	c.Name = "Marching Motherboard"
	c.Power = 2000
	c.Civ = civ.Water
	c.Family = []string{family.CyberVirus}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature,
		fx.When(fx.AnotherOwnCyberSummoned, func(card *match.Card, ctx *match.Context) {
			fx.MayDraw1(card, ctx)
		}))
}

// KelpCandle ...
func KelpCandle(c *match.Card) {

	c.Name = "Kelp Candle"
	c.Power = 1000
	c.Civ = civ.Water
	c.Family = []string{family.CyberVirus}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Blocker(), fx.CantAttackCreatures, fx.CantAttackPlayers,
		func(card *match.Card, ctx *match.Context) {
			if event, ok := ctx.Event.(*match.Battle); ok {
				if !event.Blocked || event.Defender != card {
					return
				}

				ctx.ScheduleAfter(func() {
					fx.LookTop4Put1IntoHandReorderRestOnBottomDeck(card, ctx)
				})
			}
		})
}
