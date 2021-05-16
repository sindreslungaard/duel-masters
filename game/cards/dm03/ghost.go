package dm03

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// JackViperShadowofDoom ...
func JackViperShadowofDoom(c *match.Card) {

	c.Name = "Jack Viper, Shadow of Doom"
	c.Power = 4000
	c.Civ = civ.Darkness
	c.Family = family.Ghost
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Evolution, func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.CreatureDestroyed); ok {

			ctx.ScheduleAfter(func() {

				destroyedCard, err := card.Player.GetCard(event.Card.ID, match.GRAVEYARD)

				if err != nil {
					return
				}

				if event.Card.ID != card.ID && destroyedCard.Civ == civ.Darkness {

					card.Player.MoveCard(event.Card.ID, match.GRAVEYARD, match.HAND)
					ctx.Match.Chat("Server", fmt.Sprintf("%s was moved to %s's hand from the graveyard", destroyedCard.Name, ctx.Match.PlayerRef(card.Player).Socket.User.Username))
				}

			})

		}
	})

}

// WailingShadowBelbetphlo ...
func WailingShadowBelbetphlo(c *match.Card) {

	c.Name = "Wailing Shadow Belbetphlo"
	c.Power = 1000
	c.Civ = civ.Darkness
	c.Family = family.Ghost
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Slayer)

}
