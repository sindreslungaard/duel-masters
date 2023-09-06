package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

func CharmiliaTheEnticer(c *match.Card) {

	c.Name = "Charmilia, the Enticer"
	c.Power = 3000
	c.Civ = civ.Nature
	c.Family = []string{family.SnowFaerie}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}
	c.TapAbility = true

	c.Use(fx.Creature, fx.When(fx.TapAbility, func(card *match.Card, ctx *match.Context) {

		cards := match.Filter(
			card.Player,
			ctx.Match,
			card.Player,
			match.DECK,
			"Select 1 creature from your deck that will be shown to your opponent and sent to your hand",
			0,
			1,
			false,
			func(x *match.Card) bool { return x.HasCondition(cnd.Creature) },
		)

		for _, c := range cards {
			card.Player.MoveCard(c.ID, match.DECK, match.HAND)
			ctx.Match.Chat("Server", fmt.Sprintf("%s was moved from %s's deck to their hand by Charmilia, the Enticer's tap abillity", c.Name, card.Player.Username()))
		}

		card.Player.ShuffleDeck()
		card.Tapped = true
	}))
}

func GarabonTheGlider(c *match.Card) {

	c.Name = "Garabon, the Glider"
	c.Power = 1000
	c.Civ = civ.Nature
	c.Family = []string{family.SnowFaerie}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.PowerAttacker2000)
}
