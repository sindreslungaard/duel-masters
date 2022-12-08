package dm05

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// AquaSurfer ...
func AquaSurfer(c *match.Card) {

	c.Name = "Aqua Surfer"
	c.Power = 2000
	c.Civ = civ.Water
	c.Family = family.LiquidPeople
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.ShieldTrigger, func(card *match.Card, ctx *match.Context) {

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
