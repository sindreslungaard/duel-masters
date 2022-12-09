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

	c.Use(fx.Creature, fx.ShieldTrigger, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

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
