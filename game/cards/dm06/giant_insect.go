package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

func UltraMantisScourgeOfFate(c *match.Card) {

	c.Name = "Ultra Mantis, Scourge of Fate"
	c.Power = 9000
	c.Civ = civ.Nature
	c.Family = []string{family.GiantInsect}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Nature}

	c.Use(func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.AttackPlayer); ok {

			if event.CardID != card.ID {
				return
			}

			ctx.ScheduleAfter(func() {

				blockers := make([]*match.Card, 0)

				for _, blocker := range event.Blockers {
					if ctx.Match.GetPower(blocker, false) > 8000 {
						blockers = append(blockers, blocker)
					}
				}

				event.Blockers = blockers

			})

		}

		if event, ok := ctx.Event.(*match.AttackCreature); ok {

			if event.CardID != card.ID {
				return
			}

			ctx.ScheduleAfter(func() {

				blockers := make([]*match.Card, 0)

				for _, blocker := range event.Blockers {
					if ctx.Match.GetPower(blocker, false) > 8000 {
						blockers = append(blockers, blocker)
					}
				}

				event.Blockers = blockers

			})

		}

	}, fx.Creature, fx.Evolution, fx.Doublebreaker)

}
func SplinterclawWasp(c *match.Card) {

	c.Name = "Splinterclaw Wasp"
	c.Power = 4000
	c.Civ = civ.Nature
	c.Family = []string{family.GiantInsect}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.Doublebreaker, fx.PowerAttacker3000, func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.Battle); ok {
			if !event.Blocked || event.Attacker != card {
				return
			}

			opponent := ctx.Match.Opponent(card.Player)

			ctx.Match.BreakShields(fx.SelectBackside(
				card.Player,
				ctx.Match,
				opponent,
				match.SHIELDZONE,
				"Splinterclaw Wasp: select shield to break",
				1,
				1,
				false,
			), card.ID)

			ctx.Match.ReportActionInChat(ctx.Match.Opponent(card.Player), fmt.Sprintf("Splinterclaw Wasp broke one of %s's shield", opponent.Username()))

		}
	})
}

func TrenchScarab(c *match.Card) {

	c.Name = "Trench Scarab"
	c.Power = 4000
	c.Civ = civ.Nature
	c.Family = []string{family.GiantInsect}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.CantAttackPlayers, fx.PowerAttacker4000)
}
