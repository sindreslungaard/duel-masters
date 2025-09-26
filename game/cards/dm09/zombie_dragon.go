package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// NecrodragonIzoristVhal ...
func NecrodragonIzoristVhal(c *match.Card) {

	c.Name = "Necrodragon Izorist Vhal"
	c.Power = 0
	c.Civ = civ.Darkness
	c.Family = []string{family.ZombieDragon}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Darkness}

	addPower := 0

	c.Use(fx.Creature, fx.When(fx.InTheBattlezone, func(card *match.Card, ctx *match.Context) {
		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
			c.Power = 0
			c.RemoveConditionBySource(card.ID)

			if card.Zone != match.BATTLEZONE {
				exit()
				return
			}

			addPower += len(fx.FindFilter(
				card.Player,
				match.GRAVEYARD,
				func(x *match.Card) bool {
					return len(x.Family) > 0 && x.Civ == civ.Darkness
				})) * 2000

			if addPower == 0 {
				exit()
				ctx2.Match.Destroy(card, card, match.DestroyedByMiscAbility)
				return
			}

			c.Power += addPower

			if c.Power >= 6000 {
				c.AddUniqueSourceCondition(cnd.DoubleBreaker, true, card.ID)
			}
		})
	}))

}
