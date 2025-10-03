package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// NarielTheOracle ...
func NarielTheOracle(c *match.Card) {

	c.Name = "Nariel, the Oracle"
	c.Power = 1000
	c.Civ = civ.Light
	c.Family = []string{family.LightBringer}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.When(fx.InTheBattlezone, func(card *match.Card, ctx *match.Context) {
		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
			fx.Find(
				card.Player,
				match.BATTLEZONE,
			).Map(func(x *match.Card) {
				x.RemoveConditionBySource(card.ID)
			})

			fx.Find(
				ctx2.Match.Opponent(card.Player),
				match.BATTLEZONE,
			).Map(func(x *match.Card) {
				x.RemoveConditionBySource(card.ID)
			})

			if card.Zone != match.BATTLEZONE {
				exit()
				return
			}

			if _, ok := ctx2.Event.(*match.GetPowerEvent); ok {
				return
			}

			if event, ok := ctx2.Event.(*match.AttackPlayer); ok {
				creature, err := ctx2.Match.CurrentPlayer().Player.GetCard(event.CardID, match.BATTLEZONE)
				if err != nil {
					return
				}

				if !creature.HasCondition(cnd.IgnoreCantAttack) && ctx2.Match.GetPower(creature, false) >= 3000 {
					ctx2.Match.WarnPlayer(creature.Player, fmt.Sprintf("%s can't attack due to %s's effect.", creature.Name, card.Name))
					ctx2.InterruptFlow()
				}
			}

			if event, ok := ctx2.Event.(*match.AttackCreature); ok {
				creature, err := ctx2.Match.CurrentPlayer().Player.GetCard(event.CardID, match.BATTLEZONE)
				if err != nil {
					return
				}

				if ctx2.Match.GetPower(creature, false) >= 3000 {
					ctx2.Match.WarnPlayer(creature.Player, fmt.Sprintf("%s can't attack due to %s's effect.", creature.Name, card.Name))
					ctx2.InterruptFlow()
				}
			}

			if event, ok := ctx2.Event.(*match.TapAbility); ok {
				creature, err := ctx2.Match.CurrentPlayer().Player.GetCard(event.CardID, match.BATTLEZONE)
				if err != nil {
					return
				}

				if !creature.HasCondition(cnd.IgnoreCantAttack) && ctx2.Match.GetPower(creature, false) >= 3000 {
					ctx2.Match.WarnPlayer(creature.Player, fmt.Sprintf("%s can't use tap ability due to %s's effect.", creature.Name, card.Name))
					ctx2.InterruptFlow()
				}
			}

			fx.FindFilter(
				card.Player,
				match.BATTLEZONE,
				func(x *match.Card) bool {
					return ctx2.Match.GetPower(x, false) >= 3000
				},
			).Map(func(x *match.Card) {
				x.AddUniqueSourceCondition(cnd.CantAttackCreatures, nil, card.ID)
				if !x.HasCondition(cnd.IgnoreCantAttack) {
					x.AddUniqueSourceCondition(cnd.CantAttackPlayers, nil, card.ID)
				}
			})

			fx.FindFilter(
				ctx2.Match.Opponent(card.Player),
				match.BATTLEZONE,
				func(x *match.Card) bool {
					return ctx2.Match.GetPower(x, false) >= 3000
				},
			).Map(func(x *match.Card) {
				x.AddUniqueSourceCondition(cnd.CantAttackCreatures, nil, card.ID)
				if !x.HasCondition(cnd.IgnoreCantAttack) {
					x.AddUniqueSourceCondition(cnd.CantAttackPlayers, nil, card.ID)
				}
			})
		})
	}))

}
