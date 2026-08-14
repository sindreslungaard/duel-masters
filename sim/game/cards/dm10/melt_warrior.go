package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// CoreCrashLizard ...
func CoreCrashLizard(c *match.Card) {

	c.Name = "Core-Crash Lizard"
	c.Power = 6000
	c.Civs = []string{civ.Fire}
	c.Family = []string{family.MeltWarrior}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Fire}

	// The shield goes straight to the graveyard, so it is never broken and its
	// shield trigger is not offered.
	c.Use(fx.Creature, fx.When(fx.Summoned, fx.PutOpShieldIntoGraveyard))

}

// BurnwispLizard ...
func BurnwispLizard(c *match.Card) {

	c.Name = "Burnwisp Lizard"
	c.Power = 4000
	c.Civs = []string{civ.Fire}
	c.Family = []string{family.MeltWarrior}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.When(fx.InTheBattlezone, func(card *match.Card, ctx *match.Context) {
		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {

			// The grant is re-evaluated on every event rather than cached,
			// because the set of silent skill creatures changes as they are
			// summoned and destroyed, and their conditions are rebuilt each
			// untap step.
			if card.Zone != match.BATTLEZONE {
				fx.Find(card.Player, match.BATTLEZONE).Map(func(x *match.Card) {
					x.RemoveSpecificConditionBySource(cnd.SpeedAttacker, card.ID)
				})

				exit()
				return
			}

			fx.FindFilter(
				card.Player,
				match.BATTLEZONE,
				func(x *match.Card) bool { return x.HasCondition(cnd.SilentSkill) },
			).Map(func(x *match.Card) {
				x.AddUniqueSourceCondition(cnd.SpeedAttacker, true, card.ID)
			})
		})
	}))

}
