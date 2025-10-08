package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// Upheaval ...
func Upheaval(c *match.Card) {

	c.Name = "Upheaval"
	c.Civ = civ.Darkness
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		myManaCards, err1 := card.Player.Container(match.MANAZONE)
		myHandCards, err2 := card.Player.Container(match.HAND)
		myOppManaCards, err3 := ctx.Match.Opponent(card.Player).Container(match.MANAZONE)
		myOppHandCards, err4 := ctx.Match.Opponent(card.Player).Container(match.HAND)

		if err1 == nil && err2 == nil && err3 == nil && err4 == nil {
			for _, x := range myManaCards {
				ctx.Match.MoveCard(x, match.HAND, card, true)
			}

			for _, x := range myHandCards {
				ctx.Match.MoveCard(x, match.MANAZONE, card, true)
			}

			for _, x := range myOppManaCards {
				ctx.Match.MoveCard(x, match.HAND, card, true)
			}

			for _, x := range myOppHandCards {
				ctx.Match.MoveCard(x, match.MANAZONE, card, true)
			}

			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s's effect: Both players moved their mana cards to their hand, and at the same time, their hand cards to their mana zone.", card.Name))
		}
	}))

}
