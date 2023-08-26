package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// StingerWorm ...
func StingerWorm(c *match.Card) {

	c.Name = "Stinger Worm"
	c.Power = 5000
	c.Civ = civ.Darkness
	c.Family = []string{family.ParasiteWorm}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.CardMoved); ok {

			if event.CardID == card.ID && event.To == match.BATTLEZONE {

				creatures := match.Search(card.Player, ctx.Match, card.Player, match.BATTLEZONE, "Stinger Worm: Select 1 creature from your battlezone that will be sent to your graveyard", 1, 1, false)

				for _, creature := range creatures {
					ctx.Match.Destroy(creature, card, match.DestroyedByMiscAbility)
				}

			}

		}

	})

}

// SwampWorm ...
func SwampWorm(c *match.Card) {

	c.Name = "Swamp Worm"
	c.Power = 2000
	c.Civ = civ.Darkness
	c.Family = []string{family.ParasiteWorm}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		ctx.Match.Wait(card.Player, "Waiting for your opponent to make an action")
		defer ctx.Match.EndWait(card.Player)

		creatures := match.Search(ctx.Match.Opponent(card.Player), ctx.Match, ctx.Match.Opponent(card.Player), match.BATTLEZONE, "Swamp Worm: Select 1 creature from your battlezone that will be sent to your graveyard", 1, 1, false)

		for _, creature := range creatures {
			ctx.Match.Destroy(creature, card, match.DestroyedByMiscAbility)
		}

	}))

}
