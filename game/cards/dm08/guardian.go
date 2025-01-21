package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// SolGallaHaloGuardian ...
func SolGallaHaloGuardian(c *match.Card) {

	c.Name = "Sol Galla, Halo Guardian"
	c.Power = 1000
	c.Civ = civ.Light
	c.Family = []string{family.Guardian}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Light}

	powerBoost := 0

	c.PowerModifier = func(m *match.Match, attacking bool) int { return powerBoost }

	c.Use(fx.Creature, fx.Blocker,
		fx.When(fx.SpellCasted, func(card *match.Card, ctx *match.Context) { powerBoost += 3000 }),
		fx.When(fx.EndOfTurn, func(card *match.Card, ctx *match.Context) { powerBoost = 0 }),
	)
}

// ThrumissZephyrGuardian ...
func ThrumissZephyrGuardian(c *match.Card) {
	c.Name = "Thrumiss, Zephyr Guardian"
	c.Power = 3000
	c.Civ = civ.Light
	c.Family = []string{family.Guardian}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature,
		fx.When(fx.InTheBattlezone, func(card *match.Card, ctx *match.Context) {
			ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
				if card.Zone != match.BATTLEZONE {
					fx.Find(
						card.Player,
						match.BATTLEZONE,
					).Map(func(creature *match.Card) {
						creature.RemoveConditionBySource(card.ID)
					})

					exit()
					return
				}

				tapOpponentCreature := func(card *match.Card, ctx *match.Context) {
					fx.Select(
						card.Player,
						ctx.Match,
						ctx.Match.Opponent(card.Player),
						match.BATTLEZONE,
						fmt.Sprintf("%s: you may choose one of your opponent's creatures in the battle zone and tap it.", c.Name),
						1,
						1,
						true,
					).Map(func(creature *match.Card) {
						creature.Tapped = true
						fx.RemoveBlockerFromList(creature, ctx)
					})
				}

				fx.Find(
					card.Player,
					match.BATTLEZONE,
				).Map(func(creature *match.Card) {
					creature.AddUniqueSourceCondition(cnd.WheneverThisAttacks, tapOpponentCreature, card)
				})
			})
		}),
	)
}
