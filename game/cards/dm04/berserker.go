package dm04

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// MilieusTheDaystretcher ...
func MilieusTheDaystretcher(c *match.Card) {

	c.Name = "Milieus, the Daystretcher"
	c.Power = 2500
	c.Civ = civ.Light
	c.Family = []string{family.Berserker}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Blocker, func(card *match.Card, ctx *match.Context) {

		if card.Zone != match.BATTLEZONE {
			return
		}

		if event, ok := ctx.Event.(*match.PlayCardEvent); ok {

			p := ctx.Match.CurrentPlayer()

			playedCard, err := p.Player.GetCard(event.CardID, match.HAND)

			if err != nil {
				return
			}

			if playedCard.Civ != civ.Darkness {
				return
			}

			playedCard.AddUniqueSourceCondition(cnd.IncreasedCost, 2, card.ID)
		}
	})
}
