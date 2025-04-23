package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// KyrstronLairDelver ...
func KyrstronLairDelver(c *match.Card) {

	c.Name = "Kyrstron, Lair Delver"
	c.Power = 1000
	c.Civ = civ.Fire
	c.Family = []string{family.Dragonoid}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.When(fx.Destroyed, func(card *match.Card, ctx *match.Context) {
		fx.SelectFilter(
			card.Player,
			ctx.Match,
			card.Player,
			match.HAND,
			fmt.Sprintf("%s: You may put 1 of your Dragons from your hand into the battlezone.", card.Name),
			1,
			1,
			true,
			func(c *match.Card) bool {
				return c.SharesAFamily(family.Dragons)
			},
			false,
		).Map(func(x *match.Card) {
			ctx.Match.MoveCard(x, match.BATTLEZONE, card)
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s's effect: %s was put into the battlezone", card.Name, x.Name))
		})
	}))

}
