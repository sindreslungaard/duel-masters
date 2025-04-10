package dm05

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// SyriusFirmamentElemental ...
func SyriusFirmamentElemental(c *match.Card) {

	c.Name = "Syrius, Firmament Elemental"
	c.Power = 12000
	c.Civ = civ.Light
	c.Family = []string{family.AngelCommand}
	c.ManaCost = 11
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Blocker(), fx.Triplebreaker)

}

// SyforceAuroraElemental ...
func SyforceAuroraElemental(c *match.Card) {

	c.Name = "Syforce, Aurora Elemental"
	c.Power = 7000
	c.Civ = civ.Light
	c.Family = []string{family.AngelCommand}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Blocker(), fx.Doublebreaker, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		fx.SelectFilter(
			card.Player,
			ctx.Match,
			card.Player,
			match.MANAZONE,
			fmt.Sprintf("%s: You may select 1 spell from your mana zone that will be sent to your hand", card.Name),
			1,
			1,
			true,
			func(c *match.Card) bool { return c.HasCondition(cnd.Spell) },
			false,
		).Map(func(c *match.Card) {

			c.Player.MoveCard(c.ID, match.MANAZONE, match.HAND, card.ID)
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s retrieved %s from the mana zone to their hand", c.Player.Username(), c.Name))

		})

	}))

}
