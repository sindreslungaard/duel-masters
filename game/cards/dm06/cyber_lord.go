package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

func Sopian(c *match.Card) {

	c.Name = "Sopian"
	c.Power = 2000
	c.Civ = civ.Water
	c.Family = []string{family.CyberLord}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Water}
	c.TapAbility = func(card *match.Card, ctx *match.Context) {

		creatures := match.Search(card.Player, ctx.Match, card.Player, match.BATTLEZONE, "Select 1 creature from your battlezone that will gain \"Can't be blocked this turn\"", 1, 1, false)
		for _, creature := range creatures {

			creature.AddCondition(cnd.CantBeBlocked, 1, card.ID)
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was given \"Cant be blocked this turn by %s\"", creature.Name, card.Name))

		}
	}

	c.Use(fx.Creature, fx.TapAbility)
}
