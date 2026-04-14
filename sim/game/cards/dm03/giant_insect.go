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
				func(x *match.Card) bool { return event.Card.ID == x.ID },
				false,
			).Map(func(x *match.Card) {
				ctx.InterruptFlow()

				x.Player.MoveCard(x.ID, match.BATTLEZONE, match.MANAZONE, card.ID)
				ctx.Match.ReportActionInChat(x.Player, fmt.Sprintf("%s was moved to %s's mana by %s", x.Name, x.Player.Username(), card.Name))
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
			card.Player.MoveCard(x.ID, match.MANAZONE, match.HAND, card.ID)
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was sent to %s's hand from their manazone hand by %s's effect", x.Name, card.Player.Username(), card.Name))
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
