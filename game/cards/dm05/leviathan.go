package dm05

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// KingTsunami ...
func KingTsunami(c *match.Card) {

	c.Name = "King Tsunami"
	c.Power = 12000
	c.Civ = civ.Water
	c.Family = []string{family.Leviathan}
	c.ManaCost = 12
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Triplebreaker, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		for _, creature := range append(fx.Find(card.Player, match.BATTLEZONE), fx.Find(ctx.Match.Opponent(card.Player), match.BATTLEZONE)...) {

			if creature.ID != card.ID {

				creature.Player.MoveCard(creature.ID, match.BATTLEZONE, match.HAND, card.ID)
				ctx.Match.ReportActionInChat(creature.Player, fmt.Sprintf("%s was returned to %s's hand by %s", creature.Name, creature.Player.Username(), card.Name))

			}
		}
	}))

}

// KingMazelan ...
func KingMazelan(c *match.Card) {

	c.Name = "King Mazelan"
	c.Power = 7000
	c.Civ = civ.Water
	c.Family = []string{family.Leviathan}
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Doublebreaker, fx.When(fx.Summoned, fx.MayReturnCreatureToOwnersHand))

}
