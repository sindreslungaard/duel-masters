package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// MagmadragonOgristVhal ...
func MagmadragonOgristVhal(c *match.Card) {

	c.Name = "Magmadragon Ogrist Vhal"
	c.Power = 0
	c.Civ = civ.Fire
	c.Family = []string{family.VolcanoDragon}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Fire}

	addPower := 0

	c.Use(fx.Creature, fx.When(fx.InTheBattlezone, func(card *match.Card, ctx *match.Context) {
		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
			c.Power = 0
			c.RemoveConditionBySource(card.ID)

			if card.Zone != match.BATTLEZONE {
				exit()
				return
			}

			addPower = len(fx.Find(card.Player, match.HAND)) * 3000

			if addPower == 0 {
				exit()
				ctx2.Match.Destroy(card, card, match.DestroyedByMiscAbility)
				return
			}

			c.Power += addPower

			if c.Power >= 15000 {
				c.AddUniqueSourceCondition(cnd.TripleBreaker, true, card.ID)
			} else if c.Power >= 6000 {
				c.AddUniqueSourceCondition(cnd.DoubleBreaker, true, card.ID)
			}
		})
	}))
}
