package dm03

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
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
	c.Family = []string{family.Leviathan}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.When(fx.Attacking, func(card *match.Card, ctx *match.Context) {

		cards := make(map[string][]*match.Card)

		myCards := fx.FindFilter(
			card.Player,
			match.BATTLEZONE,
			func(x *match.Card) bool { return ctx.Match.GetPower(x, false) <= 2000 },
		)

		opponentCards := fx.FindFilter(
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			func(x *match.Card) bool { return ctx.Match.GetPower(x, false) <= 2000 },
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

				ref, err := c.Player.MoveCard(vid, match.BATTLEZONE, match.HAND, card.ID)

				if err != nil {

					ref, err := ctx.Match.Opponent(c.Player).MoveCard(vid, match.BATTLEZONE, match.HAND, card.ID)

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

// KingPonitas ...
func KingPonitas(c *match.Card) {

	c.Name = "King Ponitas"
	c.Power = 4000
	c.Civ = civ.Water
	c.Family = []string{family.Leviathan}
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.When(fx.Attacking, func(card *match.Card, ctx *match.Context) {

		ctx.ScheduleAfter(func() {
			waterCards := fx.SelectFilterFullList(card.Player, ctx.Match, card.Player, match.DECK, "Select 1 water card from your deck that will be shown to your opponent and sent to your hand", 1, 1, true, func(x *match.Card) bool { return x.Civ == civ.Water }, true)

			for _, waterCard := range waterCards {

				card.Player.MoveCard(waterCard.ID, match.DECK, match.HAND, card.ID)
				ctx.Match.Chat("Server", fmt.Sprintf("%s retrieved %s from the deck to their hand", card.Player.Username(), waterCard.Name))

			}

			card.Player.ShuffleDeck()
		})
	}))
}

// LegendaryBynor ...
func LegendaryBynor(c *match.Card) {

	c.Name = "Legendary Bynor"
	c.Power = 8000
	c.Civ = civ.Water
	c.Family = []string{family.Leviathan}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Evolution, fx.Doublebreaker, func(card *match.Card, ctx *match.Context) {

		if card.Zone != match.BATTLEZONE {
			return
		}

		if event, ok := ctx.Event.(*match.AttackCreature); ok {
			legendaryBynorSpecial(card, ctx, event.CardID)
		}

		if event, ok := ctx.Event.(*match.AttackPlayer); ok {
			legendaryBynorSpecial(card, ctx, event.CardID)
		}

	})

}

func legendaryBynorSpecial(card *match.Card, ctx *match.Context, cardID string) {

	p := ctx.Match.CurrentPlayer()

	creature, err := p.Player.GetCard(cardID, match.BATTLEZONE)

	if err != nil {
		return
	}

	if creature.Civ != civ.Water || creature.ID == card.ID {
		return
	}

	if ctx.Match.IsPlayerTurn(card.Player) {
		creature.AddCondition(cnd.CantBeBlocked, true, card.ID)
	}
}
