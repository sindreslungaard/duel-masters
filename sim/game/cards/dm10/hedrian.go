package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// HourglassMutant ...
func HourglassMutant(c *match.Card) {

	c.Name = "Hourglass Mutant"
	c.Power = 2000
	c.Civs = []string{civ.Darkness}
	c.Family = []string{family.Hedrian}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.When(fx.InTheBattlezone, func(card *match.Card, ctx *match.Context) {
		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
			if card.Zone != match.BATTLEZONE {
				// Remove only this source's contribution, and do so in every zone
				// because a creature that was granted slayer may already have left
				// the battle zone before the grant ended.
				fx.FindMultiple(
					card.Player,
					[]string{
						match.BATTLEZONE,
						match.HIDDENZONE,
						match.HAND,
						match.SHIELDZONE,
						match.MANAZONE,
						match.GRAVEYARD,
						match.DECK,
					},
				).Map(func(x *match.Card) {
					x.RemoveSpecificConditionBySource(cnd.Slayer, card.ID)
				})

				exit()
				return
			}

			fx.FindFilter(
				card.Player,
				match.BATTLEZONE,
				func(x *match.Card) bool {
					return x.HasCondition(cnd.Creature) && (x.HasCiv(civ.Water) || x.HasCiv(civ.Fire))
				},
			).Map(func(x *match.Card) {
				x.AddUniqueSourceCondition(cnd.Slayer, true, card.ID)
			})
		})
	}))

}
