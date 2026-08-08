package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
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
		fx.SwapHandAndMana(card, card.Player)
		ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s's effect: %s swapped the cards in their hand and mana zone.", card.Name, card.Player.Username()))
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
