package dm07

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func KingBenthos(c *match.Card) {

	c.Name = "King Benthos"
	c.Power = 6000
	c.Civ = civ.Water
	c.Family = []string{family.Leviathan}
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Water}
	c.TapAbility = func(card *match.Card, ctx *match.Context) {
		fx.FindFilter(
			card.Player,
			match.BATTLEZONE,
			func(card *match.Card) bool { return card.Civ == civ.Water },
		).Map(func(x *match.Card) {
			x.AddUniqueSourceCondition(cnd.CantBeBlocked, nil, card.ID)
		})
	}

	c.Use(fx.Creature, fx.Doublebreaker, fx.TapAbility)
}
