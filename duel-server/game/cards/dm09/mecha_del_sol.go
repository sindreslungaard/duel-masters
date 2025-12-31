package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// PetrovaChannelerOfSuns ...
func PetrovaChannelerOfSuns(c *match.Card) {

	c.Name = "Petrova, Channeler of Suns"
	c.Power = 3500
	c.Civ = civ.Light
	c.Family = []string{family.MechaDelSol}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.CantBeSelectedByOpp,
		fx.When(fx.InTheBattlezone, func(card *match.Card, ctx *match.Context) {
			chosenFamily := fx.ChooseAFamilyFilter(
				card,
				ctx,
				fmt.Sprintf("%s's effect: Choose a race other than Mecha del Sol. Each creature of that race gets +4000 Power.", card.Name),
				func(x string) bool {
					return x != family.MechaDelSol
				},
			)

			if chosenFamily != "" {
				ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("All creatures in the battlezone of race %s were given +4000 power by %s", chosenFamily, card.Name))

				ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
					if card.Zone != match.BATTLEZONE {
						fx.FindFilter(
							card.Player,
							match.BATTLEZONE,
							func(x *match.Card) bool {
								return x.HasFamily(chosenFamily)
							},
						).Map(func(x *match.Card) {
							x.RemoveConditionBySource(card.ID)
						})

						fx.FindFilter(
							ctx2.Match.Opponent(card.Player),
							match.BATTLEZONE,
							func(x *match.Card) bool {
								return x.HasFamily(chosenFamily)
							},
						).Map(func(x *match.Card) {
							x.RemoveConditionBySource(card.ID)
						})

						ctx2.Match.ReportActionInChat(card.Player, fmt.Sprintf("All creatures in the battlezone of race %s NO more have +4000 power", chosenFamily))

						exit()
						return
					}

					fx.FindFilter(
						card.Player,
						match.BATTLEZONE,
						func(x *match.Card) bool {
							return x.HasFamily(chosenFamily)
						},
					).Map(func(x *match.Card) {
						x.AddUniqueSourceCondition(cnd.PowerAmplifier, 4000, card.ID)
					})

					fx.FindFilter(
						ctx2.Match.Opponent(card.Player),
						match.BATTLEZONE,
						func(x *match.Card) bool {
							return x.HasFamily(chosenFamily)
						},
					).Map(func(x *match.Card) {
						x.AddUniqueSourceCondition(cnd.PowerAmplifier, 4000, card.ID)
					})
				})
			}
		}),
	)

}
