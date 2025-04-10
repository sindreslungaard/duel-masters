package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// PoisonousMushroom ...
func PoisonousMushroom(c *match.Card) {

	c.Name = "Poisonous Mushroom"
	c.Power = 1000
	c.Civ = civ.Nature
	c.Family = []string{family.BalloonMushroom}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		selectedCards := fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.HAND,
			fmt.Sprintf("%s: You may select 1 card from your hand that will be sent to your manazone.", card.Name),
			1,
			1,
			true,
		)

		if len(selectedCards) > 0 {
			card.Player.MoveCard(selectedCards[0].ID, match.HAND, match.MANAZONE, card.ID)
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was moved to the manazone from your hand by %s's effect.", selectedCards[0].Name, card.Name))
		}

	}))

}
