package dm03

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// Gigamantis ...
func Gigamantis(c *match.Card) {

	c.Name = "Gigamantis"
	c.Power = 5000
	c.Civ = civ.Nature
	c.Family = []string{family.GiantInsect}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.Evolution, func(card *match.Card, ctx *match.Context) {

		if card.Zone != match.BATTLEZONE {
			return
		}

		if event, ok := ctx.Event.(*match.CreatureDestroyed); ok &&
			event.Card.ID != card.ID &&
			event.Card.Player == card.Player &&
			event.Card.Civ == civ.Nature {

			fx.SelectFilter(
				card.Player,
				ctx.Match,
				card.Player,
				match.BATTLEZONE,
				fmt.Sprintf("%s: You may put card to manazone.", card.Name),
				1,
				1,
				true,
				func(c *match.Card) bool { return event.Card.ID == c.ID },
			).Map(func(c *match.Card) {

				ctx.InterruptFlow()

				c.Player.MoveCard(c.ID, match.BATTLEZONE, match.MANAZONE)
				ctx.Match.Chat("Server", fmt.Sprintf("%s was moved to %s's mana by %s", c.Name, c.Player.Username(), card.Name))
			})

		}

	})

}

// SniperMosquito ...
func SniperMosquito(c *match.Card) {

	c.Name = "Sniper Mosquito"
	c.Power = 2000
	c.Civ = civ.Nature
	c.Family = []string{family.GiantInsect}
	c.ManaCost = 1
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.When(fx.AttackConfirmed, func(card *match.Card, ctx *match.Context) {

		fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.MANAZONE,
			"Select a card from your mana zone to be returned to your hand",
			1,
			1,
			false,
		).Map(func(x *match.Card) {
			ctx.Match.MoveCard(x, match.HAND, card)
		})
	}))
}

// SwordButterfly ...
func SwordButterfly(c *match.Card) {

	c.Name = "Sword Butterfly"
	c.Power = 2000
	c.Civ = civ.Nature
	c.Family = []string{family.GiantInsect}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.PowerAttacker3000)
}
