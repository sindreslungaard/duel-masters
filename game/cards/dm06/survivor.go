package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
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

		fx.DrawBetween(card, ctx, 0, toDraw)

	}))
}

func QTronicGargantua(c *match.Card) {

	c.Name = "Q-tronic Gargantua"
	c.Power = 9000
	c.Civ = civ.Fire
	c.Family = []string{family.Survivor}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.Evolution, fx.Survivor, fx.When(fx.AttackingPlayer, func(card *match.Card, ctx *match.Context) {

		mysurvivors := len(fx.FindFilter(card.Player, match.BATTLEZONE, func(card *match.Card) bool { return card.HasFamily(family.Survivor) }))

		if mysurvivors < 1 {
			return
		}

		card.AddUniqueSourceCondition(cnd.ShieldBreakModifier, mysurvivors-1, card.ID)
		ctx.ScheduleAfter(func() {
			card.RemoveConditionBySource(card.ID)
		})

	}))
}
