package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// ArmoredRaiderGandaval ...
func ArmoredRaiderGandaval(c *match.Card) {

	c.Name = "Armored Raider Gandaval"
	c.Power = 6000
	c.Civs = []string{civ.Fire}
	c.Family = []string{family.Human}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Fire}

	c.PowerModifier = func(m *match.Match, attacking bool) int {
		if attacking {
			return 2000 * len(fx.FindFilter(
				c.Player,
				match.BATTLEZONE,
				func(x *match.Card) bool {
					return x.ID != c.ID && x.Tapped
				},
			))
		}

		return 0
	}

	c.Use(fx.Creature, fx.Evolution, fx.Doublebreaker)

}

// MezgerCommandoLeader ...
func MezgerCommandoLeader(c *match.Card) {

	c.Name = "Mezger Commando Leader"
	c.Power = 2000
	c.Civs = []string{civ.Fire}
	c.Family = []string{family.Human}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.SpeedAttacker)

}

// GontaTheWarriorSavage ...
func GontaTheWarriorSavage(c *match.Card) {

	c.Name = "Gonta, the Warrior Savage"
	c.Power = 4000
	c.Civs = []string{civ.Fire, civ.Nature}
	c.Family = []string{family.Human, family.BeastFolk}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Fire, civ.Nature}

	// Being put into the mana zone tapped is a rule of every multicolored card
	// and is handled by the engine, so there is nothing else to implement.
	c.Use(fx.Creature)

}
