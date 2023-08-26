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

		ctx.Match.Wait(ctx.Match.Opponent(c.Player), "Waiting for your opponent to make an action")

		fx.SelectFilter(
			card.Player,
			ctx.Match,
			card.Player,
			match.MANAZONE,
			"Bone Piercer: Select 1 creature from your manazone that will be sent to your hand",
			0,
			1,
			true,
			func(c *match.Card) bool { return c.HasCondition(cnd.Creature) },
		).Map(func(x *match.Card) {
			c.Player.MoveCard(x.ID, match.MANAZONE, match.HAND)
			ctx.Match.Chat("Server", fmt.Sprintf("%s was moved to %s's hand from their mana zone by %s", x.Name, x.Player.Username(), c.Name))
		})

		ctx.Match.EndWait(ctx.Match.Opponent(c.Player))

	}))

}
