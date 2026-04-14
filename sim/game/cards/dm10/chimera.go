package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// Gigandura ...
func Gigandura(c *match.Card) {

	c.Name = "Gigandura"
	c.Power = 3000
	c.Civ = civ.Darkness
	c.Family = []string{family.Chimera}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		fx.Select(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.HAND,
			fmt.Sprintf("%s's effect: You may choose a card from your opponent's hand and put it into his mana zone. If you do, choose a card in his mana zone and put it into his hand.", card.Name),
			1,
			1,
			true,
		).Map(func(x *match.Card) {
			ctx.Match.MoveCard(x, match.MANAZONE, card)

			fx.Select(
				card.Player,
				ctx.Match,
				ctx.Match.Opponent(card.Player),
				match.MANAZONE,
				fmt.Sprintf("%s's effect: Choose a card in your opponent's mana zone and put it into his hand.", card.Name),
				1,
				1,
				false,
			).Map(func(x *match.Card) {
				ctx.Match.MoveCard(x, match.HAND, card)
			})
		})
	}))

}
