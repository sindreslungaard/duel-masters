package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// ElfX ...
func ElfX(c *match.Card) {

	c.Name = "Elf-X"
	c.Power = 2000
	c.Civ = civ.Nature
	c.Family = []string{family.TreeFolk}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.When(fx.InTheBattlezone, func(card *match.Card, ctx *match.Context) {
		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
			if card.Zone != match.BATTLEZONE {
				// we use fx.FindMultipleFilter for edge cases when elf x might leave the BZ
				// and during the same turn some other old "reduced" used creatures might be returned
				// from other zones to hand, and then re-summoned
				fx.FindMultipleFilter(
					card.Player,
					[]string{match.HAND, match.BATTLEZONE, match.GRAVEYARD, match.MANAZONE, match.SHIELDZONE},
					func(x *match.Card) bool {
						return x.HasCondition(cnd.Creature)
					},
				).Map(func(x *match.Card) {
					x.RemoveConditionBySource(card.ID)
				})

				exit()
				return
			}

			if !ctx.Match.IsPlayerTurn(card.Player) {
				return
			}

			fx.FindFilter(
				card.Player,
				match.HAND,
				func(x *match.Card) bool {
					return x.HasCondition(cnd.Creature)
				},
			).Map(func(x *match.Card) {
				x.AddUniqueSourceCondition(cnd.ReducedCost, 1, card.ID)
			})
		})
	}))

}

// EssenceElf ...
func EssenceElf(c *match.Card) {

	c.Name = "Essence Elf"
	c.Power = 1000
	c.Civ = civ.Nature
	c.Family = []string{family.TreeFolk}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.When(fx.InTheBattlezone, func(card *match.Card, ctx *match.Context) {
		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
			if card.Zone != match.BATTLEZONE {
				// we use fx.FindMultipleFilter for edge cases when essence elf might leave the BZ
				// and during the same turn some other old "reduced" used spell might be returned
				// from grave/mana/shield to hand (i.e. chargers) and then re-casted
				fx.FindMultipleFilter(
					card.Player,
					[]string{match.HAND, match.GRAVEYARD, match.MANAZONE, match.SHIELDZONE},
					func(x *match.Card) bool {
						return x.HasCondition(cnd.Spell)
					},
				).Map(func(x *match.Card) {
					x.RemoveConditionBySource(card.ID)
				})

				exit()
				return
			}

			if !ctx.Match.IsPlayerTurn(card.Player) {
				return
			}

			fx.FindFilter(
				card.Player,
				match.HAND,
				func(x *match.Card) bool {
					return x.HasCondition(cnd.Spell)
				},
			).Map(func(x *match.Card) {
				x.AddUniqueSourceCondition(cnd.ReducedCost, 1, card.ID)
			})
		})
	}))

}
