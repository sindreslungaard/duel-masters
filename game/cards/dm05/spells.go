package dm05

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// EnchantedSoil ...
func EnchantedSoil(c *match.Card) {

	c.Name = "Enchanted Soil"
	c.Civ = civ.Nature
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Spell, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			fx.SelectFilter(card.Player, ctx.Match, card.Player, match.GRAVEYARD, "Enchanted Soil: Select 2 creatures from your graveyard and put it in your manazone", 0, 2, true, func(x *match.Card) bool {
				return x.HasCondition(cnd.Creature)
			}).Map(func(x *match.Card) {
				card.Player.MoveCard(card.ID, match.GRAVEYARD, match.MANAZONE)
				ctx.Match.Chat("Server", fmt.Sprintf("%s was moved to %s's manazone", card.Name, card.Player.Username()))
			})

		}
	})
}

// SchemingHands ...
func SchemingHands(c *match.Card) {

	c.Name = "Scheming Hands"
	c.Civ = civ.Darkness
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Spell, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			fx.Select(card.Player, ctx.Match, ctx.Match.Opponent(card.Player), match.HAND, "Scheming Hands: Discard a card from your opponent's hand", 0, 1, false).Map(func(x *match.Card) {
				ctx.Match.Opponent(card.Player).MoveCard(card.ID, match.HAND, match.GRAVEYARD)
				ctx.Match.Chat("Server", fmt.Sprintf("%s was moved to %s's graveyard by Scheming Hands", x.Name, card.Player.Username()))
			})

		}
	})
}
