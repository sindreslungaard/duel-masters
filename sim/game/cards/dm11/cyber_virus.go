package dm11

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// LuckyBall ...
func LuckyBall(c *match.Card) {

	c.Name = "Lucky Ball"
	c.Power = 3000
	c.Civs = []string{civ.Water}
	c.Family = []string{family.CyberVirus}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		shields, err := ctx.Match.Opponent(card.Player).Container(match.SHIELDZONE)

		if err != nil || len(shields) > 3 {
			return
		}

		fx.DrawUpTo2(card, ctx)
	}))

}
