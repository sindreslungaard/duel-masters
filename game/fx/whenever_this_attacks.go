package fx

import (
	"duel-masters/game/civ"
	"duel-masters/game/match"
	"fmt"
)

func WheneverThisAttacksMayTapDorFCreature() match.HandlerFunc {
	return When(AttackConfirmed, func(c *match.Card, ctx *match.Context) {
		filter := func(x *match.Card) bool { return x.Civ == civ.Fire || x.Civ == civ.Darkness }
		cards := make(map[string][]*match.Card)
		cards["Your creatures"] = FindFilter(c.Player, match.BATTLEZONE, filter)
		cards["Opponent's creatures"] = FindFilter(ctx.Match.Opponent(c.Player), match.BATTLEZONE, filter)

		SelectMultipart(
			c.Player,
			ctx.Match,
			cards,
			fmt.Sprintf("%s: You may select a darkness or fire creature in the battlezone to tap", c.Name),
			1,
			1,
			true,
		).Map(func(x *match.Card) {
			x.Tapped = true
		})

	})
}

func WheneverThisAttacksMayLookAtOpShield() match.HandlerFunc {
	return When(AttackConfirmed, func(card *match.Card, ctx *match.Context) {
		SelectBackside(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.SHIELDZONE,
			fmt.Sprintf("%s: You may select 1 of your opponent's shields that will be shown to you", card.Name),
			1,
			1,
			true,
		).Map(func(x *match.Card) {
			ctx.Match.ShowCards(
				card.Player,
				"Your opponent's shield:",
				[]string{x.ImageID},
			)
		})
	})
}
