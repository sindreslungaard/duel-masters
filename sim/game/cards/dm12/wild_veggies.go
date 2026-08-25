package dm12

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// UncannyTurnip ...
func UncannyTurnip(c *match.Card) {

	c.Name = "Uncanny Turnip"
	c.Power = 1000
	c.Civs = []string{civ.Nature}
	c.Family = []string{family.WildVeggies}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.WaveStriker, fx.WhileWaveStriker(fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		// The card that goes in is not necessarily the one that comes out: the
		// deck feeds the mana zone first, and then any creature there can be
		// taken to hand.
		fx.Draw1ToMana(card, ctx)
		fx.ReturnCreatureFromManazoneToHand(card, ctx)
	})))

}

// FeverNuts ...
func FeverNuts(c *match.Card) {

	c.Name = "Fever Nuts"
	c.Power = 1000
	c.Civs = []string{civ.Nature}
	c.Family = []string{family.WildVeggies}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.When(fx.InTheBattlezone, func(card *match.Card, ctx *match.Context) {
		// Unusually for a cost reducer, this one helps the opponent too. The
		// floor (1, or the number of required civilizations for a multicolored
		// card) is enforced by Card.EffectiveManaCost.
		discount := func(player *match.Player, live bool) {
			fx.FindMultipleFilter(player, []string{match.HAND, match.MANAZONE, match.GRAVEYARD, match.DECK, match.SHIELDZONE},
				func(x *match.Card) bool { return x.HasCondition(cnd.Creature) },
			).Map(func(x *match.Card) {
				if !live {
					x.RemoveSpecificConditionBySource(cnd.ReducedCost, card.ID)
					return
				}

				x.AddUniqueSourceCondition(cnd.ReducedCost, 1, card.ID)
			})
		}

		players := []*match.Player{card.Player, ctx.Match.Opponent(card.Player)}

		for _, player := range players {
			discount(player, true)
		}

		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
			live := card.Zone == match.BATTLEZONE

			for _, player := range players {
				discount(player, live)
			}

			if !live {
				exit()
			}
		})
	}))

}
