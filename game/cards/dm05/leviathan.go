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
	c.Family = family.Leviathan
	c.ManaCost = 12
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Triplebreaker, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		myCreatures, err := card.Player.Container(match.BATTLEZONE)
		if err != nil {
			return
		}

		opponentsCreatures, err := ctx.Match.Opponent(card.Player).Container(match.BATTLEZONE)
		if err != nil {
			return
		}

		for _, creature := range append(myCreatures, opponentsCreatures...) {

			if creature.ID != card.ID {

				creature.Player.MoveCard(creature.ID, match.BATTLEZONE, match.HAND)
				ctx.Match.Chat("Server", fmt.Sprintf("%s was returned to %s's hand by %s", creature.Name, creature.Player.Username(), card.Name))

			}
		}
	}))

}

// KingMazelan ...
func KingMazelan(c *match.Card) {

	c.Name = "King Mazelan"
	c.Power = 7000
	c.Civ = civ.Water
	c.Family = family.Leviathan
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Doublebreaker, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		cards := make(map[string][]*match.Card)

		myCards, err := card.Player.Container(match.BATTLEZONE)
		if err != nil {
			return
		}

		opponentCards, err := ctx.Match.Opponent(card.Player).Container(match.BATTLEZONE)
		if err != nil {
			return
		}

		if len(myCards) < 1 && len(opponentCards) < 1 {
			return
		}

		cards["Your creatures"] = myCards
		cards["Opponent's creatures"] = opponentCards

		fx.SelectMultipart(
			card.Player,
			ctx.Match,
			cards,
			fmt.Sprintf("%s: Choose 1 creature in the battlezone that will be sent to its owner's hand", card.Name),
			1,
			1,
			true).Map(func(creature *match.Card) {
			creature.Player.MoveCard(creature.ID, match.BATTLEZONE, match.HAND)
			ctx.Match.Chat("Server", fmt.Sprintf("%s was returned to %s's hand by %s", creature.Name, creature.Player.Username(), card.Name))
		})
	}))

}
