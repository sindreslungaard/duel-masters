package dm11

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// SquawkingLunatron ...
func SquawkingLunatron(c *match.Card) {

	c.Name = "Squawking Lunatron"
	c.Power = 4000
	c.Civs = []string{civ.Water}
	c.Family = []string{family.CyberMoon}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.SilentSkill(func(card *match.Card, ctx *match.Context) {
		fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.MANAZONE,
			fmt.Sprintf("%s's effect: Return up to 3 cards from your mana zone to your hand.", card.Name),
			1,
			3,
			true,
		).Map(func(x *match.Card) {
			moved, err := card.Player.MoveCard(x.ID, match.MANAZONE, match.HAND, card.ID)

			if err != nil || moved.Zone != match.HAND {
				return
			}

			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was returned from %s's mana zone to their hand by %s", x.Name, card.Player.Username(), card.Name))
		})
	}))

}
