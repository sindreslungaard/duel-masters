package dm11

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// MelodicHunter ...
func MelodicHunter(c *match.Card) {

	c.Name = "Melodic Hunter"
	c.Power = 3000
	c.Civs = []string{civ.Water}
	c.Family = []string{family.Merfolk}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Blocker())

}

// TimeScout ...
func TimeScout(c *match.Card) {

	c.Name = "Time Scout"
	c.Power = 1000
	c.Civs = []string{civ.Water}
	c.Family = []string{family.Merfolk}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		top := ctx.Match.Opponent(card.Player).PeekDeck(1)

		if len(top) < 1 {
			return
		}

		// Looking only: the card stays on top of the deck, and the opponent is
		// shown nothing.
		ctx.Match.ShowCards(
			card.Player,
			fmt.Sprintf("%s's effect: the top card of your opponent's deck", card.Name),
			[]string{top[0].ImageID},
		)
	}))

}
