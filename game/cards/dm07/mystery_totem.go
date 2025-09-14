package dm07

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func CrypticTotem(c *match.Card) {

	c.Name = "Cryptic Totem"
	c.Power = 6000
	c.Civ = civ.Nature
	c.Family = []string{family.MysteryTotem}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.Doublebreaker, fx.When(fx.IsTapped, func(card *match.Card, ctx *match.Context) {
		fx.FilterShieldTriggers(ctx, func(x *match.Card) bool { return card.Player == x.Player })
	}))
}

func SpinningTotem(c *match.Card) {

	c.Name = "Spinning Totem"
	c.Power = 4000
	c.Civ = civ.Nature
	c.Family = []string{family.MysteryTotem}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Nature}
	c.TapAbility = func(card *match.Card, ctx *match.Context) {
		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
			if event, ok := ctx2.Event.(*match.Battle); ok {
				if !event.Blocked || event.Attacker.Player != card.Player {
					return
				}

				if event.Attacker.Civ != civ.Nature {
					return
				}

				fx.DestroyOpShield(card, ctx2)
			}

			// remove persistent effect when turn ends
			_, ok := ctx2.Event.(*match.EndOfTurnStep)
			if ok {
				exit()
			}
		})
	}

	c.Use(fx.Creature, fx.TapAbility)
}
