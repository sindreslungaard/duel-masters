package fx

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"fmt"
)

func EachPlayerDestroys1Mana(card *match.Card, ctx *match.Context) {
	eachPlayerDestroysMana(card, ctx, 1)
}

func EachPlayerDestroys2Mana(card *match.Card, ctx *match.Context) {
	eachPlayerDestroysMana(card, ctx, 2)
}

func eachPlayerDestroysMana(card *match.Card, ctx *match.Context, quantity int) {

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
		).Map(func(x *match.Card) {
			p.MoveCard(x.ID, match.MANAZONE, match.GRAVEYARD, card.ID)
			ctx.Match.ReportActionInChat(p, fmt.Sprintf("%s effect: %s moved from mana zone to graveyard", card.Name, x.Name))
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
// DestroyOpCreatureXPowerOrLess destroys one of the opponent's creatures whose
// effective non-attacking power is at most x. A cancellable selection models
// the printed "you may" wording; the prompt is never opened when no opposing
// creature is weak enough.
func DestroyOpCreatureXPowerOrLess(x int, cancellable bool, destroyType match.CreatureDestroyedContext) match.HandlerFunc {
	return func(card *match.Card, ctx *match.Context) {
		text := fmt.Sprintf("%s: Select 1 of your opponent's creatures that will be destroyed", card.Name)
		if cancellable {
			text = fmt.Sprintf("%s: You may select 1 of your opponent's creatures with power %d or less to destroy", card.Name, x)
		}

		SelectFilter(
			card.Player, ctx.Match, ctx.Match.Opponent(card.Player), match.BATTLEZONE,
			text,
			1, 1, cancellable,
			func(creature *match.Card) bool { return ctx.Match.GetPower(creature, false) <= x }, false,
		).Map(func(x *match.Card) {
			ctx.Match.Destroy(x, card, destroyType)
		})
	}
}

// DestroyAllCreaturesXPowerOrLess destroys every creature whose effective
// non-attacking power is at most x. The complete set is determined before any
// creature is destroyed so power changes caused by those destructions do not
// change which creatures the effect applies to.
func DestroyAllCreaturesXPowerOrLess(x int, destroyType match.CreatureDestroyedContext) match.HandlerFunc {
	return func(card *match.Card, ctx *match.Context) {
		toDestroy := make([]*match.Card, 0)

		players := []*match.Player{
			card.Player,
			ctx.Match.Opponent(card.Player),
		}

		for _, player := range players {
			FindFilter(
				player,
				match.BATTLEZONE,
				func(creature *match.Card) bool {
					return ctx.Match.GetPower(creature, false) <= x
				},
			).Map(func(creature *match.Card) {
				toDestroy = append(toDestroy, creature)
			})
		}

		for _, creature := range toDestroy {
			ctx.Match.Destroy(creature, card, destroyType)
		}
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

// DestroyOpBlocker destroys one of the opponent's creatures that has "blocker".
// The choice is mandatory, and no prompt is opened when the opponent controls no
// blocker at all.
func DestroyOpBlocker(card *match.Card, ctx *match.Context) {
	SelectFilter(
		card.Player,
		ctx.Match,
		ctx.Match.Opponent(card.Player),
		match.BATTLEZONE,
		fmt.Sprintf("%s's effect: Choose one of your opponent's creatures that has \"blocker\" and destroy it.", card.Name),
		1,
		1,
		false,
		func(x *match.Card) bool { return x.HasCondition(cnd.Blocker) },
		false,
	).Map(func(x *match.Card) {
		ctx.Match.Destroy(x, card, match.DestroyedByMiscAbility)
	})
}

// OwnChoosesAndDestroysCreature makes the card's controller choose one of their
// own creatures and destroy it.
func OwnChoosesAndDestroysCreature(card *match.Card, ctx *match.Context) {
	Select(
		card.Player,
		ctx.Match,
		card.Player,
		match.BATTLEZONE,
		fmt.Sprintf("%s: Select 1 creature from your battlezone that will be sent to your graveyard", card.Name),
		1,
		1,
		false,
	).Map(func(x *match.Card) {
		ctx.Match.Destroy(x, card, match.DestroyedByMiscAbility)
	})
}

// OwnChoosesManaBurn makes the card's controller choose a card in their own mana
// zone and put it into their graveyard.
func OwnChoosesManaBurn(card *match.Card, ctx *match.Context) {
	Select(
		card.Player,
		ctx.Match,
		card.Player,
		match.MANAZONE,
		fmt.Sprintf("%s effect: Select 1 card from your manazone that will be sent to your graveyard", card.Name),
		1,
		1,
		false,
	).Map(func(manaCard *match.Card) {
		card.Player.MoveCard(manaCard.ID, match.MANAZONE, match.GRAVEYARD, card.ID)
		ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s effect: %s moved from MZ to GY", card.Name, manaCard.Name))
	})
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

// DestroyUpToXOpBlockers destroys up to x of the opponent's creatures that have
// "blocker". The choice can be declined entirely, and no prompt is opened when
// the opponent controls no blocker at all.
func DestroyUpToXOpBlockers(x int) match.HandlerFunc {
	return func(card *match.Card, ctx *match.Context) {
		SelectFilter(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			fmt.Sprintf("%s's effect: Destroy up to %d of your opponent's creatures that have \"blocker\".", card.Name, x),
			1,
			x,
			true,
			func(creature *match.Card) bool { return creature.HasCondition(cnd.Blocker) },
			false,
		).Map(func(blocker *match.Card) {
			ctx.Match.Destroy(blocker, card, match.DestroyedByMiscAbility)
		})
	}
}

// OpponentChoosesCreatureToMana makes the opponent choose one of their own
// creatures and put it into their mana zone.
func OpponentChoosesCreatureToMana(card *match.Card, ctx *match.Context) {
	opponent := ctx.Match.Opponent(card.Player)

	Select(
		opponent,
		ctx.Match,
		opponent,
		match.BATTLEZONE,
		fmt.Sprintf("%s: Choose one of your creatures and put it into your mana zone", card.Name),
		1,
		1,
		false,
	).Map(func(creature *match.Card) {
		moved, err := opponent.MoveCard(creature.ID, match.BATTLEZONE, match.MANAZONE, card.ID)

		if err != nil || moved.Zone != match.MANAZONE {
			return
		}

		ctx.Match.ReportActionInChat(opponent, fmt.Sprintf("%s was put into %s's mana zone by %s", creature.Name, opponent.Username(), card.Name))
	})
}

// PlayerChoosesAndDestroysOwnCreature makes the given player pick one of their
// own creatures and destroy it. Unlike OwnChoosesAndDestroysCreature the chooser
// is passed in, which is what an effect that fires on both players' turns needs.
func PlayerChoosesAndDestroysOwnCreature(player *match.Player, card *match.Card, ctx *match.Context) {
	Select(
		player,
		ctx.Match,
		player,
		match.BATTLEZONE,
		fmt.Sprintf("%s: Choose one of your creatures and destroy it.", card.Name),
		1,
		1,
		false,
	).Map(func(creature *match.Card) {
		ctx.Match.Destroy(creature, card, match.DestroyedByMiscAbility)
	})
}

// DestroyAllCreatures destroys every creature in the battle zone, on both
// sides, sparing nothing. The full set is decided before anything is destroyed,
// so a destruction that changes the board cannot change who is caught.
func DestroyAllCreatures(destroyType match.CreatureDestroyedContext) match.HandlerFunc {
	return func(card *match.Card, ctx *match.Context) {
		toDestroy := make([]*match.Card, 0)

		for _, player := range []*match.Player{card.Player, ctx.Match.Opponent(card.Player)} {
			Find(player, match.BATTLEZONE).Map(func(creature *match.Card) {
				toDestroy = append(toDestroy, creature)
			})
		}

		for _, creature := range toDestroy {
			ctx.Match.Destroy(creature, card, destroyType)
		}
	}
}

// DestroyLowestPowerCreature destroys the creature with the lowest power in the
// battle zone, counting both players. A tie is broken by the card's controller
// choosing from among the tied creatures.
func DestroyLowestPowerCreature(card *match.Card, ctx *match.Context) {
	candidates := make([]*match.Card, 0)

	for _, player := range []*match.Player{card.Player, ctx.Match.Opponent(card.Player)} {
		Find(player, match.BATTLEZONE).Map(func(creature *match.Card) {
			candidates = append(candidates, creature)
		})
	}

	if len(candidates) < 1 {
		return
	}

	lowest := ctx.Match.GetPower(candidates[0], false)
	for _, creature := range candidates[1:] {
		if power := ctx.Match.GetPower(creature, false); power < lowest {
			lowest = power
		}
	}

	tied := make([]*match.Card, 0, len(candidates))
	for _, creature := range candidates {
		if ctx.Match.GetPower(creature, false) == lowest {
			tied = append(tied, creature)
		}
	}

	if len(tied) == 1 {
		ctx.Match.Destroy(tied[0], card, match.DestroyedByMiscAbility)
		return
	}

	// The tie is broken by this card's controller, and the tied creatures can
	// belong to either player, so they are offered as one group.
	byID := make(map[string]bool, len(tied))
	for _, creature := range tied {
		byID[creature.ID] = true
	}

	choices := map[string][]*match.Card{
		"Your creatures":            FindFilter(card.Player, match.BATTLEZONE, func(x *match.Card) bool { return byID[x.ID] }),
		"Your opponent's creatures": FindFilter(ctx.Match.Opponent(card.Player), match.BATTLEZONE, func(x *match.Card) bool { return byID[x.ID] }),
	}

	SelectMultipart(
		card.Player,
		ctx.Match,
		choices,
		fmt.Sprintf("%s's effect: Several creatures are tied for the lowest power. Choose one to destroy.", card.Name),
		1,
		1,
		false,
	).Map(func(chosen *match.Card) {
		ctx.Match.Destroy(chosen, card, match.DestroyedByMiscAbility)
	})
}

// DestroyAllCreaturesWithExactPower destroys every creature whose effective
// power is exactly the given figure, on both sides. The set is decided before
// anything is destroyed.
func DestroyAllCreaturesWithExactPower(power int, destroyType match.CreatureDestroyedContext) func(*match.Card, *match.Context) {
	return func(card *match.Card, ctx *match.Context) {
		toDestroy := make([]*match.Card, 0)

		for _, player := range []*match.Player{card.Player, ctx.Match.Opponent(card.Player)} {
			FindFilter(player, match.BATTLEZONE, func(x *match.Card) bool {
				return ctx.Match.GetPower(x, false) == power
			}).Map(func(creature *match.Card) {
				toDestroy = append(toDestroy, creature)
			})
		}

		for _, creature := range toDestroy {
			ctx.Match.Destroy(creature, card, destroyType)
		}
	}
}
