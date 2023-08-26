package dm04

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// GregoriaPrincessOfWar ...
func GregoriaPrincessOfWar(c *match.Card) {

	c.Name = "Gregoria, Princess of War"
	c.Power = 5000
	c.Civ = civ.Darkness
	c.Family = []string{family.DarkLord}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {

			demonCommands := fx.FindFilter(
				card.Player,
				match.BATTLEZONE,
				func(x *match.Card) bool { return x.HasFamily(family.DemonCommand) },
			)

			demonCommands = append(demonCommands,
				fx.FindFilter(
					ctx.Match.Opponent(card.Player),
					match.BATTLEZONE,
					func(x *match.Card) bool { return x.HasFamily(family.DemonCommand) },
				)...,
			)

			if card.Zone != match.BATTLEZONE {
				demonCommands.Map(func(x *match.Card) {
					x.RemoveConditionBySource(card.ID)
				})

				exit()
				return
			}

			demonCommands.Map(func(x *match.Card) {
				x.AddUniqueSourceCondition(cnd.PowerAmplifier, 2000, card.ID)
				x.AddUniqueSourceCondition(cnd.Blocker, true, card.ID)
			})

		})
	}))

}
