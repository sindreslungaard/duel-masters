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
	c.Civs = []string{civ.Darkness}
	c.Family = []string{family.DemonCommand}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Doublebreaker, fx.When(fx.AnotherOwnCreatureDestroyed, fx.MayUntapSelf))

}

// GiliamTheTormentor ...
func GiliamTheTormentor(c *match.Card) {

	c.Name = "Giliam, the Tormentor"
	c.Power = 5000
	c.Civs = []string{civ.Darkness}
	c.Family = []string{family.DemonCommand}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature,
		fx.LightBlocker(),
		func(card *match.Card, ctx *match.Context) {

			if event, ok := ctx.Event.(*match.CreatureDestroyed); ok && event.Card == card {
				if event.Context == match.DestroyedInBattle && event.Source.HasCiv(civ.Light) {
					ctx.InterruptFlow()
				}
			}

		})

}
