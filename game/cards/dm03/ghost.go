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

	c.Use(fx.Creature, fx.Evolution, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {

			if card.Zone != match.BATTLEZONE {
				exit()
				return
			}

			event, ok := ctx2.Event.(*match.CardMoved)

			if !ok || event.From != match.BATTLEZONE {
				return
			}

			fx.FindFilter(
				card.Player,
				match.GRAVEYARD,
				func(x *match.Card) bool { return x.ID == event.CardID && x.ID != card.ID && card.Civ == civ.Darkness },
			).Map(func(x *match.Card) {
				card.Player.MoveCard(x.ID, match.GRAVEYARD, match.HAND)
				ctx.Match.Chat("Server", fmt.Sprintf("%s was moved to %s's hand from the graveyard by Jack Viper, Shadow of Doom", x.Name, x.Player.Username()))
			})

		})

	}))

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
