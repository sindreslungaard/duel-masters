package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// Soulswap ...
func Soulswap(c *match.Card) {

	c.Name = "Soulswap"
	c.Civ = civ.Nature
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		cards := make(map[string][]*match.Card)

		myCards, err := card.Player.Container(match.BATTLEZONE)

		if err != nil {
			return
		}

		opponentCards, err := ctx.Match.Opponent(card.Player).Container(match.BATTLEZONE)

		if err != nil {
			return
		}

		if len(myCards) < 1 && len(opponentCards) < 1 {
			return
		}

		cards["Your creatures"] = myCards
		cards["Opponent's creatures"] = opponentCards

		fx.SelectMultipart(
			card.Player,
			ctx.Match,
			cards,
			fmt.Sprintf("%s's effect: You may choose a creature in the battle zone and put it into its owner's mana zone.\r\n If you do, choose a non-evolution creature in that player's mana zone that costs the same as or less than the number of cards in that mana zone.\r\n That player puts that creature into the battle zone.", card.Name),
			1,
			1,
			true,
		).Map(func(x *match.Card) {
			_, err := x.Player.MoveCard(x.ID, match.BATTLEZONE, match.MANAZONE, card.ID)

			if err == nil {
				ctx.Match.ReportActionInChat(x.Player, fmt.Sprintf("%s was put into %s's mana zone from his battle zone.", x.Name, x.Player.Username()))
			}

			manaZone, _ := x.Player.Container(match.MANAZONE)

			if manaZone != nil {
				manaZoneLen := len(manaZone)

				fx.SelectFilter(
					card.Player,
					ctx.Match,
					x.Player,
					match.MANAZONE,
					fmt.Sprintf("Choose a non-evolution creature in %s's mana zone that costs the same as or less than the number of cards in that mana zone.\r\n %s puts that creature into the battle zone.", x.Player.Username(), x.Player.Username()),
					1,
					1,
					false,
					func(x *match.Card) bool {
						return x.HasCondition(cnd.Creature) && !x.HasCondition(cnd.Evolution) && x.ManaCost <= manaZoneLen
					},
					false,
				).Map(func(x *match.Card) {
					fx.ForcePutCreatureIntoBZ(ctx, x, match.MANAZONE, card)
				})
			}
		})
	}))

}

// ThirstForTheHunt ...
func ThirstForTheHunt(c *match.Card) {

	c.Name = "Thirst for the Hunt"
	c.Civ = civ.Nature
	c.ManaCost = 1
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		fx.Find(
			card.Player,
			match.BATTLEZONE,
		).Map(func(x *match.Card) {
			x.AddUniqueSourceCondition(cnd.PowerAttacker, 1000, card.ID)
		})
	}))

}
