package dm05

import (
	"duel-masters/game/civ"
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

			cards := match.Search(card.Player, ctx.Match, card.Player, match.GRAVEYARD, "Select 2 cards from your graveyard and put it in your manazone", 0, 2, true)

			for _, card := range cards {

				card.Player.MoveCard(card.ID, match.GRAVEYARD, match.MANAZONE)
				ctx.Match.Chat("Server", fmt.Sprintf("%s was moved to %s's manazone", card.Name, card.Player.Username()))

			}

		}
	})
}
