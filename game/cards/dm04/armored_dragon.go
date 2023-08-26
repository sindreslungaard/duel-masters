package dm04

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// GalklifeDragon ...
func GalklifeDragon(c *match.Card) {

	c.Name = "Galklife Dragon"
	c.Power = 6000
	c.Civ = civ.Fire
	c.Family = []string{family.ArmoredDragon}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.Doublebreaker, func(card *match.Card, ctx *match.Context) {

		if match.AmISummoned(card, ctx) {

			fx.FindFilter(
				card.Player,
				match.BATTLEZONE,
				func(x *match.Card) bool { return ctx.Match.GetPower(x, false) <= 4000 && x.Civ == civ.Light },
			).Map(func(x *match.Card) {
				ctx.Match.Destroy(x, card, match.DestroyedByMiscAbility)
				ctx.Match.Chat("Server", fmt.Sprintf("%s was destroyed by Galklife Dragon", x.Name))
			})

			fx.FindFilter(
				ctx.Match.Opponent(card.Player),
				match.BATTLEZONE,
				func(x *match.Card) bool { return ctx.Match.GetPower(x, false) <= 4000 && x.Civ == civ.Light },
			).Map(func(x *match.Card) {
				ctx.Match.Destroy(x, card, match.DestroyedByMiscAbility)
				ctx.Match.Chat("Server", fmt.Sprintf("%s was destroyed by Galklife Dragon", x.Name))
			})
		}
	})
}
