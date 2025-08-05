package fx

import (
	"duel-masters/game/match"
	"fmt"
)

func EachPlayerDestroys1Mana(card *match.Card, ctx *match.Context) {
	EachPlayerDestroysMana(card, ctx, 1)
}

func EachPlayerDestroysMana(card *match.Card, ctx *match.Context, quantity int) {

	players := make([]*match.Player, 0)
	players = append(players, card.Player)
	players = append(players, ctx.Match.Opponent(card.Player))

	for _, p := range players {

		cards := len(Find(p, match.MANAZONE))
		if quantity > cards {
			quantity = cards
		}

		Select(
			p,
			ctx.Match,
			p,
			match.MANAZONE,
			fmt.Sprintf("%s effect: Select %v card(s) from your manazone that will be sent to your graveyard", card.Name, quantity),
			quantity,
			quantity,
			false,
		).Map(func(manaCard *match.Card) {
			p.MoveCard(manaCard.ID, match.MANAZONE, match.GRAVEYARD, card.ID)
			ctx.Match.ReportActionInChat(p, fmt.Sprintf("%s effect: %s moved from MZ to GY", card.Name, manaCard.Name))
		})

	}

}

func DestroyOpCreature(card *match.Card, ctx *match.Context) {
	Select(
		card.Player,
		ctx.Match,
		ctx.Match.Opponent(card.Player),
		match.BATTLEZONE,
		"Destroy one of your opponent's creatures",
		1, 1, false,
	).Map(func(x *match.Card) {
		ctx.Match.Destroy(x, card, match.DestroyedBySpell)
	})
}

func DestroyYourself(card *match.Card, ctx *match.Context) {
	ctx.Match.Destroy(card, card, match.DestroyedByMiscAbility)
}

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

// Destroy creature less than equal to x
func DestroyOpCreatureXPowerOrLess(x int, destroyType match.CreatureDestroyedContext) match.HandlerFunc {
	return func(card *match.Card, ctx *match.Context) {
		SelectFilter(
			card.Player, ctx.Match, ctx.Match.Opponent(card.Player), match.BATTLEZONE,
			fmt.Sprintf("%s: Select 1 of your opponent's creatures that will be destroyed", card.Name),
			1, 1, false,
			func(creature *match.Card) bool { return ctx.Match.GetPower(creature, false) <= x }, false,
		).Map(func(x *match.Card) {
			ctx.Match.Destroy(x, card, destroyType)
		})
	}
}

// Destroy opponent creature by provided cancellable option and CreatureDestroyedContext
func DestroyOpponentCreature(cancellable bool, destroyType match.CreatureDestroyedContext) match.HandlerFunc {

	return func(card *match.Card, ctx *match.Context) {
		Select(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			"Destroy one of your opponent's creatures",
			1, 1, cancellable,
		).Map(func(x *match.Card) {
			ctx.Match.Destroy(x, card, destroyType)
		})
	}

}

func OpponentChoosesAndDestroysCreature(card *match.Card, ctx *match.Context) {
	Select(
		ctx.Match.Opponent(card.Player),
		ctx.Match,
		ctx.Match.Opponent(card.Player),
		match.BATTLEZONE,
		fmt.Sprintf("%s: Select 1 creature from your battlezone that will be sent to your graveyard", card.Name),
		1,
		1,
		false,
	).Map(func(x *match.Card) {
		ctx.Match.Destroy(x, card, match.DestroyedByMiscAbility)
	})
}

func OpponentChoosesManaBurn(card *match.Card, ctx *match.Context) {
	Select(
		ctx.Match.Opponent(card.Player),
		ctx.Match,
		ctx.Match.Opponent(card.Player),
		match.MANAZONE,
		fmt.Sprintf("%s effect: Select 1 card from your manazone that will be sent to your graveyard", card.Name),
		1,
		1,
		false,
	).Map(func(manaCard *match.Card) {
		ctx.Match.Opponent(card.Player).MoveCard(manaCard.ID, match.MANAZONE, match.GRAVEYARD, card.ID)
		ctx.Match.ReportActionInChat(ctx.Match.Opponent(card.Player), fmt.Sprintf("%s effect: %s moved from MZ to GY", card.Name, manaCard.Name))
	})
}
