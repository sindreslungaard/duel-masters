package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// LahPurificationEnforcer ...
func LahPurificationEnforcer(c *match.Card) {

	c.Name = "Lah, Purification Enforcer"
	c.Power = 5500
	c.Civ = civ.Light
	c.Family = []string{family.Berserker}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature)

}

// RaylaTruthEnforcer ...
func RaylaTruthEnforcer(c *match.Card) {

	c.Name = "Rayla, Truth Enforcer"
	c.Power = 3000
	c.Civ = civ.Light
	c.Family = []string{family.Berserker}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		cards := fx.SelectFilterFullList(
			card.Player,
			ctx.Match,
			card.Player,
			match.DECK,
			"Select 1 spell from your deck that will be shown to your opponent and sent to your hand",
			1,
			1,
			true,
			func(x *match.Card) bool { return x.HasCondition(cnd.Spell) },
			true,
		)

		for _, c := range cards {
			card.Player.MoveCard(c.ID, match.DECK, match.HAND, card.ID)
			ctx.Match.Chat("Server", fmt.Sprintf("%s was moved from %s's deck to their hand", c.Name, card.Player.Username()))
		}

		card.Player.ShuffleDeck()

	}))

}
