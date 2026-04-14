package dm05

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// HornedMutant ...
func HornedMutant(c *match.Card) {

	c.Name = "Horned Mutant"
	c.Power = 3000
	c.Civ = civ.Darkness
	c.Family = []string{family.Hedrian}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		if card.Zone != match.BATTLEZONE {
			return
		}

		if event, ok := ctx.Event.(*match.PlayCardEvent); ok {

			p := ctx.Match.CurrentPlayer()

			playedCard, err := p.Player.GetCard(event.CardID, match.HAND)

			if err != nil {
				return
			}

			if playedCard.Civ != civ.Nature {
				return
			}

			playedCard.AddUniqueSourceCondition(cnd.IncreasedCost, 1, card.ID)
		}
	})

}
