package dm03

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// KingNeptas ...
func KingNeptas(c *match.Card) {

	c.Name = "King Neptas"
	c.Power = 5000
	c.Civ = civ.Water
	c.Family = family.Leviathan
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.When(fx.Attacking, func(card *match.Card, ctx *match.Context) {

		cards := make(map[string][]*match.Card)

		myCards := fx.FindFilter(
			card.Player,
			match.BATTLEZONE,
			func(x *match.Card) bool { return x.Power <= 2000 },
		)

		opponentCards := fx.FindFilter(
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			func(x *match.Card) bool { return x.Power <= 2000 },
		)

		if len(myCards) < 1 && len(opponentCards) < 1 {
			return
		}

		cards["Your creatures"] = myCards
		cards["Opponent's creatures"] = opponentCards

		ctx.Match.NewMultipartAction(card.Player, cards, 0, 1, "Choose up to 1 creatures in the battle zone and return it to its owner hand", true)

		for {

			action := <-card.Player.Action

			if action.Cancel {
				break
			}

			if len(action.Cards) < 1 || len(action.Cards) > 1 {
				break
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

	}), fx.Creature)

}
