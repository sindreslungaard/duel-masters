package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// HurricaneCrawler ...
func HurricaneCrawler(c *match.Card) {

	c.Name = "Hurricane Crawler"
	c.Power = 4000
	c.Civ = civ.Water
	c.Family = []string{family.EarthEater}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		// Snapshot the hand before moving any of its cards so the loop does not
		// range over the live container while it is being mutated.
		hand := fx.FindFilter(card.Player, match.HAND, func(x *match.Card) bool {
			return x.ID != card.ID
		})

		// "Then put that many cards" refers to the cards that actually reached the
		// mana zone, so a move prevented by a replacement effect must not be counted.
		moved := 0
		for _, handCard := range hand {
			manaCard, err := card.Player.MoveCard(handCard.ID, match.HAND, match.MANAZONE, card.ID)
			if err == nil && manaCard.Zone == match.MANAZONE {
				moved++
			}
		}

		if moved < 1 {
			return
		}

		ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s effect: %s put %d card(s) from their hand into their mana zone", card.Name, card.Player.Username(), moved))

		// The mana zone is not guaranteed to still hold the cards that were just
		// added: every move above dispatched events another card can react to.
		// Select clamps min/max down to the cards actually present and returns an
		// empty collection for an empty zone, so a shrunken mana zone yields as
		// many cards as remain instead of an impossible prompt.
		fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.MANAZONE,
			fmt.Sprintf("%s effect: Select %d card(s) from your mana zone that will be put into your hand", card.Name, moved),
			moved,
			moved,
			false,
		).Map(func(x *match.Card) {
			card.Player.MoveCard(x.ID, match.MANAZONE, match.HAND, card.ID)
		})

		ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s effect: %s put %d card(s) from their mana zone into their hand", card.Name, card.Player.Username(), moved))

	}))

}
