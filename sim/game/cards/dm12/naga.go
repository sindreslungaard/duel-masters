package dm12

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// CruelNagaAvatarOfFate ...
func CruelNagaAvatarOfFate(c *match.Card) {

	c.Name = "Cruel Naga, Avatar of Fate"
	c.Power = 9000
	c.Civs = []string{civ.Water, civ.Darkness}
	c.Family = []string{family.Naga}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Water, civ.Darkness}

	c.Use(
		fx.Creature,
		fx.PutIntoManaZoneTapped,
		fx.Doublebreaker,
		fx.CantBeBlocked,
		fx.VortexEvolution(
			family.Merfolk, func(x *match.Card) bool { return x.HasFamily(family.Merfolk) },
			family.Chimera, func(x *match.Card) bool { return x.HasFamily(family.Chimera) },
		),
		fx.When(fx.LeftBattlezone, fx.DestroyAllCreatures(match.DestroyedByMiscAbility)),
	)

}
