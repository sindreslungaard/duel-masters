package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func QTronicHypermind(c *match.Card) {

	c.Name = "Q-tronic Hypermind"
	c.Power = 8000
	c.Civ = civ.Water
	c.Family = []string{family.Survivor}
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Doublebreaker, fx.Evolution, fx.Survivor, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		mysurvivors := len(fx.FindFilter(card.Player, match.BATTLEZONE, func(card *match.Card) bool { return card.HasFamily(family.Survivor) }))
		enemysurvivors := len(fx.FindFilter(ctx.Match.Opponent(card.Player), match.BATTLEZONE, func(card *match.Card) bool { return card.HasFamily(family.Survivor) }))
		toDraw := mysurvivors + enemysurvivors

		if toDraw < 1 {
			return
		}

		card.Player.DrawCards(toDraw)

	}))
}
