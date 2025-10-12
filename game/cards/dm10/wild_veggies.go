package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// ScowlingTomato ...
func ScowlingTomato(c *match.Card) {

	c.Name = "Scowling Tomato"
	c.Power = 2000
	c.Civ = civ.Nature
	c.Family = []string{family.WildVeggies}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature)

}

// ShamanBrocolli ...
func ShamanBrocolli(c *match.Card) {

	c.Name = "Shaman Brocolli"
	c.Power = 1000
	c.Civ = civ.Nature
	c.Family = []string{family.WildVeggies}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.When(fx.WouldBeDestroyed, func(card *match.Card, ctx *match.Context) {
		_, err := card.Player.MoveCard(card.ID, match.BATTLEZONE, match.MANAZONE, card.ID)

		if err == nil {
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was moved to mana zone instead of being destroyed.", card.Name))
		}
	}))

}
