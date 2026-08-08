package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// TerradragonCusdalf ...
func TerradragonCusdalf(c *match.Card) {
	c.Name = "Terradragon Cusdalf"
	c.Power = 7000
	c.Civ = civ.Nature
	c.Family = []string{family.EarthDragon}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Nature}

	c.Use(
		fx.Creature,
		fx.Doublebreaker,
		fx.PowerAttacker4000,
		func(card *match.Card, ctx *match.Context) {
			if event, ok := ctx.Event.(*match.UntapManaEvent); ok &&
				card.Zone == match.BATTLEZONE &&
				event.Player == card.Player {
				ctx.InterruptFlow()
			}
		},
	)
}

// TerradragonDakmaBalgarow ...
func TerradragonDakmaBalgarow(c *match.Card) {
	c.Name = "Terradragon Dakma Balgarow"
	c.Power = 1000
	c.Civ = civ.Nature
	c.Family = []string{family.EarthDragon}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Nature}

	c.PowerModifier = func(m *match.Match, attacking bool) int {
		if c.Zone != match.BATTLEZONE {
			return 0
		}

		shields := len(fx.Find(c.Player, match.SHIELDZONE)) + len(fx.Find(m.Opponent(c.Player), match.SHIELDZONE))

		return shields * 2000
	}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {
		// GetPower dispatches handlers synchronously. Do not recursively ask for
		// power while already handling that calculation.
		if _, calculatingPower := ctx.Event.(*match.GetPowerEvent); calculatingPower {
			return
		}

		hasDoubleBreaker := terradragonDakmaBalgarowHasOwnCondition(card, cnd.DoubleBreaker)
		hasTripleBreaker := terradragonDakmaBalgarowHasOwnCondition(card, cnd.TripleBreaker)

		if card.Zone != match.BATTLEZONE {
			if hasDoubleBreaker {
				card.RemoveSpecificConditionBySource(cnd.DoubleBreaker, card.ID)
			}
			if hasTripleBreaker {
				card.RemoveSpecificConditionBySource(cnd.TripleBreaker, card.ID)
			}
			return
		}

		attacking := false
		switch event := ctx.Event.(type) {
		case *match.AttackPlayer:
			attacking = event.CardID == card.ID
		case *match.AttackCreature:
			attacking = event.CardID == card.ID
		case *match.SelectShields:
			attacking = event.Attacker == card
		}

		power := ctx.Match.GetPower(card, attacking)

		// "Triple breaker" replaces "double breaker" rather than stacking with it.
		wantsTripleBreaker := power >= 15000
		wantsDoubleBreaker := power >= 6000 && !wantsTripleBreaker

		if wantsDoubleBreaker && !hasDoubleBreaker {
			card.AddUniqueSourceCondition(cnd.DoubleBreaker, true, card.ID)
		} else if !wantsDoubleBreaker && hasDoubleBreaker {
			card.RemoveSpecificConditionBySource(cnd.DoubleBreaker, card.ID)
		}

		if wantsTripleBreaker && !hasTripleBreaker {
			card.AddUniqueSourceCondition(cnd.TripleBreaker, true, card.ID)
		} else if !wantsTripleBreaker && hasTripleBreaker {
			card.RemoveSpecificConditionBySource(cnd.TripleBreaker, card.ID)
		}
	})
}

func terradragonDakmaBalgarowHasOwnCondition(card *match.Card, id string) bool {
	for _, condition := range card.Conditions() {
		if condition.ID == id && condition.Src == card.ID {
			return true
		}
	}
	return false
}
