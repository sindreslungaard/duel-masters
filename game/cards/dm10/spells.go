package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// StaticWarp ...
func StaticWarp(c *match.Card) {

	c.Name = "Static Warp"
	c.Civ = civ.Light
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Light}

	textFmt := "%s's effect: Choose a creature in your battle zone. Tap the rest of the creatures in the battle zone."

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		players := []*match.Player{card.Player, ctx.Match.Opponent(card.Player)}

		for _, p := range players {
			fx.Select(
				p,
				ctx.Match,
				p,
				match.BATTLEZONE,
				fmt.Sprintf(textFmt, card.Name),
				1,
				1,
				false,
			).Map(func(x *match.Card) {
				fx.FindFilter(
					p,
					match.BATTLEZONE,
					func(y *match.Card) bool {
						return y.ID != x.ID
					},
				).Map(func(y *match.Card) {
					y.Tapped = true
				})
			})
		}

		ctx.Match.BroadcastState()
	}))

}
