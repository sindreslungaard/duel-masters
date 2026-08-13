package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// ChargeWhipper ...
func ChargeWhipper(c *match.Card) {

	c.Name = "Charge Whipper"
	c.Power = 2000
	c.Civs = []string{civ.Water}
	c.Family = []string{family.CyberVirus}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.SilentSkill(func(card *match.Card, ctx *match.Context) {

		// The two halves are one transaction: the printed "if you do" means a
		// shield only comes back when a card actually went in.
		added := false

		fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.HAND,
			fmt.Sprintf("%s's effect: You may add a card from your hand to your shields face down.", card.Name),
			1,
			1,
			true,
		).Map(func(x *match.Card) {
			moved, err := card.Player.MoveCard(x.ID, match.HAND, match.SHIELDZONE, card.ID)

			if err != nil || moved.Zone != match.SHIELDZONE {
				return
			}

			added = true
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s added a card from their hand to their shields", card.Player.Username()))
		})

		if !added {
			return
		}

		// The shield is moved, not broken, so its shield trigger is never
		// offered. That is exactly what the reminder text describes.
		fx.SelectBackside(
			card.Player,
			ctx.Match,
			card.Player,
			match.SHIELDZONE,
			fmt.Sprintf("%s's effect: Choose one of your shields and put it into your hand. You can't use its shield trigger.", card.Name),
			1,
			1,
			false,
		).Map(func(x *match.Card) {
			moved, err := card.Player.MoveCard(x.ID, match.SHIELDZONE, match.HAND, card.ID)

			if err != nil || moved.Zone != match.HAND {
				return
			}

			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s put one of their shields into their hand with %s", card.Player.Username(), card.Name))
		})
	}))

}
