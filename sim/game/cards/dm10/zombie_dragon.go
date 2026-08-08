package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// NecrodragonBryzenaga ...
func NecrodragonBryzenaga(c *match.Card) {

	c.Name = "Necrodragon Bryzenaga"
	c.Power = 9000
	c.Civ = civ.Darkness
	c.Family = []string{family.ZombieDragon}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Doublebreaker, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		// Snapshot the shield zone before moving anything out of it.
		shields := fx.Find(card.Player, match.SHIELDZONE)
		if len(shields) < 1 {
			return
		}

		// These shields are put into the hand, not broken, so BreakShields must
		// not be used: it would dispatch break events and hand turbo rush and
		// other "whenever a shield is broken" effects a false trigger. The shield
		// triggers are still offered, because this card explicitly allows it.
		shieldTriggers := make([]*match.Card, 0, len(shields))
		for _, shield := range shields {
			moved, err := card.Player.MoveCard(shield.ID, match.SHIELDZONE, match.HAND, card.ID)
			if err != nil || moved.Zone != match.HAND {
				continue
			}

			if moved.HasCondition(cnd.ShieldTrigger) {
				shieldTriggers = append(shieldTriggers, moved)
			}
		}

		ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s put all of %s's shields into their hand", card.Name, card.Player.Username()))

		ctx.Match.ResolveShieldTriggers(shieldTriggers, card)

	}))

}
