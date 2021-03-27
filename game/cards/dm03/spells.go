package dm03

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// BoomerangComet ...
func BoomerangComet(c *match.Card) {

	c.Name = "Boomerang Comet"
	c.Civ = civ.Light
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.ShieldTrigger, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			cards := match.Filter(card.Player, ctx.Match, card.Player, match.MANAZONE, "Select 1 card from your mana zone that will be sent to your hand", 1, 1, false, func(x *match.Card) bool { return true })

			for _, card := range cards {

				card.Player.MoveCard(card.ID, match.MANAZONE, match.HAND)
				ctx.Match.Chat("Server", fmt.Sprintf("%s retrieved %s from the mana zone to their hand", card.Player.Username(), card.Name))

			}

			card.Player.MoveCard(card.ID, match.HAND, match.MANAZONE)
		}
	})
}

// LogicSphere ...
func LogicSphere(c *match.Card) {

	c.Name = "Logic Sphere"
	c.Civ = civ.Light
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.ShieldTrigger, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			spells := match.Filter(card.Player, ctx.Match, card.Player, match.MANAZONE, "Select 1 spell from your mana zone that will be sent to your hand", 1, 1, false, func(x *match.Card) bool { return x.HasCondition(cnd.Spell) })

			for _, spell := range spells {

				card.Player.MoveCard(spell.ID, match.MANAZONE, match.HAND)
				ctx.Match.Chat("Server", fmt.Sprintf("%s retrieved %s from the mana zone to their hand", spell.Player.Username(), spell.Name))

			}
		}
	})
}

// SundropArmor ...
func SundropArmor(c *match.Card) {

	c.Name = "Sundrop Armor"
	c.Civ = civ.Light
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			fx.SelectFilter(
				card.Player,
				ctx.Match,
				card.Player,
				match.HAND,
				"Select 1 card from your hand that will be put as a shield",
				1,
				1,
				false,
				func(c *match.Card) bool { return c.ID != card.ID },
			).Map(func(x *match.Card) {
				ctx.Match.MoveCard(x, match.SHIELDZONE, card)
			})

		}
	})
}
