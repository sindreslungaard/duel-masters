package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// KingOquanos ...
func KingOquanos(c *match.Card) {

	c.Name = "King Oquanos"
	c.Power = 2000
	c.Civ = civ.Water
	c.Family = []string{family.Leviathan}
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {
		if event, ok := ctx.Event.(*match.GetPowerEvent); ok {
			if card.Zone == match.BATTLEZONE && event.Card == card {
				event.Power += len(fx.FindFilter(
					ctx.Match.Opponent(card.Player),
					match.MANAZONE,
					func(x *match.Card) bool {
						return x.Tapped
					},
				)) * 2000
			}
		}
	})

}
