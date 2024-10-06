package dm07

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

func TangleFistTheWeaver(c *match.Card) {

	c.Name = "Tangle Fist, the Weaver"
	c.Power = 2000
	c.Civ = civ.Nature
	c.Family = []string{family.BeastFolk}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}
	c.TapAbility = func(card *match.Card, ctx *match.Context) {
		fx.Select(card.Player, ctx.Match, card.Player, match.HAND,
			"Select up to 3 cards to put into the manazone", 0, 3, false,
		).Map(func(x *match.Card) {
			card.Player.MoveCard(x.ID, match.HAND, match.MANAZONE, card.ID)
			ctx.Match.ReportActionInChat(
				card.Player,
				fmt.Sprintf("%s effect: %s moved to manazone from hand", card.Name, x.Name),
			)
		})
	}

	c.Use(fx.Creature, fx.TapAbility)
}
