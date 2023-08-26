package dm05

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// DeathCruzerTheAnnihilator ...
func DeathCruzerTheAnnihilator(c *match.Card) {

	c.Name = "Death Cruzer, the Annihilator"
	c.Power = 13000
	c.Civ = civ.Darkness
	c.Family = []string{family.DemonCommand}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Triplebreaker, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		fx.Find(card.Player, match.BATTLEZONE).Map(func(x *match.Card) {

			if x.ID != card.ID {
				ctx.Match.Destroy(x, card, match.DestroyedByMiscAbility)
			}

		})

	}))

}

// VashunaSwordDancer ...
func VashunaSwordDancer(c *match.Card) {

	c.Name = "Vashuna, Sword Dancer"
	c.Power = 7000
	c.Civ = civ.Darkness
	c.Family = []string{family.DemonCommand}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Doublebreaker, func(card *match.Card, ctx *match.Context) {

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
