package dm05

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// LaGuileSeekerOfSkyfire ...
func LaGuileSeekerOfSkyfire(c *match.Card) {

	c.Name = "La Guile, Seeker of Skyfire"
	c.Power = 7500
	c.Civ = civ.Light
	c.Family = []string{family.MechaThunder}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Doublebreaker)

}

// LaByleSeekerOfTheWinds ...
func LaByleSeekerOfTheWinds(c *match.Card) {

	c.Name = "La Byle, Seeker of the Winds"
	c.Power = 5000
	c.Civ = civ.Light
	c.Family = []string{family.MechaThunder}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Blocker, func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.CreatureDestroyed); ok {

			if event.Source == card && event.Blocked {
				card.Tapped = false
			}

		}

	})

}
