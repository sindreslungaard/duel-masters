package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// SirenConcerto ...
func SirenConcerto(c *match.Card) {

	c.Name = "Siren Concerto"
	c.Civ = civ.Water
	c.ManaCost = 1
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.MANAZONE,
			fmt.Sprintf("%s's effect: Put a card from your mana zone into your hand.", card.Name),
			1,
			1,
			false,
		).Map(func(x *match.Card) {
			_, err := card.Player.MoveCard(x.ID, match.MANAZONE, match.HAND, card.ID)
			if err == nil {
				ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s put %s into his hand from his mana zone.", card.Player.Username(), x.Name))
			}
		})

		fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.HAND,
			fmt.Sprintf("%s's effect: Put a card from your hand into your mana zone.", card.Name),
			1,
			1,
			false,
		).Map(func(x *match.Card) {
			_, err := card.Player.MoveCard(x.ID, match.HAND, match.MANAZONE, card.ID)
			if err == nil {
				ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s put %s into his mana zone from his hand.", card.Player.Username(), x.Name))
			}
		})
	}))

}
