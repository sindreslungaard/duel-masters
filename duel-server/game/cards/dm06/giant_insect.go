package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func UltraMantisScourgeOfFate(c *match.Card) {

	c.Name = "Ultra Mantis, Scourge of Fate"
	c.Power = 9000
	c.Civ = civ.Nature
	c.Family = []string{family.GiantInsect}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.Evolution, fx.Doublebreaker, fx.CantBeBlockedByPowerUpTo8000)

}
func SplinterclawWasp(c *match.Card) {

	c.Name = "Splinterclaw Wasp"
	c.Power = 4000
	c.Civ = civ.Nature
	c.Family = []string{family.GiantInsect}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.Doublebreaker, fx.PowerAttacker3000, fx.When(fx.Blocked, fx.DestroyOpShield))
}

func TrenchScarab(c *match.Card) {

	c.Name = "Trench Scarab"
	c.Power = 4000
	c.Civ = civ.Nature
	c.Family = []string{family.GiantInsect}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.CantAttackPlayers, fx.PowerAttacker4000)
}
