package dm03

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// BonePiercer ...
func BonePiercer(c *match.Card) {

	c.Name = "Bone Piercer"
	c.Power = 1000
	c.Civ = civ.Darkness
	c.Family = []string{family.BrainJacker}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.When(fx.Destroyed, func(card *match.Card, ctx *match.Context) {
		fx.SelectFilter(
			card.Player,
			ctx.Match,
			card.Player,
			match.MANAZONE,
			fmt.Sprintf("%s: You may return 1 creature from your mana zone to your hand", card.Name),
			1,
			1,
			true,
			func(c *match.Card) bool { return c.HasCondition(cnd.Creature) },
			false,
		).Map(func(x *match.Card) {
			x.Player.MoveCard(x.ID, match.MANAZONE, match.HAND, card.ID)
			ctx.Match.ReportActionInChat(x.Player, fmt.Sprintf("%s was moved to %s's hand from their mana zone by %s", x.Name, x.Player.Username(), card.Name))
		})
	}))

}
