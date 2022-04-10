package dm05

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// BallusDogfightEnforcerQ ...
func BallusDogfightEnforcerQ(c *match.Card) {

	c.Name = "Ballus, Dogfight Enforcer Q"
	c.Power = 3000
	c.Civ = civ.Light
	c.Family = family.Berserker
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Survivor, fx.When(fx.EndOfMyTurn, func(card *match.Card, ctx *match.Context) {

		if card.Zone != match.BATTLEZONE {
			return
		}

		fx.FindFilter(
			card.Player,
			match.BATTLEZONE,
			func(x *match.Card) bool { return x.HasCondition(cnd.Survivor) },
		).Map(func(x *match.Card) {
			x.Tapped = false
			ctx.Match.Chat("Server", fmt.Sprintf("%s was untapped by %s's survivor ability", x.Name, card.Name))
		})

	}))

}
