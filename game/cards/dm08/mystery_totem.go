package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// DracodanceTotem ...
func DracodanceTotem(c *match.Card) {

	c.Name = "Dracodance Totem"
	c.Power = 1000
	c.Civ = civ.Nature
	c.Family = []string{family.MysteryTotem}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature,
		fx.When(fx.WouldBeDestroyed, func(card *match.Card, ctx *match.Context) {
			dragonsInMyMana := fx.FindFilter(
				card.Player,
				match.MANAZONE,
				func(x *match.Card) bool {
					return x.SharesAFamily(family.Dragons)
				},
			)

			if len(dragonsInMyMana) > 0 {

				card.Player.MoveCard(card.ID, match.BATTLEZONE, match.MANAZONE, card.ID)
				ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was moved to manazone instead of being destroyed.", card.Name))

				fx.SelectFilter(
					card.Player,
					ctx.Match,
					card.Player,
					match.MANAZONE,
					fmt.Sprintf("%s's effect: Choose 1 Dragon from your manazone and put it into your hand.", card.Name),
					1,
					1,
					false,
					func(x *match.Card) bool {
						return x.SharesAFamily(family.Dragons)
					},
					false,
				).Map(func(x *match.Card) {
					card.Player.MoveCard(x.ID, match.MANAZONE, match.HAND, card.ID)
					ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was put into %s's hand from their manazone.", x.Name, card.Player.Username()))
				})

			}
		}),
	)
}
