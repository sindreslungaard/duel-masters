package promo

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// QTronicOmnistrain ...
func QTronicOmnistrain(c *match.Card) {

	c.Name = "Q-tronic Omnistrain"
	c.Power = 3000
	c.Civs = []string{civ.Nature}
	c.Family = []string{family.Survivor}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.ShieldTrigger, fx.Evolution,
		fx.When(fx.InTheBattlezone, fx.AuraForOwnCreatures(cnd.AddFamily, family.Survivor, func(source, candidate *match.Card) bool {
			return true
		})),
	)

}
