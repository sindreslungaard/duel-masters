package dm11

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// NialVizierOfDexterity ...
func NialVizierOfDexterity(c *match.Card) {

	c.Name = "Nial, Vizier of Dexterity"
	c.Power = 2500
	c.Civs = []string{civ.Light}
	c.Family = []string{family.Initiate}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.When(fx.EndOfMyTurnCreatureBZ, fx.MayUntapSelf))

}

// AsraVizierOfSafety ...
func AsraVizierOfSafety(c *match.Card) {

	c.Name = "Asra, Vizier of Safety"
	c.Power = 2000
	c.Civs = []string{civ.Light}
	c.Family = []string{family.Initiate}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light}

	c.PowerModifier = fx.WaveStrikerPower(c, 4000)

	c.Use(fx.Creature, fx.WaveStriker, func(card *match.Card, ctx *match.Context) {
		if !fx.WaveStrikerActive(card, ctx.Match) {
			card.RemoveSpecificConditionBySource(cnd.Blocker, card.ID)
			return
		}

		// ForceBlocker both marks the creature and offers it as a blocker for
		// the attack being resolved, which is what a granted blocker needs.
		fx.ForceBlocker(card, ctx, card.ID)
	})

}
