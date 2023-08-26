package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// HunterFish ...
func HunterFish(c *match.Card) {

	c.Name = "Hunter Fish"
	c.Power = 3000
	c.Civ = civ.Water
	c.Family = []string{family.Fish}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Blocker, fx.CantAttackPlayers, fx.CantAttackCreatures)

}

// Seamine ...
func Seamine(c *match.Card) {

	c.Name = "Seamine"
	c.Power = 4000
	c.Civ = civ.Water
	c.Family = []string{family.Fish}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Blocker)

}

// UnicornFish ...
func UnicornFish(c *match.Card) {

	c.Name = "Unicorn Fish"
	c.Power = 1000
	c.Civ = civ.Water
	c.Family = []string{family.Fish}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

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

			ctx.Match.NewMultipartAction(card.Player, cards, 1, 1, "Unicorn Fish: Choose 1 creature in the battlezone that will be sent to its owners hands", true)

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
							ctx.Match.Chat("Server", fmt.Sprintf("%s was moved to %s's hand", ref.Name, ctx.Match.PlayerRef(ref.Player).Socket.User.Username))
						}

					} else {
						ctx.Match.Chat("Server", fmt.Sprintf("%s was moved to %s's hand", ref.Name, ctx.Match.PlayerRef(ref.Player).Socket.User.Username))
					}

				}

				break

			}

			ctx.Match.CloseAction(c.Player)

		}

	})

}
