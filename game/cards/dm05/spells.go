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

// CyclonePanic ...
func CyclonePanic(c *match.Card) {

	c.Name = "Cyclone Panic"
	c.Civ = civ.Fire
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			// p1
			cards1 := fx.Find(card.Player, match.HAND)
			n1 := len(cards1)

			for _, c1 := range cards1 {
				card.Player.MoveCard(c1.ID, match.HAND, match.DECK)
			}

			card.Player.ShuffleDeck()
			card.Player.DrawCards(n1)

			// p2
			cards2 := fx.Find(ctx.Match.Opponent(card.Player), match.HAND)
			n2 := len(cards1)

			for _, c2 := range cards2 {
				ctx.Match.Opponent(card.Player).MoveCard(c2.ID, match.HAND, match.DECK)
			}

			ctx.Match.Opponent(card.Player).ShuffleDeck()
			ctx.Match.Opponent(card.Player).DrawCards(n2)

			ctx.Match.Chat("Server", "Cyclone Panic: Both players shuffled their deck and replaced the cards in their hand with new ones")

		}
	})
}

// GlorySnow ...
func GlorySnow(c *match.Card) {

	c.Name = "Glory Snow"
	c.Civ = civ.Light
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.ShieldTrigger, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			mana, _ := card.Player.Container(match.MANAZONE)
			opnnt_mana, _ := ctx.Match.Opponent(card.Player).Container(match.MANAZONE)

			if len(mana) < len(opnnt_mana) {
				cards := card.Player.PeekDeck(2)

				for _, toMove := range cards {

					card.Player.MoveCard(toMove.ID, match.DECK, match.MANAZONE)
					ctx.Match.Chat("Server", fmt.Sprintf("%s put %s into the manazone from the top of their deck", card.Player.Username(), toMove.Name))

				}
			}

		}

	})

}
