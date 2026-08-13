package dm12

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// ValkyerStarstormElemental ...
func ValkyerStarstormElemental(c *match.Card) {

	c.Name = "Valkyer, Starstorm Elemental"
	c.Power = 7000
	c.Civs = []string{civ.Light}
	c.Family = []string{family.AngelCommand}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Blocker(), fx.CantAttackPlayers)

}

// KilstineNebulaElemental ...
func KilstineNebulaElemental(c *match.Card) {

	c.Name = "Kilstine, Nebula Elemental"
	c.Power = 5000
	c.Civs = []string{civ.Light}
	c.Family = []string{family.AngelCommand}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.WaveStriker, func(card *match.Card, ctx *match.Context) {
		active := fx.WaveStrikerActive(card, ctx.Match)

		// Swept across every zone when inactive, so a creature that left the
		// battle zone cannot carry a stale bonus back into it.
		fx.FindMultiple(card.Player, []string{match.BATTLEZONE, match.HAND, match.GRAVEYARD, match.MANAZONE, match.SHIELDZONE, match.HIDDENZONE}).Map(func(x *match.Card) {
			if !active || x.Zone != match.BATTLEZONE || x.ID == card.ID {
				x.RemoveSpecificConditionBySource(cnd.PowerAmplifier, card.ID)
				x.RemoveSpecificConditionBySource(cnd.DoubleBreaker, card.ID)
				x.RemoveSpecificConditionBySource(cnd.Blocker, card.ID)
				return
			}

			x.AddUniqueSourceCondition(cnd.PowerAmplifier, 5000, card.ID)
			x.AddUniqueSourceCondition(cnd.DoubleBreaker, true, card.ID)

			// ForceBlocker both marks the creature and offers it for the attack
			// being resolved, which a granted blocker needs in order to block.
			fx.ForceBlocker(x, ctx, card.ID)
		})
	})

}
