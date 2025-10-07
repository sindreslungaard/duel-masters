package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// BenzoTheHiddenFury ...
func BenzoTheHiddenFury(c *match.Card) {

	c.Name = "Benzo, the Hidden Fury"
	c.Power = 2000
	c.Civ = civ.Darkness
	c.Family = []string{family.PandorasBox}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.SHIELDZONE,
			fmt.Sprintf("%s's effect: Choose one of your shields and put it into your hand. You can use the 'Shield Trigger' ability of that shield.", card.Name),
			1,
			1,
			false,
		).Map(func(x *match.Card) {
			_, err := x.Player.MoveCard(x.ID, match.SHIELDZONE, match.HAND, card.ID)

			if err == nil {
				ctx.Match.ReportActionInChat(x.Player, fmt.Sprintf("%s was put into %s's hand from his shields.", x.Name, card.Player.Username()))

				fx.SelectFilter(
					card.Player,
					ctx.Match,
					card.Player,
					match.HAND,
					fmt.Sprintf("%s's effect: You can use the 'Shield Trigger' ability of %s by choosing it.", card.Name, x.Name),
					1,
					1,
					true,
					func(y *match.Card) bool {
						return y.ID == x.ID && y.HasCondition(cnd.ShieldTrigger)
					},
					false,
				).Map(func(y *match.Card) {
					if y.HasCondition(cnd.Spell) {
						ctx.Match.CastSpell(y, true)
					} else {
						ctx.Match.MoveCard(y, match.BATTLEZONE, card)
					}

					ctx.Match.HandleFx(match.NewContext(ctx.Match, &match.ShieldTriggerPlayedEvent{
						Card:   y,
						Source: card.ID,
					}))
				})
			}
		})
	}))

}

// DedreenTheHiddenCorrupter ...
func DedreenTheHiddenCorrupter(c *match.Card) {

	c.Name = "Dedreen, the Hidden Corrupter"
	c.Power = 4000
	c.Civ = civ.Darkness
	c.Family = []string{family.PandorasBox}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		oppShields, _ := ctx.Match.Opponent(card.Player).Container(match.SHIELDZONE)

		if oppShields != nil && len(oppShields) <= 3 {
			fx.OpponentDiscardsRandomCard(card, ctx)
		}
	}))

}
