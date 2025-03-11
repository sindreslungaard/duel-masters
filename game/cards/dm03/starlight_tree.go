package dm03

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
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

	// This card only checks for light cards at the untap step, when it should check every time your BZ get modified.
	// TODO: minor bug - fix required
	// TODO: see if it can be fixed with persistent effect
	c.Use(fx.Creature,

		// if _, ok := ctx.Event.(*match.UntapStep); ok {

		// 	if match.ContainerHas(card.Player, match.MANAZONE, func(x *match.Card) bool { return x.Civ != civ.Light }) {
		// 		card.RemoveCondition(cnd.Blocker)
		// 	} else {
		// 		card.AddCondition(cnd.Blocker, true, card.ID)
		// 	}
		// }

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
					card.AddUniqueSourceCondition(cnd.Blocker, true, card.ID)
				}

			})

		}))

}
