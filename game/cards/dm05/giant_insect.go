package dm05

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// BloodwingMantis ...
func BloodwingMantis(c *match.Card) {

	c.Name = "Bloodwing Mantis"
	c.Power = 6000
	c.Civ = civ.Nature
	c.Family = []string{family.GiantInsect}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.Doublebreaker, fx.When(fx.AttackConfirmed, func(card *match.Card, ctx *match.Context) {

		fx.SelectFilter(
			card.Player,
			ctx.Match,
			card.Player,
			match.MANAZONE,
			"Select 2 creatures from your mana zone to be returned to your hand",
			2,
			2,
			false,
			func(c *match.Card) bool { return c.HasCondition(cnd.Creature) },
			false,
		).Map(func(c *match.Card) {
			c.Player.MoveCard(c.ID, match.MANAZONE, match.HAND, card.ID)
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was sent to %s's hand from their manazone by %s's effect", c.Name, c.Player.Username(), card.Name))
		})

	}))

}

// ScissorScarab ...
func ScissorScarab(c *match.Card) {

	c.Name = "Scissor Scarab"
	c.Power = 5000
	c.Civ = civ.Nature
	c.Family = []string{family.GiantInsect}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		fx.SearchDeckTakeCards(
			card,
			ctx,
			1,
			func(x *match.Card) bool { return x.HasFamily(family.GiantInsect) },
			"Giant Insect",
		)

	}))

}

// AmbushScorpion ...
func AmbushScorpion(c *match.Card) {

	c.Name = "Ambush Scorpion"
	c.Power = 3000
	c.Civ = civ.Nature
	c.Family = []string{family.GiantInsect}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.PowerAttacker3000, fx.When(fx.Destroyed, func(card *match.Card, ctx *match.Context) {

		fx.SelectFilterSelectablesOnly(
			card.Player,
			ctx.Match,
			card.Player,
			match.MANAZONE,
			fmt.Sprintf("%s: choose an %s from your manazone and send it to battlezone", card.Name, card.Name),
			1,
			1,
			true,
			func(c *match.Card) bool { return c.Name == card.Name },
		).Map(func(c *match.Card) {

			c.Player.MoveCard(c.ID, match.MANAZONE, match.BATTLEZONE, card.ID)

			if ctx.Match.IsPlayerTurn(card.Player) {
				c.AddCondition(cnd.SummoningSickness, nil, nil)
			}

			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was moved to %s's battlezone from their manazone", c.Name, c.Player.Username()))

		})

	}))

}

// ObsidianScarab ...
func ObsidianScarab(c *match.Card) {

	c.Name = "Obsidian Scarab"
	c.Power = 5000
	c.Civ = civ.Nature
	c.Family = []string{family.GiantInsect}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.Doublebreaker, fx.PowerAttacker3000, fx.When(fx.Destroyed, func(card *match.Card, ctx *match.Context) {

		fx.SelectFilterSelectablesOnly(
			card.Player,
			ctx.Match,
			card.Player,
			match.MANAZONE,
			fmt.Sprintf("%s: choose an %s from your manazone and send it to battlezone", card.Name, card.Name),
			1,
			1,
			true,
			func(c *match.Card) bool { return c.Name == card.Name },
		).Map(func(c *match.Card) {

			c.Player.MoveCard(c.ID, match.MANAZONE, match.BATTLEZONE, card.ID)

			if ctx.Match.IsPlayerTurn(card.Player) {
				c.AddCondition(cnd.SummoningSickness, nil, nil)
			}

			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was moved to %s's battlezone from their manazone", c.Name, c.Player.Username()))

		})

	}))

}
