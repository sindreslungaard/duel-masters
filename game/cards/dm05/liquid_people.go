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
	c.Family = []string{family.LiquidPeople}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.ShieldTrigger, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		cards := make(map[string][]*match.Card)

		cards["Your creatures"] = fx.Find(card.Player, match.BATTLEZONE)
		cards["Opponent's creatures"] = fx.Find(ctx.Match.Opponent(card.Player), match.BATTLEZONE)

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
