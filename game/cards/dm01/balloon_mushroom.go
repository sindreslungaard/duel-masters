package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
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

		hand, err := card.Player.Container(match.HAND)

		if err != nil {
			return
		}

		ctx.Match.NewAction(card.Player, hand, 1, 1, "Select 1 card from your hand that will be sent to your manazone. Choose close to cancel.", true)

		defer ctx.Match.CloseAction(card.Player)

		for {

			action := <-card.Player.Action

			if action.Cancel {
				break
			}

			if len(action.Cards) != 1 || !match.AssertCardsIn(hand, action.Cards...) {
				ctx.Match.DefaultActionWarning(card.Player)
				continue
			}

			card.Player.MoveCard(action.Cards[0], match.HAND, match.MANAZONE, card.ID)

			break

		}

	}))

}
