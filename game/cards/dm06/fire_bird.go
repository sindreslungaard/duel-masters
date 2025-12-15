package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func CoccoLupia(c *match.Card) {

	c.Name = "Cocco Lupia"
	c.Power = 1000
	c.Civ = civ.Fire
	c.Family = []string{family.FireBird}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {
		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
			if card.Zone != match.BATTLEZONE {
				// we use fx.FindMultipleFilter for edge cases when Cocco Lupia might leave the BZ
				// and during the same turn some other old "reduced" used creatures might be returned
				// from other zones to hand, and then re-summoned
				fx.FindMultipleFilter(
					card.Player,
					[]string{match.HAND, match.BATTLEZONE, match.GRAVEYARD, match.MANAZONE, match.SHIELDZONE},
					func(x *match.Card) bool {
						return x.HasCondition(cnd.Creature) && x.SharesAFamily(family.Dragons)
					},
				).Map(func(x *match.Card) {
					x.RemoveConditionBySource(card.ID)
				})

				exit()
				return
			}

			if !ctx2.Match.IsPlayerTurn(card.Player) {
				return
			}

			fx.FindFilter(
				card.Player,
				match.HAND,
				func(x *match.Card) bool {
					return x.HasCondition(cnd.Creature) && x.SharesAFamily(family.Dragons)
				},
			).Map(func(x *match.Card) {
				manaCost := x.ManaCost

				for _, condition := range x.Conditions() {
					if condition.ID == cnd.ReducedCost {
						manaCost -= condition.Val.(int)
						if manaCost < 1 {
							manaCost = 1
						}
					}

					if condition.ID == cnd.IncreasedCost {
						manaCost += condition.Val.(int)
					}
				}

				if manaCost <= 2 {
					return
				}

				subtraction := 2
				if manaCost == 3 {
					subtraction = 1
				}

				x.AddUniqueSourceCondition(cnd.ReducedCost, subtraction, card.ID)
			})
		})
	})

}
