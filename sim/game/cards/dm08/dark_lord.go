package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// MegariaEmpressOfDread ...
func MegariaEmpressOfDread(c *match.Card) {

	c.Name = "Megaria, Empress of Dread"
	c.Power = 5000
	c.Civ = civ.Darkness
	c.Family = []string{family.DarkLord}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.When(fx.InTheBattlezone, func(card *match.Card, ctx *match.Context) {
		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
			if card.Zone != match.BATTLEZONE {
				fx.Find(
					card.Player,
					match.BATTLEZONE,
				).Map(func(x *match.Card) {
					x.RemoveConditionBySource(card.ID)
				})

				fx.Find(
					ctx2.Match.Opponent(card.Player),
					match.BATTLEZONE,
				).Map(func(x *match.Card) {
					x.RemoveConditionBySource(card.ID)
				})

				exit()
				return
			}

			fx.Find(
				card.Player,
				match.BATTLEZONE,
			).Map(func(x *match.Card) {
				x.AddUniqueSourceCondition(cnd.Slayer, true, card.ID)
			})

			fx.Find(
				ctx2.Match.Opponent(card.Player),
				match.BATTLEZONE,
			).Map(func(x *match.Card) {
				x.AddUniqueSourceCondition(cnd.Slayer, true, card.ID)
			})
		})
	}))

}
