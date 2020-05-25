package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// StingerWorm ...
func StingerWorm(c *match.Card) {

	c.Name = "Stinger Worm"
	c.Power = 5000
	c.Civ = civ.Darkness
	c.Family = family.BeastFolk
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.CardMoved); ok {

			if event.CardID == card.ID && event.To == match.BATTLEZONE {

				creatures := match.Search(card.Player, ctx.Match, card.Player, match.BATTLEZONE, "Stinger Worm: Select 1 creature from your playzone that will be sent to your graveyard", 1, 1, false)

				for _, creature := range creatures {
					ctx.Match.Destroy(creature, card)
				}

			}

		}

	})

}
