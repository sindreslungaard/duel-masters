package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// MagmadragonMelgars ...
func MagmadragonMelgars(c *match.Card) {

	c.Name = "Magmadragon Melgars"
	c.Power = 4000
	c.Civ = civ.Fire
	c.Family = []string{family.VolcanoDragon}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature)

}

// MagmadragonJagalzor ...
func MagmadragonJagalzor(c *match.Card) {

	c.Name = "Magmadragon Jagalzor"
	c.Power = 6000
	c.Civ = civ.Fire
	c.Family = []string{family.VolcanoDragon}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Fire}

	turboRush := false

	c.Use(fx.Creature, fx.Doublebreaker,
		fx.When(fx.TurboRushCondition, func(card *match.Card, ctx *match.Context) { turboRush = true }),
		fx.When(fx.EndOfMyTurn, func(card *match.Card, ctx *match.Context) { turboRush = false }),
		fx.When(fx.InTheBattlezone, func(card *match.Card, ctx *match.Context) {
			ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
				if card.Zone != match.BATTLEZONE {
					fx.Find(
						card.Player,
						match.BATTLEZONE,
					).Map(func(x *match.Card) {
						x.RemoveConditionBySource(card.ID)
					})

					exit()
					return
				}

				if turboRush {
					fx.Find(card.Player, match.BATTLEZONE).Map(func(x *match.Card) { x.AddUniqueSourceCondition(cnd.SpeedAttacker, true, card.ID) })
				}
			})
		}),
	)

}
