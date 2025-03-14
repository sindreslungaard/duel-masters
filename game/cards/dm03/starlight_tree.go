package dm03

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// SparkleFlower ...
func SparkleFlower(c *match.Card) {

	c.Name = "Sparkle Flower"
	c.Power = 3000
	c.Civ = civ.Light
	c.Family = []string{family.StarlightTree}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature,

		fx.When(fx.InTheBattlezone, func(card *match.Card, ctx *match.Context) {

			ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {

				if card.Zone != match.BATTLEZONE {
					card.RemoveConditionBySource(card.ID)

					exit()
					return
				}

				if match.ContainerHas(card.Player, match.MANAZONE, func(x *match.Card) bool { return x.Civ != civ.Light }) {
					card.RemoveConditionBySource(card.ID)
				} else {
					fx.ForceBlocker(card, ctx2, card.ID)
				}

			})

		}))

}
