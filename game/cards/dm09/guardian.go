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
			if _, ok := ctx2.Event.(*match.EndOfTurnStep); ok {
				ctx2.ScheduleAfter(func() {
					affectedCreatures := fx.FindFilter(
						card.Player,
						match.BATTLEZONE,
						func(x *match.Card) bool {
							return x.HasFamily(family)
						},
					)

					append(affectedCreatures, fx.FindFilter(
						ctx2.Match.Opponent(card.Player),
						match.BATTLEZONE,
						func(x *match.Card) bool {
							return x.HasFamily(family)
						},
					)...).Map(func(x *match.Card) {
						x.Tapped = false
						ctx2.Match.BroadcastState()
					})

					exit()
				})
			}
		})
	}
}
