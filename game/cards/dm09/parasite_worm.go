package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// VenomWorm ...
func VenomWorm(c *match.Card) {

	c.Name = "Venom Worm"
	c.Power = 1000
	c.Civ = civ.Darkness
	c.Family = []string{family.ParasiteWorm}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Darkness}
	c.TapAbility = venomWormTapAbility

	c.Use(fx.Creature, fx.TapAbility)

}

func venomWormTapAbility(card *match.Card, ctx *match.Context) {
	family := fx.ChooseAFamily(card, ctx, fmt.Sprintf("%s's effect: Choose a race. Each creature of that race gets 'slayer' until the end of the turn", card.Name))

	allCreatures := fx.FindFilter(
		card.Player,
		match.BATTLEZONE,
		func(x *match.Card) bool {
			return x.HasFamily(family)
		})

	allCreatures = append(allCreatures, fx.FindFilter(
		ctx.Match.Opponent(card.Player),
		match.BATTLEZONE,
		func(x *match.Card) bool {
			return x.HasFamily(family)
		})...)

	allCreatures.Map(func(x *match.Card) {
		x.AddUniqueSourceCondition(cnd.Slayer, true, card.ID)
	})

	ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("Every creature in the battlezone of %s race was given 'Slayer' until the end of the turn due to %s's effect.", family, card.Name))
}
