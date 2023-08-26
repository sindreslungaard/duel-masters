package dm04

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// DoboulgyserGiantRockBeast ...
func DoboulgyserGiantRockBeast(c *match.Card) {

	c.Name = "Doboulgyser, Giant Rock Beast"
	c.Power = 8000
	c.Civ = civ.Fire
	c.Family = []string{family.RockBeast}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.Evolution, fx.Doublebreaker, func(card *match.Card, ctx *match.Context) {

		if match.AmISummoned(card, ctx) {

			fx.SelectFilter(
				card.Player,
				ctx.Match,
				ctx.Match.Opponent(card.Player),
				match.BATTLEZONE,
				"Doboulgyser: You may select 1 opponent's creature with 3000 or less power and destroy it.",
				0,
				1,
				true,
				func(x *match.Card) bool { return ctx.Match.GetPower(x, false) <= 3000 },
			).Map(func(x *match.Card) {

				ctx.Match.Destroy(x, card, match.DestroyedByMiscAbility)
				ctx.Match.Chat("Server", fmt.Sprintf("%s was destroyed by Doboulgyser", x.Name))
			})

		}
	})

}

// Magmarex ...
func Magmarex(c *match.Card) {

	c.Name = "Magmarex"
	c.Power = 3000
	c.Civ = civ.Fire
	c.Family = []string{family.RockBeast}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.ShieldTrigger, func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.CardMoved); ok {

			if event.CardID != card.ID || event.To != match.BATTLEZONE {
				return
			}

			fx.FindFilter(
				card.Player,
				match.BATTLEZONE,
				func(x *match.Card) bool { return ctx.Match.GetPower(x, false) == 1000 },
			).Map(func(x *match.Card) {
				ctx.Match.Destroy(x, card, match.DestroyedByMiscAbility)
				ctx.Match.Chat("Server", fmt.Sprintf("%s was destroyed by Magmarex", x.Name))
			})

			fx.FindFilter(
				ctx.Match.Opponent(card.Player),
				match.BATTLEZONE,
				func(x *match.Card) bool { return ctx.Match.GetPower(x, false) == 1000 },
			).Map(func(x *match.Card) {
				ctx.Match.Destroy(x, card, match.DestroyedByMiscAbility)
				ctx.Match.Chat("Server", fmt.Sprintf("%s was destroyed by Magmarex", x.Name))
			})

		}
	})

}