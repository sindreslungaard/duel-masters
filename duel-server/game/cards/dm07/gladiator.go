package dm07

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

func RondobilTheExplorer(c *match.Card) {

	c.Name = "Rondobil, the Explorer"
	c.Power = 5000
	c.Civ = civ.Light
	c.Family = []string{family.Gladiator}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Light}
	c.TapAbility = func(card *match.Card, ctx *match.Context) {
		fx.Select(card.Player, ctx.Match, card.Player, match.BATTLEZONE,
			"Select creature to add to shields", 1, 1, false,
		).Map(func(x *match.Card) {
			card.Player.MoveCard(x.ID, match.BATTLEZONE, match.SHIELDZONE, card.ID)
			ctx.Match.ReportActionInChat(
				card.Player,
				fmt.Sprintf("Rondobil, the Explorer effect: %s was added to shieldzone", c.Name),
			)
		})
	}

	c.Use(fx.Creature, fx.TapAbility)
}
