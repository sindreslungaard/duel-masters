package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

func NeonCluster(c *match.Card) {

	c.Name = "Neon Cluster"
	c.Power = 4000
	c.Civ = civ.Water
	c.Family = []string{family.CyberCluster}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Water}
	c.TapAbility = true

	c.Use(fx.Creature, fx.When(fx.TapAbility, func(card *match.Card, ctx *match.Context) {
		ctx.Match.Chat("Server", fmt.Sprintf("%s activated %s's tap ability to draw 2 cards", card.Player.Username(), card.Name))
		card.Player.DrawCards(2)
		card.Tapped = true
	}))

}
