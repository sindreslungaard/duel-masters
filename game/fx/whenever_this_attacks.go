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

func WheneverOneOfMyCreaturesAttacksOppDiscardsRandom() match.HandlerFunc {
	return When(OneOfMyCreaturesAttacksConfirmed, OpponentDiscardsRandomCard)
}

func LookAtOppShields(card *match.Card, ctx *match.Context) {
	ctx.Match.ShowCards(
		card.Player,
		"Your opponent's shield:",
		Find(
			ctx.Match.Opponent(card.Player),
			match.SHIELDZONE,
		).ProjectImageIDs(),
	)
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
			if ev, ok := ctx.Event.(*match.AttackConfirmed); ok && ev.Player {
				ctx.Match.ShowCardsNonDismissible(
					card.Player,
					"Your opponent's shield:",
					[]string{x.ImageID},
				)
			} else {
				ctx.Match.ShowCards(
					card.Player,
					"Your opponent's shield:",
					[]string{x.ImageID},
				)
			}
		})
	})
}
