package promo

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// NeveTheLeveler ...
func NeveTheLeveler(c *match.Card) {

	c.Name = "Neve, the Leveler"
	c.Power = 4000
	c.Civs = []string{civ.Nature}
	c.Family = []string{family.SnowFaerie}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		myCreatures := fx.Find(card.Player, match.BATTLEZONE)
		oppCreatures := fx.Find(ctx.Match.Opponent(card.Player), match.BATTLEZONE)

		extra := len(oppCreatures) - len(myCreatures)
		if extra < 1 {
			return
		}

		fx.SearchDeckTakeCards(
			card,
			ctx,
			extra,
			func(x *match.Card) bool { return x.HasCondition(cnd.Creature) },
			"creature(s)",
		)
	}))

}
