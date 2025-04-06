package dm03

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// ArmoredWarriorQuelos ...
func ArmoredWarriorQuelos(c *match.Card) {

	c.Name = "Armored Warrior Quelos"
	c.Power = 2000
	c.Civ = civ.Fire
	c.Family = []string{family.Armorloid}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.When(fx.AttackConfirmed, func(card *match.Card, ctx *match.Context) {

		creatures := fx.SelectFilter(
			card.Player,
			ctx.Match,
			card.Player,
			match.MANAZONE,
			fmt.Sprintf("%s: Select 1 of your non-fire cards from mana zone and put it into your graveyard.", card.Name),
			1,
			1,
			false,
			func(x *match.Card) bool { return x.Civ != civ.Fire },
			false,
		)

		for _, creature := range creatures {
			card.Player.MoveCard(creature.ID, match.MANAZONE, match.GRAVEYARD, card.ID)
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was sent to graveyard from %s's mana zone by %s's effect.", creature.Name, card.Player.Username(), card.Name))
		}

		ctx.Match.Wait(card.Player, "Waiting for your opponent to make an action")
		defer ctx.Match.EndWait(card.Player)

		opponentCreatures := fx.SelectFilter(
			ctx.Match.Opponent(card.Player),
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.MANAZONE,
			fmt.Sprintf("%s: Select 1 of your non-fire cards from mana zone and put it into your graveyard.", card.Name),
			1,
			1,
			false,
			func(x *match.Card) bool { return x.Civ != civ.Fire },
			false,
		)

		for _, creature := range opponentCreatures {
			ctx.Match.Opponent(card.Player).MoveCard(creature.ID, match.MANAZONE, match.GRAVEYARD, card.ID)
			ctx.Match.ReportActionInChat(ctx.Match.Opponent(card.Player), fmt.Sprintf("%s was sent to graveyard from %s's mana zone by %s's effect.", creature.Name, ctx.Match.Opponent(card.Player).Username(), card.Name))
		}

	}))

}
