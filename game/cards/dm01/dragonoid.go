package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// DeadlyFighterBraidClaw ...
func DeadlyFighterBraidClaw(c *match.Card) {

	c.Name = "Deadly Fighter Braid Claw"
	c.Power = 1000
	c.Civ = civ.Fire
	c.Family = []string{family.Dragonoid}
	c.ManaCost = 1
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.ForceAttack)

}

// ExplosiveFighterUcarn ...
func ExplosiveFighterUcarn(c *match.Card) {

	c.Name = "Explosive Fighter Ucarn"
	c.Power = 9000
	c.Civ = civ.Fire
	c.Family = []string{family.Dragonoid}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.Doublebreaker, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.MANAZONE,
			fmt.Sprintf("%s: Select 2 cards from your manazone that will be sent to your graveyard", card.Name),
			2,
			2,
			false,
		).Map(func(x *match.Card) {
			x.Player.MoveCard(x.ID, match.MANAZONE, match.GRAVEYARD, card.ID)
			ctx.Match.ReportActionInChat(x.Player, fmt.Sprintf("%s was sent from %s's manazone to their graveyard", x.Name, x.Player.Username()))
		})
	}))

}

// FireSweeperBurningHellion ...
func FireSweeperBurningHellion(c *match.Card) {

	c.Name = "Fire Sweeper Burning Hellion"
	c.Power = 3000
	c.Civ = civ.Fire
	c.Family = []string{family.Dragonoid}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.PowerAttacker2000)

}

// OnslaughterTriceps ...
func OnslaughterTriceps(c *match.Card) {

	c.Name = "Onslaughter Triceps"
	c.Power = 5000
	c.Civ = civ.Fire
	c.Family = []string{family.Dragonoid}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.MANAZONE,
			fmt.Sprintf("%s: Select 1 card from your manazone that will be sent to your graveyard", card.Name),
			1,
			1,
			false,
		).Map(func(x *match.Card) {
			x.Player.MoveCard(x.ID, match.MANAZONE, match.GRAVEYARD, card.ID)
			ctx.Match.ReportActionInChat(x.Player, fmt.Sprintf("%s was sent from %s's manazone to their graveyard", x.Name, card.Player.Username()))
		})
	}))

}

// SuperExplosiveVolcanodon ...
func SuperExplosiveVolcanodon(c *match.Card) {

	c.Name = "Super Explosive Volcanodon"
	c.Power = 2000
	c.Civ = civ.Fire
	c.Family = []string{family.Dragonoid}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.PowerAttacker4000)

}
