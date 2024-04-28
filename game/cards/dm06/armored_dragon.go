package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

func BazagazealDragon(c *match.Card) {
	c.Name = "Bazagazeal Dragon"
	c.Power = 8000
	c.Civ = civ.Fire
	c.Family = []string{family.ArmoredDragon}
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.Doublebreaker, fx.SpeedAttacker, fx.AttackUntapped, fx.When(fx.EndOfMyTurn, func(card *match.Card, ctx *match.Context) {
		if card.Zone != match.BATTLEZONE {
			return
		}

		card.Player.MoveCard(card.ID, match.BATTLEZONE, match.HAND, card.ID)
		ctx.Match.Chat("Server", fmt.Sprintf("%s was returned to the %s's hand", c.Name, c.Player.Username()))
	}))

}

func BolmeteusSteelDragon(c *match.Card) {
	c.Name = "Bolmeteus Steel Dragon"
	c.Power = 7000
	c.Civ = civ.Fire
	c.Family = []string{family.ArmoredDragon}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.Doublebreaker, func(card *match.Card, ctx *match.Context) {

		if card.Zone != match.BATTLEZONE {
			return
		}

		if event, ok := ctx.Event.(*match.ShieldTriggerEvent); ok && event.Source == card.ID {

			ctx.InterruptFlow()

		}

		if event, ok := ctx.Event.(*match.BrokenShieldEvent); ok && event.Source == card.ID {

			ctx.InterruptFlow()

			ctx.Match.Opponent(card.Player).MoveCard(event.CardID, match.HAND, match.GRAVEYARD, card.ID)

		}
	})
}
