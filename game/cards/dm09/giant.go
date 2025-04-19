package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// StratosphereGiant ...
func StratosphereGiant(c *match.Card) {

	c.Name = "Stratosphere Giant"
	c.Power = 6000
	c.Civ = civ.Nature
	c.Family = []string{family.Giant}
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.Triplebreaker, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		fx.SelectFilter(
			ctx.Match.Opponent(card.Player),
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.HAND,
			"Choose up to 2 creatures in your hand and put them into the battlezone.",
			1,
			2,
			true,
			func(c *match.Card) bool {
				return fx.CanBeSummoned(ctx.Match.Opponent(card.Player), c)
			},
			false,
		).Map(func(x *match.Card) {
			cardPlayedCtx := match.NewContext(ctx.Match, &match.CardPlayedEvent{
				CardID: x.ID,
			})
			ctx.Match.HandleFx(cardPlayedCtx)

			if !cardPlayedCtx.Cancelled() {

				if !x.HasCondition(cnd.Evolution) {
					x.AddCondition(cnd.SummoningSickness, nil, nil)
				}

				ctx.Match.Opponent(card.Player).MoveCard(x.ID, match.HAND, match.BATTLEZONE, x.ID)
				ctx.Match.ReportActionInChat(ctx.Match.Opponent(card.Player), fmt.Sprintf("%s was moved to the battle zone by %s's effect", x.Name, card.Name))

			}
		})
	}))
}
