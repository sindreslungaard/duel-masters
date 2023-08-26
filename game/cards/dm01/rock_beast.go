package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// Meteosaur ...
func Meteosaur(c *match.Card) {

	c.Name = "Meteosaur"
	c.Power = 2000
	c.Civ = civ.Fire
	c.Family = []string{family.RockBeast}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.CardMoved); ok {

			if event.CardID == card.ID && event.To == match.BATTLEZONE {

				creatures := match.Filter(
					card.Player,
					ctx.Match,
					ctx.Match.Opponent(card.Player),
					match.BATTLEZONE,
					"Meteosaur: Select 1 of your opponent's creatures with power 2000 or less and destroy it",
					1,
					1,
					true,
					func(x *match.Card) bool { return ctx.Match.GetPower(x, false) <= 2000 },
				)

				for _, creature := range creatures {
					ctx.Match.Destroy(creature, card, match.DestroyedByMiscAbility)
				}

			}
		}

	})

}

// Stonesaur ...
func Stonesaur(c *match.Card) {

	c.Name = "Stonesaur"
	c.Power = 4000
	c.Civ = civ.Fire
	c.Family = []string{family.RockBeast}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.PowerAttacker2000)

}
