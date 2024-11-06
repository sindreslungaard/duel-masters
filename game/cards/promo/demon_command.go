package promo

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// OlgateNightmareSamurai ...
func OlgateNightmareSamurai(c *match.Card) {

	c.Name = "Olgate, Nightmare Samurai"
	c.Power = 6000
	c.Civ = civ.Darkness
	c.Family = []string{family.DemonCommand}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Doublebreaker, fx.When(fx.AnotherOwnCreatureDestroyed, fx.MayUntapSelf))

}

// GiliamTheTormentor ...
func GiliamTheTormentor(c *match.Card) {

	c.Name = "Giliam, the Tormentor"
	c.Power = 5000
	c.Civ = civ.Darkness
	c.Family = []string{family.DemonCommand}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature,
		fx.ConditionalBlocker(func(target *match.Card) bool {
			return target.Civ == civ.Light
		}),
		func(card *match.Card, ctx *match.Context) {

			if event, ok := ctx.Event.(*match.CreatureDestroyed); ok && event.Card == card {

				if event.Context == match.DestroyedInBattle && event.Source.Civ == civ.Light {
					ctx.InterruptFlow()
				}

			}

		})

}
