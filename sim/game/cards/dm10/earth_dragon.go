package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// TerradragonCusdalf ...
func TerradragonCusdalf(c *match.Card) {
	c.Name = "Terradragon Cusdalf"
	c.Power = 7000
	c.Civs = []string{civ.Nature}
	c.Family = []string{family.EarthDragon}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Nature}

	c.Use(
		fx.Creature,
		fx.Doublebreaker,
		fx.PowerAttacker4000,
		func(card *match.Card, ctx *match.Context) {
			if event, ok := ctx.Event.(*match.UntapManaEvent); ok &&
				card.Zone == match.BATTLEZONE &&
				event.Player == card.Player {
				ctx.InterruptFlow()
			}
		},
	)
}

// TerradragonDakmaBalgarow ...
func TerradragonDakmaBalgarow(c *match.Card) {
	c.Name = "Terradragon Dakma Balgarow"
	c.Power = 1000
	c.Civs = []string{civ.Nature}
	c.Family = []string{family.EarthDragon}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Nature}

	c.PowerModifier = func(m *match.Match, attacking bool) int {
		if c.Zone != match.BATTLEZONE {
			return 0
		}

		shields := len(fx.Find(c.Player, match.SHIELDZONE)) + len(fx.Find(m.Opponent(c.Player), match.SHIELDZONE))

		return shields * 2000
	}

	c.Use(fx.Creature, fx.PowerBreakerTiers(6000, 15000))
}

// UltimateDragon ...
func UltimateDragon(c *match.Card) {
	c.Name = "Ultimate Dragon"
	c.Power = 5000
	c.Civs = []string{civ.Fire}
	c.Family = []string{family.ArmoredDragon}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Fire}

	otherDragons := fx.CountOtherOwnCreaturesWithFamily(family.Dragons)

	c.PowerModifier = func(m *match.Match, attacking bool) int {
		if c.Zone != match.BATTLEZONE {
			return 0
		}

		return otherDragons(c) * 5000
	}

	// Crew breaker-Dragon
	c.Use(fx.Creature, fx.CrewBreaker(otherDragons))
}
