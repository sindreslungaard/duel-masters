package dm03

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// Flametropus ...
func Flametropus(c *match.Card) {

	c.Name = "Flametropus"
	c.Power = 3000
	c.Civ = civ.Fire
	c.Family = []string{family.RockBeast}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.When(fx.Attacking, func(card *match.Card, ctx *match.Context) {

		fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.MANAZONE,
			"Flametropus: You may select 1 card from your mana zone that will be sent to your graveyard.",
			1,
			1,
			true,
		).Map(func(x *match.Card) {

			card.Player.MoveCard(x.ID, match.MANAZONE, match.GRAVEYARD)
			ctx.Match.Chat("Server", fmt.Sprintf("%s was moved from %s manazone to their graveyard by Flametropus", x.Name, card.Player.Username()))

			card.AddCondition(cnd.DoubleBreaker, nil, card.ID)
			card.AddCondition(cnd.PowerAttacker, 3000, card.ID)
			ctx.Match.Chat("Server", fmt.Sprintf("%s was given power attacker +3000 and double breaker until the end of the turn", card.Name))

		})

	}))

}
