package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// ExplosiveTrooperZalmez ...
func ExplosiveTrooperZalmez(c *match.Card) {

	c.Name = "Explosive Trooper Zalmez"
	c.Power = 2000
	c.Civ = civ.Fire
	c.Family = []string{family.Armorloid}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		oppShields, err := ctx.Match.Opponent(card.Player).Container(match.SHIELDZONE)

		if err == nil && len(oppShields) <= 2 {
			fx.SelectFilter(
				card.Player,
				ctx.Match,
				ctx.Match.Opponent(card.Player),
				match.BATTLEZONE,
				fmt.Sprintf("%s's effect: You may destroy one of your opponent's creatures that has power 3000 or less.", card.Name),
				1,
				1,
				true,
				func(x *match.Card) bool {
					return ctx.Match.GetPower(x, false) <= 3000
				},
				false,
			).Map(func(x *match.Card) {
				ctx.Match.Destroy(x, card, match.DestroyedByMiscAbility)
			})
		}
	}))

}

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
