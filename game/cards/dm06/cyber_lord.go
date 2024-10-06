package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

func Sopian(c *match.Card) {

	c.Name = "Sopian"
	c.Power = 2000
	c.Civ = civ.Water
	c.Family = []string{family.CyberLord}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Water}
	c.TapAbility = fx.GiveOwnCreatureCantBeBlocked

	c.Use(fx.Creature, fx.TapAbility)
}

func Kyuroro(c *match.Card) {

	c.Name = "Kyuroro"
	c.Power = 2000
	c.Civ = civ.Water
	c.Family = []string{family.CyberLord}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.When(fx.BreakShield, func(card *match.Card, ctx *match.Context) {

		event, ok := ctx.Event.(*match.BreakShieldEvent)
		if !ok {
			return
		}

		if event.Source.Player == card.Player ||
			!event.Source.HasCondition(cnd.Creature) ||
			len(event.Cards) < 1 {
			return
		}

		ctx.Match.Wait(ctx.Match.Opponent(card.Player), "Waiting for your opponent to make an action")
		defer ctx.Match.EndWait(ctx.Match.Opponent(card.Player))

		nrOfShields := len(event.Cards)

		newShieldsSelection := fx.SelectBackside(
			card.Player,
			ctx.Match,
			event.Cards[0].Player,
			match.SHIELDZONE,
			fmt.Sprintf("Kyuroro: choose %d shield(s) that your opponent will break", nrOfShields),
			nrOfShields,
			nrOfShields,
			false,
		)

		event.Cards = newShieldsSelection
	}))
}
