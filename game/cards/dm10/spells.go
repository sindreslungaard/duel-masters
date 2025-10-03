package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// Transmogrify ...
func Transmogrify(c *match.Card) {

	c.Name = "Transmogrify"
	c.Civ = civ.Water
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
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
			fmt.Sprintf("%s's effect: You may destroy a creature. If you do, its owner reveals cards from the top of this deck until he reveals a non-evolution creature. He puts that creature into the battle zone and puts the rest of those cards into his graveyard.", card.Name),
			1,
			1,
			true,
		).Map(func(x *match.Card) {
			ctx.Match.Destroy(x, card, match.DestroyedBySpell)

			for {
				topCards := x.Player.PeekDeck(1)

				if len(topCards) < 1 {
					return
				}

				topCard := topCards[0]

				if topCard.HasCondition(cnd.Creature) && !topCard.HasCondition(cnd.Evolution) {
					fx.ForcePutCreatureIntoBZ(ctx, topCard, match.DECK, card)
					return
				}

				ctx.Match.MoveCard(topCard, match.GRAVEYARD, card)
			}
		})
	}))

}
