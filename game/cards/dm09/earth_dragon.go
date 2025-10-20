package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// TerradragonAnristVhal ...
func TerradragonAnristVhal(c *match.Card) {

	c.Name = "Terradragon Anrist Vhal"
	c.Power = 0
	c.Civ = civ.Nature
	c.Family = []string{family.EarthDragon}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Nature}

	addPower := 0

	c.Use(fx.Creature, fx.When(fx.InTheBattlezone, func(card *match.Card, ctx *match.Context) {
		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
			addPower = 0
			c.Power = 0
			card.RemoveConditionBySource(card.ID)

			if card.Zone != match.BATTLEZONE {
				exit()
				return
			}

			addPower += len(fx.FindFilter(
				card.Player,
				match.BATTLEZONE,
				func(x *match.Card) bool {
					return x.ID != card.ID && x.Civ == civ.Nature
				})) * 2000

			if addPower == 0 {
				card.RemoveConditionBySource(card.ID)
				exit() //NOTE: We call exit() before calling Destroy below to avoid stack overflow error!
				ctx2.Match.Destroy(card, card, match.DestroyedByMiscAbility)
				return
			}

			c.Power += addPower

			if c.Power >= 6000 {
				card.AddUniqueSourceCondition(cnd.DoubleBreaker, true, card.ID)
			}
		})
	}))

}
