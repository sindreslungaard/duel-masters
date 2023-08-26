package dm05

import (
	"fmt"

	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func Gigazoul(c *match.Card) {

	c.Name = "Gigazoul"
	c.Power = 3000
	c.Civ = civ.Darkness
	c.Family = []string{family.Chimera}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		if _, ok := ctx.Event.(*match.UntapStep); ok {

			opponentShields := fx.Find(ctx.Match.Opponent(card.Player), match.SHIELDZONE)
			n := len(opponentShields)

			if n == 0 {
				card.AddCondition(cnd.CantAttackPlayers, true, card.ID)
				card.AddCondition(cnd.CantAttackCreatures, true, card.ID)
			} else {
				card.RemoveCondition(cnd.CantAttackPlayers)
				card.RemoveCondition(cnd.CantAttackCreatures)
			}

		}

		if event, ok := ctx.Event.(*match.CardMoved); ok {

			if event.From == match.SHIELDZONE || event.To == match.SHIELDZONE {

				opponentShields := fx.Find(ctx.Match.Opponent(card.Player), match.SHIELDZONE)
				n := len(opponentShields)

				if n == 0 {
					card.AddCondition(cnd.CantAttackPlayers, true, card.ID)
					card.AddCondition(cnd.CantAttackCreatures, true, card.ID)
				} else {
					card.RemoveCondition(cnd.CantAttackPlayers)
					card.RemoveCondition(cnd.CantAttackCreatures)
				}

			}

		}

		if event, ok := ctx.Event.(*match.AttackPlayer); ok {

			// Is this event for me or someone else?
			if event.CardID != card.ID || !card.HasCondition(cnd.CantAttackPlayers) {
				return
			}

			ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("%s can't attack players", card.Name))

			ctx.InterruptFlow()

		}

		if event, ok := ctx.Event.(*match.AttackCreature); ok {

			// Is this event for me or someone else?
			if event.CardID != card.ID || !card.HasCondition(cnd.CantAttackCreatures) {
				return
			}

			ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("%s can't attack creatures", card.Name))

			ctx.InterruptFlow()

		}

	})

}

func Gigakail(c *match.Card) {
	c.Name = "Gigakail"
	c.Power = 4000
	c.Civ = civ.Darkness
	c.Family = []string{family.Chimera}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.ConditionalSlayer(func(target *match.Card) bool {
		return target.Civ == civ.Nature || target.Civ == civ.Light
	}))
}

func GigalingQ(c *match.Card) {

	c.Name = "Gigaling Q"
	c.Power = 2000
	c.Civ = civ.Darkness
	c.Family = []string{family.Chimera}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Survivor, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
			if card.Zone != match.BATTLEZONE {
				fx.FindFilter(
					card.Player,
					match.BATTLEZONE,
					func(x *match.Card) bool { return x.HasCondition(cnd.Survivor) },
				).Map(func(x *match.Card) {
					x.RemoveConditionBySource(card.ID)
				})

				exit()
				return
			}

			fx.FindFilter(
				card.Player,
				match.BATTLEZONE,
				func(x *match.Card) bool { return x.HasCondition(cnd.Survivor) },
			).Map(func(x *match.Card) {
				x.AddUniqueSourceCondition(cnd.Slayer, true, card.ID)
			})

		})

	}))

}
