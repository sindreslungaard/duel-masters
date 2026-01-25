package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// CarnivalTotem ...
func CarnivalTotem(c *match.Card) {

	c.Name = "Carnival Totem"
	c.Power = 7000
	c.Civ = civ.Nature
	c.Family = []string{family.MysteryTotem}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.Doublebreaker, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		myManaCards, _ := card.Player.Container(match.MANAZONE)
		myHandCards, _ := card.Player.Container(match.HAND)

		if myManaCards != nil && myHandCards != nil {
			for _, manaCard := range myManaCards {
				ctx.Match.MoveCard(manaCard, match.HAND, card)
			}

			for _, handCard := range myHandCards {
				ctx.Match.MoveCard(handCard, match.MANAZONE, card) //@TODO after merging fix the call to MoveCard without chat report for all cards
				handCard.Tapped = true
			}
		}
	}))

}

// JigglyTotem ...
func JigglyTotem(c *match.Card) {

	c.Name = "Jiggly Totem"
	c.Power = 1000
	c.Civ = civ.Nature
	c.Family = []string{family.MysteryTotem}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.PowerModifier = func(m *match.Match, attacking bool) int {
		if attacking {
			return len(fx.FindFilter(
				c.Player,
				match.MANAZONE,
				func(x *match.Card) bool {
					return x.Tapped
				},
			)) * 1000
		}

		return 0
	}

	c.Use(fx.Creature)

}
