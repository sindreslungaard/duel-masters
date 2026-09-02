package promo

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// dynoMantisTheMightspinnerBreakThreshold is the power at or above which this
// creature's other creatures break one more shield.
const dynoMantisTheMightspinnerBreakThreshold = 5000

// DynoMantisTheMightspinner ...
func DynoMantisTheMightspinner(c *match.Card) {

	c.Name = "Dyno Mantis, the Mightspinner"
	c.Power = 7000
	c.Civs = []string{civ.Nature}
	c.Family = []string{family.GiantInsect}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.Evolution, fx.Doublebreaker, fx.When(fx.InTheBattlezone, func(card *match.Card, ctx *match.Context) {
		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
			// GetPower dispatches handlers synchronously, and this effect
			// calls it below, so there is nothing to do while a power
			// calculation is already in flight.
			if _, ok := ctx2.Event.(*match.GetPowerEvent); ok {
				return
			}

			live := card.Zone == match.BATTLEZONE

			// A creature currently selecting how many shields it breaks is
			// evaluated as attacking, so a "power attacker" bonus - which
			// only applies while attacking - still counts here: it is still
			// in effect at this point, before the battle that would end it.
			attacker, _ := ctx2.Event.(*match.SelectShields)

			// Every owned zone is swept, not just the battle zone, so a
			// creature that left play does not carry a stale bonus back into
			// it before the end of the turn clears conditions.
			fx.FindMultiple(card.Player, []string{
				match.BATTLEZONE,
				match.HAND,
				match.GRAVEYARD,
				match.MANAZONE,
				match.SHIELDZONE,
				match.HIDDENZONE,
			}).Map(func(x *match.Card) {
				isAttacking := attacker != nil && attacker.Attacker.ID == x.ID

				if live && x.ID != card.ID && x.Zone == match.BATTLEZONE && ctx2.Match.GetPower(x, isAttacking) >= dynoMantisTheMightspinnerBreakThreshold {
					x.AddUniqueSourceCondition(cnd.ShieldBreakModifier, 1, card.ID)
					return
				}

				x.RemoveSpecificConditionBySource(cnd.ShieldBreakModifier, card.ID)
			})

			if !live {
				exit()
			}
		})
	}))

}
