package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// UliyaTheEntrancer ...
func UliyaTheEntrancer(c *match.Card) {

	c.Name = "Uliya, the Entrancer"
	c.Power = 5000
	c.Civ = civ.Darkness
	c.Family = []string{family.DarkLord}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		fx.SelectBackside(
			card.Player,
			ctx.Match,
			card.Player,
			match.SHIELDZONE,
			fmt.Sprintf("%s's effect: Choose one of your shields and put it into your hand. You can use the 'shield trigger' ability of that shield.", card.Name),
			1,
			1,
			false,
		).Map(func(x *match.Card) {
			_, err := x.Player.MoveCard(x.ID, match.SHIELDZONE, match.HAND, card.ID)

			if err == nil {
				ctx.Match.ReportActionInChat(x.Player, fmt.Sprintf("%s was put into %s's hand from his shields.", x.Name, x.Player.Username()))

				fx.SelectFilter(
					card.Player,
					ctx.Match,
					card.Player,
					match.HAND,
					fmt.Sprintf("You can use the 'shield trigger' ability of %s if you choose it.", x.Name),
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
