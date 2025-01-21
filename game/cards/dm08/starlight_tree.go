package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// SolarGrass ...
func SolarGrass(c *match.Card) {

	c.Name = "Solar Grass"
	c.Power = 3000
	c.Civ = civ.Light
	c.Family = []string{family.StarlightTree}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	turboRush := false

	c.Use(
		fx.Creature,
		fx.When(fx.TurboRushCondition, func(card *match.Card, ctx *match.Context) { turboRush = true }),
		fx.When(fx.EndOfMyTurn, func(card *match.Card, ctx *match.Context) { turboRush = false }),
		fx.When(func(c *match.Card, ctx *match.Context) bool { return turboRush },
			fx.When(fx.WheneverThisAttacksAndIsntBlocked, func(card *match.Card, ctx *match.Context) {
				fx.FindFilter(
					card.Player,
					match.BATTLEZONE,
					func(creature *match.Card) bool { return creature.Name != card.Name },
				).Map(func(creature *match.Card) {
					creature.Tapped = false
					ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s untapped due to %s", creature.Name, card.Name))
				})
			})),
	)

}
