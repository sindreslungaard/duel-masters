package dm03

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// SnipStrikerBullraizer ...
func SnipStrikerBullraizer(c *match.Card) {

	c.Name = "Snip Striker Bullraizer"
	c.Power = 3000
	c.Civ = civ.Fire
	c.Family = family.Dragonoid
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.When(fx.AttackingPlayer, func(card *match.Card, ctx *match.Context) {

		creatures, err := card.Player.Container(match.BATTLEZONE)

		if err != nil {
			return
		}

		oppCreatures, err := ctx.Match.Opponent(card.Player).Container(match.BATTLEZONE)

		if err != nil {
			return
		}

		if len(creatures) < len(oppCreatures) {

			ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("%s can't attack when the opponent has more creatures in the battle zone than you.", card.Name))

			ctx.InterruptFlow()
		}
	}))

}
