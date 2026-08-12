package dm11

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// dragonFromMana is the ability Rollicking Totem and Royal Durian share:
// "Put a creature that has Dragon in its race from your mana zone into the
// battle zone."
func dragonFromMana() match.HandlerFunc {
	return fx.PutCreatureFromManaIntoBZ(
		func(x *match.Card) bool {
			return x.HasCondition(cnd.Creature) && x.SharesAFamily(family.Dragons)
		},
		"a creature that has Dragon in its race",
	)
}

// RollickingTotem ...
func RollickingTotem(c *match.Card) {

	c.Name = "Rollicking Totem"
	c.Power = 4000
	c.Civs = []string{civ.Nature}
	c.Family = []string{family.MysteryTotem}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.SilentSkill(dragonFromMana()))

}
