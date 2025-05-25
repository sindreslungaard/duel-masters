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

	c.Use(fx.Creature, fx.CantBeSelectedByOpp, //TODO implement filter in Select fxs
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
							ctx.Match.Opponent(card.Player),
							match.BATTLEZONE,
							func(x *match.Card) bool {
								return x.HasFamily(chosenFamily)
							},
						).Map(func(x *match.Card) {
							x.RemoveConditionBySource(card.ID)
						})

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
						ctx.Match.Opponent(card.Player),
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

// BruiserDragon ...
func BruiserDragon(c *match.Card) {

	c.Name = "Bruiser Dragon"
	c.Power = 5000
	c.Civ = civ.Fire
	c.Family = []string{family.ArmoredDragon}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.When(fx.Destroyed, func(card *match.Card, ctx *match.Context) {
		fx.SelectBackside(
			card.Player,
			ctx.Match,
			card.Player,
			match.SHIELDZONE,
			fmt.Sprintf("%s's effect: Put 1 of your shields into your graveyard.", card.Name),
			1,
			1,
			false,
		).Map(func(x *match.Card) {
			ctx.Match.MoveCard(x, match.GRAVEYARD, card)
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s's effect: One of your shields was put into graveyard", card.Name))
		})
	}))

}
