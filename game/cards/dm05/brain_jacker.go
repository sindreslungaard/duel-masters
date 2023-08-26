package dm05

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// SkullsweeperQ ...
func SkullsweeperQ(c *match.Card) {

	c.Name = "Skullsweeper Q"
	c.Power = 1000
	c.Civ = civ.Darkness
	c.Family = []string{family.BrainJacker}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Survivor, func(card *match.Card, ctx *match.Context) {

		if !ctx.Match.IsPlayerTurn(card.Player) || card.Zone != match.BATTLEZONE {
			return
		}

		event, ok := ctx.Event.(*match.AttackConfirmed)

		if !ok {
			return
		}

		creature, err := card.Player.GetCard(event.CardID, match.BATTLEZONE)

		if err != nil {
			return
		}

		if creature.HasCondition(cnd.Survivor) {
			fx.Select(
				ctx.Match.Opponent(card.Player),
				ctx.Match,
				ctx.Match.Opponent(card.Player),
				match.HAND,
				"Skullsweeper ability: Select 1 card from your hand that will be sent to your graveyard",
				1,
				1,
				false,
			).Map(func(x *match.Card) {
				x.Player.MoveCard(x.ID, match.HAND, match.GRAVEYARD)
				ctx.Match.Chat("Server", fmt.Sprintf("%s was moved from %s's hand to his graveyard due to %s's survivor ability", x.Name, x.Player.Username(), creature.Name))
			})
		}

	})

}

// JewelSpider ...
func JewelSpider(c *match.Card) {

	c.Name = "Jewel Spider"
	c.Power = 1000
	c.Civ = civ.Darkness
	c.Family = []string{family.BrainJacker}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.When(fx.Destroyed, func(card *match.Card, ctx *match.Context) {

		fx.SelectBackside(
			card.Player,
			ctx.Match,
			card.Player,
			match.SHIELDZONE,
			fmt.Sprintf("%s: You may put a shield to your hand.", card.Name),
			1,
			1,
			true,
		).Map(func(c *match.Card) {
			c.Player.MoveCard(c.ID, match.SHIELDZONE, match.HAND)
			ctx.Match.Chat("Server", fmt.Sprintf("A shield was moved to %s's hand by %s", card.Player.Username(), card.Name))
		})

	}))

}
