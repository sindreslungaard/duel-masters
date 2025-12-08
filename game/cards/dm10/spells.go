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

// StaticWarp ...
func StaticWarp(c *match.Card) {

	c.Name = "Static Warp"
	c.Civ = civ.Light
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Light}

	textFmt := "%s's effect: Choose a creature in your battle zone. Tap the rest of the creatures in the battle zone."

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		players := []*match.Player{card.Player, ctx.Match.Opponent(card.Player)}

		for _, p := range players {
			fx.Select(
				p,
				ctx.Match,
				p,
				match.BATTLEZONE,
				fmt.Sprintf(textFmt, card.Name),
				1,
				1,
				false,
			).Map(func(x *match.Card) {
				fx.FindFilter(
					p,
					match.BATTLEZONE,
					func(y *match.Card) bool {
						return y.ID != x.ID
					},
				).Map(func(y *match.Card) {
					y.Tapped = true
				})
			})
		}

		ctx.Match.BroadcastState()
	}))

}
