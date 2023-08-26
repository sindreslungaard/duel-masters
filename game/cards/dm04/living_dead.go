package dm04

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// SkeletonThiefTheRevealer ...
func SkeletonThiefTheRevealer(c *match.Card) {

	c.Name = "Skeleton Thief, the Revealer"
	c.Power = 2000
	c.Civ = civ.Darkness
	c.Family = []string{family.LivingDead}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		fx.SelectFilter(
			card.Player,
			ctx.Match,
			card.Player,
			match.GRAVEYARD,
			"Skeleton Thief, the Revealer: You may return a Living Dead from your graveyard to your hand",
			1,
			1,
			true,
			func(x *match.Card) bool { return x.HasFamily(family.LivingDead) },
		).Map(func(x *match.Card) {
			card.Player.MoveCard(x.ID, match.GRAVEYARD, match.HAND)
			ctx.Match.Chat("Server", fmt.Sprintf("%s was moved to %s's hand from their graveyard by Skeleton Thief, the Revealer", x.Name, card.Player.Username()))
		})

	}))

}
