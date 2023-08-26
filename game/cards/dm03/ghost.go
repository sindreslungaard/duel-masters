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
	c.Family = []string{family.Ghost}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Evolution, func(card *match.Card, ctx *match.Context) {

		if card.Zone != match.BATTLEZONE {
			return
		}

		if event, ok := ctx.Event.(*match.CreatureDestroyed); ok &&
			event.Card.ID != card.ID &&
			event.Card.Player == card.Player &&
			event.Card.Civ == civ.Darkness {

			fx.SelectFilter(
				card.Player,
				ctx.Match,
				card.Player,
				match.BATTLEZONE,
				fmt.Sprintf("%s: You may return card to hand.", card.Name),
				1,
				1,
				true,
				func(c *match.Card) bool { return event.Card.ID == c.ID },
			).Map(func(c *match.Card) {

				ctx.InterruptFlow()

				c.Player.MoveCard(c.ID, match.BATTLEZONE, match.HAND)
				ctx.Match.Chat("Server", fmt.Sprintf("%s was moved to %s's hand by %s", c.Name, c.Player.Username(), card.Name))
			})

		}

	})

}

// WailingShadowBelbetphlo ...
func WailingShadowBelbetphlo(c *match.Card) {

	c.Name = "Wailing Shadow Belbetphlo"
	c.Power = 1000
	c.Civ = civ.Darkness
	c.Family = []string{family.Ghost}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Slayer)

}
