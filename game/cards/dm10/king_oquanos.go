package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// KingOquanos ...
func KingOquanos(c *match.Card) {

	c.Name = "King Oquanos"
	c.Power = 2000
	c.Civ = civ.Water
	c.Family = []string{family.Leviathan}
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.When(fx.InTheBattlezone, func(card *match.Card, ctx *match.Context) {
		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
			card.Power = 2000
			card.RemoveConditionBySource(card.ID)

			if card.Zone != match.BATTLEZONE {
				exit()
				return
			}

			card.Power += len(fx.FindFilter(
				ctx2.Match.Opponent(card.Player),
				match.MANAZONE,
				func(x *match.Card) bool {
					return x.Tapped
				})) * 2000

			if card.Power >= 6000 {
				card.AddUniqueSourceCondition(cnd.DoubleBreaker, true, card.ID)
			}
		})
	}))

}
