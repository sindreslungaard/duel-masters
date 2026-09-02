package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// GeneralDarkFiend ...
func GeneralDarkFiend(c *match.Card) {

	c.Name = "General Dark Fiend"
	c.Power = 6000
	c.Civs = []string{civ.Darkness}
	c.Family = []string{family.DarkLord}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Doublebreaker, fx.When(fx.Attacking, func(card *match.Card, ctx *match.Context) {

		// Backside selection keeps the shield's face hidden from its own
		// controller too, which is what "without looking" means here. The
		// choice is mandatory, so it is not cancellable, and MoveCard (rather
		// than BreakShields) never offers the shield trigger on the way to
		// the graveyard.
		fx.SelectBackside(
			card.Player,
			ctx.Match,
			card.Player,
			match.SHIELDZONE,
			fmt.Sprintf("%s: Choose one of your shields without looking and put it into your graveyard.", card.Name),
			1,
			1,
			false,
		).Map(func(shield *match.Card) {
			ctx.Match.MoveCard(shield, match.GRAVEYARD, card)
		})

	}))

}
