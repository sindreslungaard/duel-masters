package dm04

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// ThreeEyedDragonfly ...
func ThreeEyedDragonfly(c *match.Card) {

	c.Name = "Three-Eyed Dragonfly"
	c.Power = 4000
	c.Civ = civ.Nature
	c.Family = family.GiantInsect
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.When(fx.Attacking, func(card *match.Card, ctx *match.Context) {

		selected := fx.SelectFilter(
			card.Player,
			ctx.Match,
			card.Player,
			match.BATTLEZONE,
			"Select one creature to destroy.",
			0,
			1,
			true,
			func(c *match.Card) bool { return c.ID != card.ID },
		)

		if len(selected) > 0 {
			selected.Map(func(x *match.Card) {
				ctx.Match.Destroy(x, card, match.DestroyedByMiscAbility)
				ctx.Match.Chat("Server", fmt.Sprintf("%s was destroyed by Three-Eyed Dragonfly", x.Name))
			})

			card.AddCondition(cnd.PowerAttacker, 2000, card.ID)
			card.AddCondition(cnd.DoubleBreaker, nil, card.ID)
			card.AddCondition(cnd.ActionNotCancellable, nil, card.ID)
			ctx.Match.Chat("Server", fmt.Sprintf("%s was given power attacker +2000 and double breaker until the end of the turn", card.Name))

		}

	}))

}
