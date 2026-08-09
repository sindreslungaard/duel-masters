package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// BatDoctorShadowOfUndeath ...
func BatDoctorShadowOfUndeath(c *match.Card) {

	c.Name = "Bat Doctor, Shadow of Undeath"
	c.Power = 2000
	c.Civs = []string{civ.Darkness}
	c.Family = []string{family.Ghost}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.When(fx.Destroyed, func(card *match.Card, ctx *match.Context) {
		ctx.ScheduleAfter(func() {
			fx.SelectFilter(
				card.Player,
				ctx.Match,
				card.Player,
				match.GRAVEYARD,
				fmt.Sprintf("%s's effect: You may return another creature from your graveyard to your hand.", card.Name),
				1,
				1,
				true,
				func(x *match.Card) bool {
					return x.ID != card.ID && x.HasCondition(cnd.Creature)
				},
				false,
			).Map(func(x *match.Card) {
				card.Player.MoveCard(x.ID, match.GRAVEYARD, match.HAND, card.ID)
				ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was returned to %s's hand from its graveyard by %s's effect.", x.Name, card.Player.Username(), card.Name))
			})
		})
	}))

}

// IceVaporShadowOfAnguish ...
func IceVaporShadowOfAnguish(c *match.Card) {

	c.Name = "Ice Vapor, Shadow of Anguish"
	c.Power = 1000
	c.Civs = []string{civ.Darkness}
	c.Family = []string{family.Ghost}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Darkness}

	// Both choices belong to the opponent, and both are mandatory when they have
	// a card to give up.
	c.Use(fx.Creature, fx.When(fx.OpponentCastASpell, func(card *match.Card, ctx *match.Context) {
		fx.OpDiscardsXCards(1)(card, ctx)
		fx.OpponentChoosesManaBurn(card, ctx)
	}))

}
