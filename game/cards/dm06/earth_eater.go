package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

func HazardCrawler(c *match.Card) {

	c.Name = "Hazard Crawler"
	c.Power = 6000
	c.Civ = civ.Water
	c.Family = []string{family.EarthEater}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Blocker, fx.CantAttackCreatures, fx.CantAttackPlayers)
}

func MidnightCrawler(c *match.Card) {

	c.Name = "Midnight Crawler"
	c.Power = 6000
	c.Civ = civ.Water
	c.Family = []string{family.EarthEater}
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Doublebreaker, fx.When(fx.Summoned, fx.ReturnOpCardFromMZToHand))
}

func ThrashCrawler(c *match.Card) {

	c.Name = "Thrash Crawler"
	c.Power = 5000
	c.Civ = civ.Water
	c.Family = []string{family.EarthEater}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Blocker, fx.CantAttackCreatures, fx.CantAttackPlayers, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.MANAZONE,
			"Thrash Crawler: Choose a card from your mana zone that will be returned to your hand.",
			1,
			1,
			false,
		).Map(func(x *match.Card) {
			x.Player.MoveCard(x.ID, match.MANAZONE, match.HAND, card.ID)
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s got moved to %s hand from his mana zone by Thrash Crawler", x.Name, x.Player.Username()))
		})
	}))
}
