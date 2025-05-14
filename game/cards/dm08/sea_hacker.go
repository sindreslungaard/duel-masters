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
		opp := ctx.Match.Opponent(card.Player)

		oppHandIds := make([]string, 0)

		fx.Find(opp, match.HAND).Map(func(x *match.Card) {
			oppHandIds = append(oppHandIds, x.ImageID)
		})

		ctx.Match.ShowCards(card.Player, "Your opponent's hand:", oppHandIds)

		peekedCards := opp.PeekDeck(1)

		if len(peekedCards) > 0 {
			oppTopDeckIds := make([]string, 0)
			oppTopDeckIds = append(oppTopDeckIds, peekedCards[0].ImageID)

			ctx.Match.ShowCards(card.Player, "Your opponent's top card of his deck:", oppTopDeckIds)
		}

	}))
}

// Vikorakys ...
func Vikorakys(c *match.Card) {

	c.Name = "Vikorakys"
	c.Power = 1000
	c.Civ = civ.Water
	c.Family = []string{family.SeaHacker}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Water}

	turboRush := false

	c.Use(
		fx.Creature,
		fx.When(fx.TurboRushCondition, func(card *match.Card, ctx *match.Context) { turboRush = true }),
		fx.When(fx.EndOfMyTurn, func(card *match.Card, ctx *match.Context) { turboRush = false }),
		fx.When(func(c *match.Card, ctx *match.Context) bool { return turboRush }, fx.When(fx.AttackConfirmed, fx.SearchDeckTakeXCards(1))),
	)

}
