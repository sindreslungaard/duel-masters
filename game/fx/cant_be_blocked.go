package fx

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"fmt"
)

// CantBeBlocked allows the card to attack without being blocked
func CantBeBlocked(card *match.Card, ctx *match.Context) {

	if _, ok := ctx.Event.(*match.UntapStep); ok {

		card.AddCondition(cnd.CantBeBlocked, nil, card.ID)

	}

}

func GiveOwnCreatureCantBeBlocked(card *match.Card, ctx *match.Context) {
	Select(card.Player, ctx.Match, card.Player, match.BATTLEZONE,
		"Choose a card to receive 'Can't be blocked this turn'", 1, 1, false,
	).Map(func(x *match.Card) {
		x.AddCondition(cnd.CantBeBlocked, nil, card.ID)
		ctx.Match.ReportActionInChat(card.Player,
			fmt.Sprintf("%s tap effect: %s can't be blocked this turn", card.Name, x.Name))
	})
}

func CantBeBlockedByDarkness(card *match.Card, ctx *match.Context) {
	cantBeBlockedIf(card, ctx, func(blocker *match.Card) bool {
		return blocker.Civ == civ.Darkness
	})
}

func CantBeBlockedByLight(card *match.Card, ctx *match.Context) {
	cantBeBlockedIf(card, ctx, func(blocker *match.Card) bool {
		return blocker.Civ == civ.Light
	})
}

func CantBeBlockedByPower4000OrMore(card *match.Card, ctx *match.Context) {
	cantBeBlockedByPowerXOrMore(card, ctx, 4000)
}

func CantBeBlockedByPowerUpTo4000(card *match.Card, ctx *match.Context) {
	cantBeBlockedByPowerUpTo(card, ctx, 4000)
}

func CantBeBlockedByPowerUpTo5000(card *match.Card, ctx *match.Context) {
	cantBeBlockedByPowerUpTo(card, ctx, 5000)
}

func CantBeBlockedByPowerUpTo8000(card *match.Card, ctx *match.Context) {
	cantBeBlockedByPowerUpTo(card, ctx, 8000)
}

func CantBeBlockedByPowerUpTo3000(card *match.Card, ctx *match.Context) {
	cantBeBlockedByPowerUpTo(card, ctx, 3000)
}

func cantBeBlockedByPowerUpTo(card *match.Card, ctx *match.Context, power int) {
	cantBeBlockedIf(card, ctx, func(blocker *match.Card) bool {
		return ctx.Match.GetPower(blocker, false) > power
	})
}

func cantBeBlockedByPowerXOrMore(card *match.Card, ctx *match.Context, power int) {
	cantBeBlockedIf(card, ctx, func(blocker *match.Card) bool {
		return ctx.Match.GetPower(blocker, false) < power
	})
}

// CantBeBlockedIf
func cantBeBlockedIf(card *match.Card, ctx *match.Context, test func(blocker *match.Card) bool) {

	filter := func(blockers []*match.Card) []*match.Card {
		filtered := []*match.Card{}

		for _, b := range blockers {
			if !test(b) {
				filtered = append(filtered, b)
			}
		}

		return filtered
	}

	if event, ok := ctx.Event.(*match.SelectBlockers); ok &&
		event.CardWhoAttacked == card &&
		event.CardWhoAttacked.Zone == match.BATTLEZONE {

		ctx.ScheduleAfter(func() {
			event.Blockers = filter(event.Blockers)
		})

	}

}
