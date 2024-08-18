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
		).Map(func(c *match.Card) {
			c.Player.MoveCard(c.ID, match.MANAZONE, match.HAND, card.ID)
			ctx.Match.Chat("Server", fmt.Sprintf("%s was returned to %s's hand from manazone", c.Name, c.Player.Username()))
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

		fx.SelectFilterFullList(
			card.Player,
			ctx.Match,
			card.Player,
			match.DECK,
			fmt.Sprintf("%s: Select 1 Giant Insect from your deck that will be shown to your opponent and sent to your hand", card.Name),
			1,
			1,
			true,
			func(c *match.Card) bool { return c.HasFamily(family.GiantInsect) },
			true,
		).Map(func(c *match.Card) {
			card.Player.MoveCard(c.ID, match.DECK, match.HAND, card.ID)
			ctx.Match.Chat("Server", fmt.Sprintf("%s was moved from %s's deck to their hand", c.Name, card.Player.Username()))
		})

		card.Player.ShuffleDeck()

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

		fx.SelectFilter(
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

			ctx.Match.Chat("Server", fmt.Sprintf("%s was moved to %s's battlezone from their manazone", c.Name, c.Player.Username()))

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

		fx.SelectFilter(
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

			ctx.Match.Chat("Server", fmt.Sprintf("%s was moved to %s's battlezone from their manazone", c.Name, c.Player.Username()))

		})

	}))

}
