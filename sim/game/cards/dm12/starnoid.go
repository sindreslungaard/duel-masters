package dm12

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// WiseStarnoidAvatarOfHope ...
func WiseStarnoidAvatarOfHope(c *match.Card) {

	c.Name = "Wise Starnoid, Avatar of Hope"
	c.Power = 9000
	c.Civs = []string{civ.Light, civ.Water}
	c.Family = []string{family.Starnoid}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light, civ.Water}

	c.Use(
		fx.Creature,
		fx.PutIntoManaZoneTapped,
		fx.Doublebreaker,
		fx.VortexEvolution(
			family.LightBringer, func(x *match.Card) bool { return x.HasFamily(family.LightBringer) },
			family.CyberLord, func(x *match.Card) bool { return x.HasFamily(family.CyberLord) },
		),
		fx.When(fx.Attacking, fx.TopCardToShield),
		fx.When(fx.LeftBattlezone, fx.TopCardToShield),
	)

}
