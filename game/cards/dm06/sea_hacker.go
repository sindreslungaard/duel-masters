package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

func Aeropica(c *match.Card) {

	c.Name = "Aeropica"
	c.Power = 4000
	c.Civ = civ.Water
	c.Family = []string{family.SeaHacker}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Water}
	c.TapAbility = true

	c.Use(fx.Creature, fx.When(fx.TapAbility, func(card *match.Card, ctx *match.Context) {

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

		card.Tapped = true

	}))
}

func Zepimeteus(c *match.Card) {

	c.Name = "Zepimeteus"
	c.Power = 2000
	c.Civ = civ.Water
	c.Family = []string{family.SeaHacker}
	c.ManaCost = 1
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Blocker, fx.CantAttackCreatures, fx.CantAttackPlayers)
}

func PromephiusQ(c *match.Card) {

	c.Name = "Promephius Q"
	c.Power = 2000
	c.Civ = civ.Water
	c.Family = []string{family.SeaHacker, family.Survivor}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature)
}
