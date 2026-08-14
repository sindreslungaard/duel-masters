package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// CarnivalTotem ...
func CarnivalTotem(c *match.Card) {
	c.Name = "Carnival Totem"
	c.Power = 7000
	c.Civs = []string{civ.Nature}
	c.Family = []string{family.MysteryTotem}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.Doublebreaker, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		fx.SwapHandAndMana(card, card.Player)
		ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s's effect: %s swapped the cards in their hand and mana zone.", card.Name, card.Player.Username()))
	}))
}

// JigglyTotem ...
func JigglyTotem(c *match.Card) {

	c.Name = "Jiggly Totem"
	c.Power = 1000
	c.Civs = []string{civ.Nature}
	c.Family = []string{family.MysteryTotem}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.PowerModifier = func(m *match.Match, attacking bool) int {
		if attacking {
			return len(fx.FindFilter(
				c.Player,
				match.MANAZONE,
				func(x *match.Card) bool {
					return x.Tapped
				},
			)) * 1000
		}

		return 0
	}

	c.Use(fx.Creature)

}

// TechnoTotem ...
func TechnoTotem(c *match.Card) {

	c.Name = "Techno Totem"
	c.Power = 5000
	c.Civs = []string{civ.Light, civ.Nature}
	c.Family = []string{family.MysteryTotem}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light, civ.Nature}

	c.TapAbility = func(card *match.Card, ctx *match.Context) {
		fx.TapOpCreature(card, ctx)
	}

	c.Use(fx.Creature, fx.PutIntoManaZoneTapped, fx.TapAbility, fx.When(fx.InTheBattlezone, func(card *match.Card, ctx *match.Context) {
		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
			// The grant is removed whenever it stops applying, which includes the
			// creature untapping again, so it is recomputed on every event.
			fx.FindFilter(
				card.Player,
				match.BATTLEZONE,
				func(x *match.Card) bool { return x.ID != card.ID },
			).Map(func(x *match.Card) {
				x.RemoveSpecificConditionBySource(cnd.PowerAttacker, card.ID)
			})

			if card.Zone != match.BATTLEZONE {
				exit()
				return
			}

			if !card.Tapped {
				return
			}

			fx.FindFilter(
				card.Player,
				match.BATTLEZONE,
				func(x *match.Card) bool { return x.ID != card.ID && x.HasCondition(cnd.Creature) },
			).Map(func(x *match.Card) {
				x.AddUniqueSourceCondition(cnd.PowerAttacker, 1500, card.ID)
			})
		})
	}))

}
