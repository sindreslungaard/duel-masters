package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// FortressShell ...
func FortressShell(c *match.Card) {

	c.Name = "Fortress Shell"
	c.Power = 5000
	c.Civ = civ.Nature
	c.Family = []string{family.ColonyBeetle}
	c.ManaCost = 9
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		if match.AmISummoned(card, ctx) {

			manaCards := match.Search(
				card.Player,
				ctx.Match,
				ctx.Match.Opponent(card.Player),
				match.MANAZONE,
				"Select up to 2 cards in your opponent's mana zone that will be sent to their graveyard",
				1,
				2,
				true,
			)

			for _, mana := range manaCards {
				ctx.Match.Opponent(card.Player).MoveCard(mana.ID, match.MANAZONE, match.GRAVEYARD)
				ctx.Match.Chat("Server", fmt.Sprintf("%s was sent from %s's manazone to their graveyard by %s", mana.Name, mana.Player.Username(), card.Name))
			}

		}

	})

}
