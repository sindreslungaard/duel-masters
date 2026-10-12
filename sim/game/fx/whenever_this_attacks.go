package fx

import (
	"duel-masters/game/civ"
	"duel-masters/game/match"
	"fmt"
)

func WheneverThisAttacksReturnCardFromMZToHand() match.HandlerFunc {
	return When(AttackConfirmed, func(c *match.Card, ctx *match.Context) {
		ReturnMyCardFromMZToHand(c, ctx)
	})
}

func WheneverThisAttacksMayTapDorFCreature() match.HandlerFunc {
	return When(AttackConfirmed, func(c *match.Card, ctx *match.Context) {
		filter := func(x *match.Card) bool { return x.HasCiv(civ.Fire) || x.HasCiv(civ.Darkness) }
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
			TapCreature(x, ctx, c)
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
		opponent := ctx.Match.Opponent(card.Player)
		chooser := ResolveShieldChooser(ctx, card.Player, opponent)

		SelectBackside(
			chooser,
			ctx.Match,
			opponent,
			match.SHIELDZONE,
			fmt.Sprintf("%s: You may select 1 of %s shields that will be shown to you%s", card.Name, ShieldPossessive(chooser, opponent), MeloppeNote(chooser, card.Player)),
			1,
			1,
			true,
		).Map(func(x *match.Card) {
			// Grabbed before the pop-up so the number still matches: whoever
			// picked x (its own owner, if Meloppe redirected) never sees this
			// reveal, so it's the only way card.Player can tell them which
			// shield it was.
			description := DescribeShield(x)

			// Reported before the pop-up, which can block on the attacker
			// dismissing it: the defender who picked shouldn't have to wait on
			// that to find out what happened.
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s of %s was shown to %s", description, opponent.Username(), card.Player.Username()))

			if ev, ok := ctx.Event.(*match.AttackConfirmed); ok && ev.Player {
				ctx.Match.ShowCardsNonDismissible(
					card.Player,
					fmt.Sprintf("Your opponent's %s:", description),
					[]string{x.ImageID},
				)
			} else {
				ctx.Match.ShowCards(
					card.Player,
					fmt.Sprintf("Your opponent's %s:", description),
					[]string{x.ImageID},
				)
			}
		})
	})
}
