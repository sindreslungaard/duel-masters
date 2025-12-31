package dm04

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// MissileBoy ...
func MissileBoy(c *match.Card) {

	c.Name = "Missile Boy"
	c.Power = 1000
	c.Civ = civ.Fire
	c.Family = []string{family.Human}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Fire}

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

			if playedCard.Civ != civ.Light {
				return
			}

			playedCard.AddUniqueSourceCondition(cnd.IncreasedCost, 1, card.ID)
		}

	})

}