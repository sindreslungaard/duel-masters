package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// KaluteVizierOfEternity ...
func KaluteVizierOfEternity(c *match.Card) {

	c.Name = "Kalute, Vizier of Eternity"
	c.Power = 1000
	c.Civ = civ.Light
	c.Family = []string{family.Initiate}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature,
		func(card *match.Card, ctx *match.Context) {
			if event, ok := ctx.Event.(*match.CreatureDestroyed); ok {
				if event.Card.ID == card.ID {
					if len(fx.FindFilter(
						card.Player,
						match.BATTLEZONE,
						func(x *match.Card) bool {
							return x.Name == card.Name && x.ID != card.ID
						},
					)) > 0 {
						ctx.InterruptFlow()

						_, err := card.Player.MoveCard(card.ID, match.BATTLEZONE, match.HAND, card.ID)

						if err != nil {
							ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was moved to hand instead of being destroyed.", card.Name))
						}
					}
				}
			}
		})
}
