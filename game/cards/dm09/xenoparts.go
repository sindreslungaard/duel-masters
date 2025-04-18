package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// GigiosHammer ...
func GigiosHammer(c *match.Card) {

	c.Name = "Gigio's Hammer"
	c.Power = 2000
	c.Civ = civ.Fire
	c.Family = []string{family.Xenoparts}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Fire}
	c.TapAbility = gigiosHammerTapAbility

	c.Use(fx.Creature, fx.TapAbility)
}

func gigiosHammerTapAbility(card *match.Card, ctx *match.Context) {
	family := fx.ChooseAFamily(card, ctx, fmt.Sprintf("%s's effect: Choose a race. Each creature of that race attacks this turn if able and gets 'Power attacker +4000'.", card.Name))

	fx.FindFilter(
		card.Player,
		match.BATTLEZONE,
		func(x *match.Card) bool {
			return x.HasFamily(family) && x.HasCondition(cnd.Creature)
		},
	).Map(func(x *match.Card) {
		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
			if _, ok := ctx2.Event.(*match.EndOfTurnStep); ok {
				x.RemoveConditionBySource(card.ID)
				exit()
				return
			}

			x.AddUniqueSourceCondition(cnd.PowerAttacker, 4000, card.ID)
			fx.ForceAttack(x, ctx2)
		})
	})
}
