package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// RapidReincarnation ...
func RapidReincarnation(c *match.Card) {

	c.Name = "Rapid Reincarnation"
	c.Civ = civ.Light
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.BATTLEZONE,
			fmt.Sprintf("%s's effect: You may destroy one of your creatures.\r\nIf you do, choose a creature in your hand that costs the same as or less than the number of cards in your mana zone\r\nand put it into the battle zone.", card.Name),
			1,
			1,
			true,
		).Map(func(x *match.Card) {
			manaZone, _ := card.Player.Container(match.MANAZONE)

			if manaZone != nil {
				ctx.Match.Destroy(x, card, match.DestroyedBySpell)

				manaCount := len(manaZone)

				fx.SelectFilter(
					card.Player,
					ctx.Match,
					card.Player,
					match.HAND,
					fmt.Sprintf("%s's effect: Choose a creature in your hand that costs the same as or less than the number of cards in your mana zone\r\nand put it into the battle zone.", card.Name),
					1,
					1,
					false,
					func(x *match.Card) bool {
						return x.HasCondition(cnd.Creature) && x.ManaCost <= manaCount && fx.CanBeSummoned(card.Player, x)
					},
					false,
				).Map(func(x *match.Card) {
					fx.ForcePutCreatureIntoBZ(ctx, x, match.HAND, card)
				})
			}
		})
	}))

}
