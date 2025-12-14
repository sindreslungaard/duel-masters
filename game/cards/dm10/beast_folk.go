package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// AdventureBoar ...
func AdventureBoar(c *match.Card) {

	c.Name = "Adventure Boar"
	c.Civ = civ.Nature
	c.Power = 1000
	c.Family = []string{family.BeastFolk}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.PowerAttacker2000)

}

// EarthRipperTalonOfRage ...
func EarthRipperTalonOfRage(c *match.Card) {

	c.Name = "Earth Ripper, Talon of Rage"
	c.Power = 6000
	c.Civ = civ.Nature
	c.Family = []string{family.BeastFolk}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.Evolution, fx.Doublebreaker, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		fx.FindFilter(
			card.Player,
			match.MANAZONE,
			func(x *match.Card) bool {
				return x.Tapped
			},
		).Map(func(x *match.Card) {
			_, err := x.Player.MoveCard(x.ID, match.MANAZONE, match.HAND, card.ID)

			if err == nil {
				ctx.Match.ReportActionInChat(x.Player, fmt.Sprintf("%s was moved from %s's mana zone to his hand due to %s's effect.", x.Name, card.Player.Username(), card.Name))
			}
		})
	}))

}
