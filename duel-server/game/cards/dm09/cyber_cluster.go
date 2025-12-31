package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// TentacleCluster ...
func TentacleCluster(c *match.Card) {

	c.Name = "Tentacle Cluster"
	c.Power = 2000
	c.Civ = civ.Water
	c.Family = []string{family.CyberCluster}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.When(fx.WheneverThisAttacksPlayerAndIsntBlocked, fx.ReturnCreatureToOwnersHand))

}

// ScoutCluster ...
func ScoutCluster(c *match.Card) {

	c.Name = "Scout Cluster"
	c.Power = 4000
	c.Civ = civ.Water
	c.Family = []string{family.CyberCluster}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Blocker(), fx.When(fx.AnotherOwnCreatureSummoned, func(card *match.Card, ctx *match.Context) {
		card.Player.MoveCard(card.ID, match.BATTLEZONE, match.HAND, card.ID)
		ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was returned to the hand", card.Name))
	}))

}
