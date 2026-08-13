package dm12

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// Cloned Spike-Horn grows by this much for every copy of itself lying in a
// graveyard, and breaks two shields once that has carried it far enough.
const (
	clonedSpikeHornPowerPerCopy    = 3000
	clonedSpikeHornDoubleBreakerAt = 6000
)

// RadioactiveHornTheStrange ...
func RadioactiveHornTheStrange(c *match.Card) {

	c.Name = "Radioactive Horn, the Strange"
	c.Power = 1000
	c.Civs = []string{civ.Nature}
	c.Family = []string{family.HornedBeast}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.Doublebreaker)

}

// SpectralHornGlitalis ...
func SpectralHornGlitalis(c *match.Card) {

	c.Name = "Spectral Horn Glitalis"
	c.Power = 4000
	c.Civs = []string{civ.Nature, civ.Light}
	c.Family = []string{family.HornedBeast, family.RainbowPhantom}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Nature, civ.Light}

	c.Use(fx.Creature, fx.PutIntoManaZoneTapped)

}

// ClonedSpikeHorn ...
func ClonedSpikeHorn(c *match.Card) {

	c.Name = "Cloned Spike-Horn"
	c.Power = 3000
	c.Civs = []string{civ.Nature}
	c.Family = []string{family.HornedBeast}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.PowerModifier = func(m *match.Match, attacking bool) int {
		return fx.ClonedCopiesInGraveyards(c, m) * clonedSpikeHornPowerPerCopy
	}

	c.Use(fx.Creature, fx.PowerBreakerTiers(clonedSpikeHornDoubleBreakerAt, 0))

}

// TurtleHornTheImposing ...
func TurtleHornTheImposing(c *match.Card) {

	c.Name = "Turtle Horn, the Imposing"
	c.Power = 2000
	c.Civs = []string{civ.Nature}
	c.Family = []string{family.HornedBeast}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Nature}

	// The search is mandatory but taking a creature is not, which is what
	// SearchDeckTake1Creature already offers, along with showing the card and
	// shuffling afterwards.
	c.Use(fx.Creature, fx.When(fx.OpponentPlayedShieldTrigger, func(card *match.Card, ctx *match.Context) {
		if card.Zone != match.BATTLEZONE {
			return
		}

		fx.SearchDeckTake1Creature(card, ctx)
	}))

}
