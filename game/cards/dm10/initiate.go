package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// LemikVizierOfThought ...
func LemikVizierOfThought(c *match.Card) {

	c.Name = "Lemik, Vizier of Thought"
	c.Power = 3000
	c.Civ = civ.Light
	c.Family = []string{family.Initiate}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.When(fx.InTheBattlezone, func(card *match.Card, ctx *match.Context) {
		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
			if card.Zone != match.BATTLEZONE {
				fx.Find(
					card.Player,
					match.BATTLEZONE,
				).Map(func(x *match.Card) {
					x.RemoveConditionBySource(card.ID)
				})

				exit()
				return
			}

			fx.FindFilter(
				card.Player,
				match.BATTLEZONE,
				func(x *match.Card) bool { return x.Civ == civ.Water || x.Civ == civ.Nature },
			).Map(func(x *match.Card) {
				fx.ForceBlocker(x, ctx2, card.ID)
			})
		})
	}))

}
