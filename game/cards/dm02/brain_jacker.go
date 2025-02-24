package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// AmberPiercer ...
func AmberPiercer(c *match.Card) {

	c.Name = "Amber Piercer"
	c.Power = 2000
	c.Civ = civ.Darkness
	c.Family = []string{family.BrainJacker}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.When(fx.AttackConfirmed, func(card *match.Card, ctx *match.Context) {

		fx.SelectFilter(
			card.Player,
			ctx.Match,
			card.Player,
			match.GRAVEYARD,
			"You may select a creature to return from your graveyard to your hand.",
			1,
			1,
			true,
			func(x *match.Card) bool { return x.HasCondition(cnd.Creature) },
			false,
		).Map(func(x *match.Card) {
			card.Player.MoveCard(x.ID, match.GRAVEYARD, match.HAND, card.ID)
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was sent to %s's hand from their graveyard by %s's effect", x.Name, card.Player.Username(), card.Name))
		})

	}))

}
