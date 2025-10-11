package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// SiegeRollerBagash ...
func SiegeRollerBagash(c *match.Card) {

	c.Name = "Siege Roller Bagash"
	c.Power = 3000
	c.Civ = civ.Fire
	c.Family = []string{family.Armorloid}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Fire}

	c.PowerModifier = func(m *match.Match, attacking bool) int {
		if attacking {
			return len(fx.FindFilter(
				c.Player,
				match.BATTLEZONE,
				func(x *match.Card) bool {
					return x.ID != c.ID && x.Tapped
				},
			)) * 1000
		}

		return 0
	}

	c.Use(fx.Creature)

}

// SmashWarriorStagrandu ...
func SmashWarriorStagrandu(c *match.Card) {

	c.Name = "Smash Warrior Stagrandu"
	c.Power = 1000
	c.Civ = civ.Fire
	c.Family = []string{family.Armorloid}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.AttackUntapped, func(card *match.Card, ctx *match.Context) {
		if event, ok := ctx.Event.(*match.SelectBlockers); ok {
			if event.Attacker.ID == card.ID && event.AttackedCardID != "" {
				attackedCard, _ := ctx.Match.Opponent(card.Player).GetCard(event.AttackedCardID, match.BATTLEZONE)

				if attackedCard != nil && ctx.Match.GetPower(attackedCard, false) >= 6000 {
					// This creature gets +9000 Power until the end of the turn
					ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
						if card.Zone != match.BATTLEZONE {
							card.RemoveConditionBySource(card.ID)
							exit()
							return
						}

						if _, ok := ctx2.Event.(*match.EndOfTurnStep); ok {
							card.RemoveConditionBySource(card.ID)
							exit()
							return
						}

						card.AddUniqueSourceCondition(cnd.PowerAmplifier, 9000, card.ID)
					})
				}
			}
		}
	})

}
