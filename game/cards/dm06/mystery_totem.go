package dm06

import (
	"duel-masters/game/civ"
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

		ctx.ScheduleAfter(func() {
			if event, ok := ctx.Event.(*match.AttackCreature); ok {

				if event.CardID != card.ID {
					return
				}

				cards := make([]*match.Card, 0)

				for _, blocker := range event.Blockers {

					if ctx.Match.GetPower(blocker, false) > 5000 {
						cards = append(cards, blocker)
					}
				}

				event.Blockers = cards
			}
		})

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

func ForbiddingTotem(c *match.Card) {

	c.Name = "Cannon Shell"
	c.Power = 1000
	c.Civ = civ.Nature
	c.Family = []string{family.MysteryTotem}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.ShieldTrigger)

}
