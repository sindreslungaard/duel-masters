package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// MicuteTheOracle ...
func MicuteTheOracle(c *match.Card) {

	c.Name = "Micute, the Oracle"
	c.Power = 4000
	c.Civ = civ.Light
	c.Family = []string{family.LightBringer}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature,
		fx.When(fx.AnotherOwnGuardianSummoned, func(card *match.Card, ctx *match.Context) {
			fx.Select(
				card.Player,
				ctx.Match,
				ctx.Match.Opponent(card.Player),
				match.BATTLEZONE,
				fmt.Sprintf("%s's effect: You may choose one of your opponent's creatures in the battlezone and tap it.", card.Name),
				1,
				1,
				true,
			).Map(func(x *match.Card) {
				x.Tapped = true
				ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was tapped by %s's effect.", x.Name, card.Name))
			})
		}))
}
