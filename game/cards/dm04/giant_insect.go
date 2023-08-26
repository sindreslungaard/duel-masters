package dm04

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// ThreeEyedDragonfly ...
func ThreeEyedDragonfly(c *match.Card) {

	c.Name = "Three-Eyed Dragonfly"
	c.Power = 4000
	c.Civ = civ.Nature
	c.Family = []string{family.GiantInsect}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Nature}

	// id of the card to destroy that the user selects when attacking
	selectedCard := ""

	c.Use(
		fx.Creature,

		fx.When(fx.Attacking, func(card *match.Card, ctx *match.Context) {

			// reset this before each attack attempt
			selectedCard = ""

			cards := fx.SelectFilter(
				card.Player,
				ctx.Match,
				card.Player,
				match.BATTLEZONE,
				"Select one creature to destroy.",
				0,
				1,
				true,
				func(c *match.Card) bool { return c.ID != card.ID },
			)

			if len(cards) > 0 {
				selectedCard = cards[0].ID
				card.AddUniqueSourceCondition(cnd.DoubleBreaker, true, card.ID)
				card.AddCondition(cnd.PowerAmplifier, 2000, card.ID)
			}
		}),

		fx.When(fx.AttackConfirmed, func(card *match.Card, ctx *match.Context) {
			if selectedCard == "" {
				return
			}

			x, err := card.Player.GetCard(selectedCard, match.BATTLEZONE)

			if err != nil {
				card.RemoveConditionBySource(card.ID)
				return
			}

			ctx.Match.Destroy(x, card, match.DestroyedByMiscAbility)
			ctx.Match.Chat("Server", fmt.Sprintf("%s was destroyed by Three-Eyed Dragonfly", x.Name))
		}),
	)

}
