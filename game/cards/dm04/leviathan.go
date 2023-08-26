package dm04

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// KingAquakamui ...
func KingAquakamui(c *match.Card) {

	c.Name = "King Aquakamui"
	c.Power = 5000
	c.Civ = civ.Water
	c.Family = []string{family.Leviathan}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.CardMoved); ok {

			if event.CardID == card.ID && event.To == match.BATTLEZONE {

				fx.FindFilter(
					card.Player,
					match.GRAVEYARD,
					func(x *match.Card) bool { return x.HasFamily(family.AngelCommand) || x.HasFamily(family.DemonCommand) },
				).Map(func(x *match.Card) {

					x.Player.MoveCard(x.ID, match.GRAVEYARD, match.HAND)
					ctx.Match.Chat("Server", fmt.Sprintf("%s was moved to %s's hand from their graveyard by King Aquakamui", x.Name, card.Player.Username()))
				})
			}
		}

		if card.Zone != match.BATTLEZONE {
			return
		}

		if event, ok := ctx.Event.(*match.GetPowerEvent); ok {

			if event.Card.HasFamily(family.AngelCommand) || event.Card.HasFamily(family.DemonCommand) {
				event.Power += 2000
			}
		}
	})
}
