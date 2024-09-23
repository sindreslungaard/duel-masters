package fx

import (
	"duel-masters/game/match"
	"fmt"
)

func destroyOpCreature2000OrLess(card *match.Card, ctx *match.Context, destroyType match.CreatureDestroyedContext) {
	SelectFilter(
		card.Player, ctx.Match, ctx.Match.Opponent(card.Player), match.BATTLEZONE,
		fmt.Sprintf("%s: Select 1 of your opponent's creatures that will be destroyed", card.Name),
		1, 1, false,
		func(x *match.Card) bool { return ctx.Match.GetPower(x, false) <= 2000 }, false,
	).Map(func(x *match.Card) {
		ctx.Match.Destroy(x, card, destroyType)
	})
}

func DestroyBySpellOpCreature2000OrLess(card *match.Card, ctx *match.Context) {
	destroyOpCreature2000OrLess(card, ctx, match.DestroyedBySpell)
}

func DestroyByMiscOpCreature2000OrLess(card *match.Card, ctx *match.Context) {
	destroyOpCreature2000OrLess(card, ctx, match.DestroyedByMiscAbility)
}
