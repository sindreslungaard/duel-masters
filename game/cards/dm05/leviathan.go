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

	c.Use(fx.Creature, fx.Triplebreaker, func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.CardMoved); ok {

			if event.CardID == card.ID && event.To == match.BATTLEZONE {

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
			}
		}
	})

}

// KingMazelan ...
func KingMazelan(c *match.Card) {

	c.Name = "King Mazelan"
	c.Power = 7000
	c.Civ = civ.Water
	c.Family = family.Leviathan
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Doublebreaker, func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.CardMoved); ok {

			if event.CardID != card.ID || event.To != match.BATTLEZONE {
				return
			}

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

			ctx.Match.NewMultipartAction(card.Player, cards, 1, 1, "Choose 1 creature in the battlezone that will be sent to its owners hands", true)

			for {

				action := <-card.Player.Action

				if action.Cancel {
					break
				}

				if len(action.Cards) != 1 {
					ctx.Match.DefaultActionWarning(card.Player)
					continue
				}

				for _, vid := range action.Cards {

					ref, err := c.Player.MoveCard(vid, match.BATTLEZONE, match.HAND)

					if err != nil {

						ref, err := ctx.Match.Opponent(c.Player).MoveCard(vid, match.BATTLEZONE, match.HAND)

						if err == nil {
							ctx.Match.Chat("Server", fmt.Sprintf("%s was moved to %s's hand by %s", ref.Name, ctx.Match.PlayerRef(ref.Player).Socket.User.Username, card.Name))
						}

					} else {
						ctx.Match.Chat("Server", fmt.Sprintf("%s was moved to %s's hand by %s", ref.Name, ctx.Match.PlayerRef(ref.Player).Socket.User.Username, card.Name))
					}

				}

				break

			}

			ctx.Match.CloseAction(c.Player)

		}

	})

}
