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

		if card.Zone != match.BATTLEZONE {
			return
		}

		if event, ok := ctx.Event.(*match.CreatureDestroyed); ok {

			if event.Card.ID != c.ID && card.Civ == civ.Darkness && event.Card.Player == card.Player {

				card.Player.MoveCard(event.Card.ID, match.BATTLEZONE, match.HAND)
				ctx.Match.Chat("Server", fmt.Sprintf("Jack Viper, Shadow of Doom: %s was destroyed by %s and returned to the hand", event.Card.Name, event.Source.Name))

				ctx.InterruptFlow()
			}
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
