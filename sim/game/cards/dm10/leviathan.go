package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// KingOquanos ...
func KingOquanos(c *match.Card) {
	c.Name = "King Oquanos"
	c.Power = 2000
	c.Civ = civ.Water
	c.Family = []string{family.Leviathan}
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Water}

	c.PowerModifier = func(m *match.Match, attacking bool) int {
		if c.Zone != match.BATTLEZONE {
			return 0
		}

		return len(fx.FindFilter(
			m.Opponent(c.Player),
			match.MANAZONE,
			func(card *match.Card) bool { return card.Tapped },
		)) * 2000
	}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {
		// GetPower dispatches handlers synchronously. Do not recursively ask for
		// power while already handling that calculation.
		if _, calculatingPower := ctx.Event.(*match.GetPowerEvent); calculatingPower {
			return
		}

		if card.Zone != match.BATTLEZONE {
			if kingOquanosHasOwnDoubleBreaker(card) {
				card.RemoveSpecificConditionBySource(cnd.DoubleBreaker, card.ID)
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

		hasDoubleBreaker := kingOquanosHasOwnDoubleBreaker(card)
		power := ctx.Match.GetPower(card, attacking)
		if power >= 6000 && !hasDoubleBreaker {
			card.AddUniqueSourceCondition(cnd.DoubleBreaker, true, card.ID)
		} else if power < 6000 && hasDoubleBreaker {
			card.RemoveSpecificConditionBySource(cnd.DoubleBreaker, card.ID)
		}
	})
}

func kingOquanosHasOwnDoubleBreaker(card *match.Card) bool {
	for _, condition := range card.Conditions() {
		if condition.ID == cnd.DoubleBreaker && condition.Src == card.ID {
			return true
		}
	}
	return false
}
