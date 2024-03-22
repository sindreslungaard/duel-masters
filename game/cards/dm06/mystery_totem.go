package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

func BlissTotemAvatarOfLuck(c *match.Card) {

	c.Name = "Bliss Totem, Avatar of Luck"
	c.Power = 5000
	c.Civ = civ.Nature
	c.Family = []string{family.MysteryTotem}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Nature}
	c.TapAbility = true

	c.Use(fx.Creature, fx.When(fx.TapAbility, func(card *match.Card, ctx *match.Context) {

		fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.GRAVEYARD,
			"Bliss Totem, Avatar of Luck: Select up to 3 cards from your graveyard and put it in your manazone",
			0,
			3,
			true,
		).Map(func(c *match.Card) {
			card.Player.MoveCard(c.ID, match.GRAVEYARD, match.MANAZONE)
			ctx.Match.Chat("Server", fmt.Sprintf("%s was moved to %s's manazone by %s", c.Name, c.Player.Username(), card.Name))

			card.Tapped = true
		})

	}))
}

func ClobberTotem(c *match.Card) {

	c.Name = "Clobber Totem"
	c.Power = 4000
	c.Civ = civ.Nature
	c.Family = []string{family.MysteryTotem}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Nature}

	c.Use(func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.AttackCreature); ok {

			if event.CardID != card.ID {
				return
			}

			ctx.ScheduleAfter(func() {
				cards := make([]*match.Card, 0)

				for _, blocker := range event.Blockers {

					if ctx.Match.GetPower(blocker, false) > 5000 {
						cards = append(cards, blocker)
					}
				}

				event.Blockers = cards
			})
		}

		if event, ok := ctx.Event.(*match.AttackPlayer); ok {

			if event.CardID != card.ID {
				return
			}

			ctx.ScheduleAfter(func() {
				cards := make([]*match.Card, 0)

				for _, blocker := range event.Blockers {

					if ctx.Match.GetPower(blocker, false) > 5000 {
						cards = append(cards, blocker)
					}
				}

				event.Blockers = cards
			})
		}

	}, fx.Creature, fx.PowerAttacker2000, fx.Doublebreaker)
}

// ForbiddingTotem
func ForbiddingTotem(c *match.Card) {

	c.Name = "Forbidding Totem"
	c.Power = 4000
	c.Civ = civ.Nature
	c.Family = []string{family.MysteryTotem}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		if c.Zone == match.BATTLEZONE && ctx.Match.IsPlayerTurn(ctx.Match.Opponent(c.Player)) {

			if _, ok := ctx.Event.(*match.AttackPlayer); ok {

				attackableCreatures := fx.FindFilter(
					c.Player,
					match.BATTLEZONE,
					func(c *match.Card) bool {
						return c.HasFamily(family.MysteryTotem) && (c.Tapped || card.HasCondition(cnd.AttackUntapped))
					},
				)

				if len(attackableCreatures) == 0 {
					return
				}

				ctx.Match.WarnPlayer(ctx.Match.Opponent(card.Player), "Creatures can only attack Mystery Totems")

				ctx.InterruptFlow()

			}

			if _, ok := ctx.Event.(*match.AttackCreature); ok {

				attackableCreatures := fx.FindFilter(
					c.Player,
					match.BATTLEZONE,
					func(c *match.Card) bool {
						return c.HasFamily(family.MysteryTotem) && (c.Tapped || card.HasCondition(cnd.AttackUntapped))
					},
				)

				if len(attackableCreatures) == 0 {
					return
				}

				ctx.Match.WarnPlayer(ctx.Match.Opponent(card.Player), "Creatures can only attack Mystery Totems")

				ctx.InterruptFlow()

			}

		}

	})
}
