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
		myManaCards := fx.Find(card.Player, match.MANAZONE)
		myHandCards := fx.FindFilter(card.Player, match.HAND, func(x *match.Card) bool { return x.ID != card.ID })
		myOppManaCards := fx.Find(ctx.Match.Opponent(card.Player), match.MANAZONE)
		myOppHandCards := fx.Find(ctx.Match.Opponent(card.Player), match.HAND)

		for _, x := range myManaCards {
			ctx.Match.MoveCard(x, match.HAND, card, true)
		}

		for _, x := range myHandCards {
			ctx.Match.MoveCard(x, match.MANAZONE, card, true)
			x.Tapped = true
		}

		for _, x := range myOppManaCards {
			ctx.Match.MoveCard(x, match.HAND, card, true)
		}

		for _, x := range myOppHandCards {
			ctx.Match.MoveCard(x, match.MANAZONE, card, true)
			x.Tapped = true
		}

		ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s's effect: Both players moved their mana cards to their hand, and at the same time, their hand cards to their mana zone.", card.Name))
	}))

}
