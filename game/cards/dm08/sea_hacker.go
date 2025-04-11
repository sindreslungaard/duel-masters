package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// Lalicious ...
func Lalicious(c *match.Card) {

	c.Name = "Lalicious"
	c.Power = 4000
	c.Civ = civ.Water
	c.Family = []string{family.SeaHacker}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.When(fx.AttackConfirmed, func(card *match.Card, ctx *match.Context) {

		oppHandIds := make([]string, 0)

		fx.Find(ctx.Match.Opponent(card.Player), match.HAND).Map(func(x *match.Card) {
			oppHandIds = append(oppHandIds, x.ImageID)
		})

		ctx.Match.ShowCards(card.Player, "Your opponent's hand:", oppHandIds)

		oppTopDeckIds := make([]string, 0)

		oppTopDeckIds = append(oppTopDeckIds, ctx.Match.Opponent(card.Player).PeekDeck(1)[0].ImageID)

		ctx.Match.ShowCards(card.Player, "Your opponent's top card of his deck:", oppTopDeckIds)

	}))
}
