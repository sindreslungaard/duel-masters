package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// StainedGlass ...
func StainedGlass(c *match.Card) {

	c.Name = "Stained Glass"
	c.Power = 1000
	c.Civ = civ.Water
	c.Family = []string{family.CyberVirus}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.When(fx.AttackConfirmed, func(card *match.Card, ctx *match.Context) {

		fx.SelectFilter(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			fmt.Sprintf("%s: You may select 1 of your opponent's fire or nature creatures that will be returned to their hand", card.Name),
			1,
			1,
			true,
			func(x *match.Card) bool { return x.Civ == civ.Fire || x.Civ == civ.Nature },
			false,
		).Map(func(x *match.Card) {
			x.Player.MoveCard(x.ID, match.BATTLEZONE, match.HAND, card.ID)
			ctx.Match.ReportActionInChat(ctx.Match.Opponent(card.Player), fmt.Sprintf("%s was sent to %s's hand from the battle zone by %s's effect", x.Name, x.Player.Username(), card.Name))
		})

	}))

}
