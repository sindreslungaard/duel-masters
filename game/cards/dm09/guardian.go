package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// TraRionPenumbraGuardian ...
func TraRionPenumbraGuardian(c *match.Card) {

	c.Name = "Tra Rion, Penumbra Guardian"
	c.Power = 5500
	c.Civ = civ.Light
	c.Family = []string{family.Guardian}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Light}
	c.TapAbility = traRionPenumbraGuardianTapAbility

	c.Use(fx.Creature, fx.TapAbility)

}

func traRionPenumbraGuardianTapAbility(card *match.Card, ctx *match.Context) {
	family := fx.ChooseAFamily(
		card,
		ctx,
		fmt.Sprintf("%s's effect: Choose a race. At the end of this turn, untap all creatures of that race in the battlezone.", card.Name),
	)

	if family != "" {
		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
			affectedCreatures := fx.FindFilter(
				card.Player,
				match.BATTLEZONE,
				func(x *match.Card) bool {
					return x.HasFamily(family) && !x.HasCondition(card.ID+"-custom")
				},
			)

			affectedCreatures = append(affectedCreatures, fx.FindFilter(
				ctx2.Match.Opponent(card.Player),
				match.BATTLEZONE,
				func(x *match.Card) bool {
					return x.HasFamily(family) && !x.HasCondition(card.ID+"-custom")
				},
			)...).Map(func(x *match.Card) {
				ctx2.Match.ApplyPersistentEffect(func(ctx3 *match.Context, exit2 func()) {
					x.AddUniqueSourceCondition(card.ID+"-custom", true, card.ID)

					if _, ok := ctx3.Event.(*match.EndOfTurnStep); ok {
						x.RemoveConditionBySource(card.ID)

						x.Tapped = false

						exit2()
						return
					}
				})
			})

			if _, ok := ctx2.Event.(*match.EndOfTurnStep); ok {
				affectedCreatures.Map(func(x *match.Card) {
					x.RemoveConditionBySource(card.ID)
				})
				exit()
				return
			}
		})
	}
}
