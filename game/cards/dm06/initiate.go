package dm06

import (
	"fmt"

	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func CrazeValkyrieTheDrastic(c *match.Card) {

	c.Name = "Craze Valkyrie, the Drastic"
	c.Power = 7500
	c.Civ = civ.Light
	c.Family = []string{family.Initiate}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Evolution, fx.Doublebreaker, func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.CardMoved); ok {

			if event.CardID == card.ID && event.To == match.BATTLEZONE {

				creatures := match.Search(card.Player, ctx.Match, ctx.Match.Opponent(card.Player), match.BATTLEZONE, "Craze Valkyrie, the Drastic: Choose up to 2 of your opponent's creature and tap them. Close to not tap any creatures.", 1, 2, true)

				for _, creature := range creatures {
					creature.Tapped = true

					ctx.Match.Chat("Server", fmt.Sprintf("%s was tapped by %s's %s", creature.Name, card.Player.Username(), card.Name))
				}

			}

		}

	})

}
