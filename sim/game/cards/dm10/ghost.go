package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// MummyWrapShadowOfFatigue ...
func MummyWrapShadowOfFatigue(c *match.Card) {

	c.Name = "Mummy Wrap, Shadow of Fatigue"
	c.Power = 1000
	c.Civs = []string{civ.Darkness}
	c.Family = []string{family.Ghost}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Darkness}
	c.TapAbility = func(card *match.Card, ctx *match.Context) {
		fx.PlayerDiscardsRandomCard(card, ctx)
		fx.OpponentDiscardsRandomCard(card, ctx)
	}

	c.Use(fx.Creature, fx.TapAbility)

}

// SparkChemistShadowOfWhim ...
func SparkChemistShadowOfWhim(c *match.Card) {

	c.Name = "Spark Chemist, Shadow of Whim"
	c.Power = 3000
	c.Civs = []string{civ.Darkness}
	c.Family = []string{family.Ghost}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.When(fx.Summoned, fx.ReturnAllCardsFromManaToHand))

}

// ZeroNemesisShadowOfPanic ...
func ZeroNemesisShadowOfPanic(c *match.Card) {

	c.Name = "Zero Nemesis, Shadow of Panic"
	c.Civs = []string{civ.Darkness}
	c.Power = 6000
	c.Family = []string{family.Ghost}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Evolution, fx.Doublebreaker, fx.WheneverOneOfMyCreaturesAttacksOppDiscardsRandom())

}

// GalekTheShadowWarrior ...
func GalekTheShadowWarrior(c *match.Card) {

	c.Name = "Galek, the Shadow Warrior"
	c.Power = 2000
	c.Civs = []string{civ.Darkness, civ.Fire}
	c.Family = []string{family.Ghost, family.Human}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Darkness, civ.Fire}

	c.Use(fx.Creature, fx.PutIntoManaZoneTapped, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		fx.DestroyOpBlocker(card, ctx)
		fx.OpponentDiscardsRandomCard(card, ctx)
	}))

}

// DolmarksTheShadowWarrior ...
func DolmarksTheShadowWarrior(c *match.Card) {

	c.Name = "Dolmarks, the Shadow Warrior"
	c.Power = 4000
	c.Civs = []string{civ.Darkness, civ.Fire}
	c.Family = []string{family.Ghost, family.Human}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness, civ.Fire}

	// Each half is mandatory and each player chooses from their own zones, in the
	// printed order: the controller pays first, then the opponent.
	c.Use(fx.Creature, fx.PutIntoManaZoneTapped, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		fx.OwnChoosesAndDestroysCreature(card, ctx)
		fx.OwnChoosesManaBurn(card, ctx)
		fx.OpponentChoosesAndDestroysCreature(card, ctx)
		fx.OpponentChoosesManaBurn(card, ctx)
	}))

}
