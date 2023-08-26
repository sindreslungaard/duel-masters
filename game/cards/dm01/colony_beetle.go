package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// DomeShell ...
func DomeShell(c *match.Card) {

	c.Name = "Dome Shell"
	c.Power = 3000
	c.Civ = civ.Nature
	c.Family = []string{family.ColonyBeetle}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.PowerAttacker2000)

}

// StormShell ...
func StormShell(c *match.Card) {

	c.Name = "Storm Shell"
	c.Power = 2000
	c.Civ = civ.Nature
	c.Family = []string{family.ColonyBeetle}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.CardMoved); ok {

			if event.CardID != card.ID || event.To != match.BATTLEZONE {
				return
			}

			opponent := ctx.Match.Opponent(card.Player)

			battlezone, err := opponent.Container(match.BATTLEZONE)

			if err != nil {
				return
			}

			if len(battlezone) < 1 {
				return
			}

			ctx.Match.Wait(card.Player, "Waiting for your opponent to make an action")

			ctx.Match.NewAction(opponent, battlezone, 1, 1, "Storm Shell: Select 1 card from your battlezone that will be sent to your manazone", false)

			defer func() {
				ctx.Match.EndWait(card.Player)
				ctx.Match.CloseAction(opponent)
			}()

			for {

				action := <-opponent.Action

				if len(action.Cards) != 1 || !match.AssertCardsIn(battlezone, action.Cards...) {
					ctx.Match.ActionWarning(opponent, "Your selection of cards does not fulfill the requirements")
					continue
				}

				movedCard, err := opponent.MoveCard(action.Cards[0], match.BATTLEZONE, match.MANAZONE)

				if err != nil {
					break
				}

				ctx.Match.Chat("Server", fmt.Sprintf("%s was moved from %s's battlezone to their manazone", movedCard.Name, opponent.Username()))

				break

			}

		}

	})

}

// TowerShell ...
func TowerShell(c *match.Card) {

	c.Name = "Tower Shell"
	c.Power = 5000
	c.Civ = civ.Nature
	c.Family = []string{family.ColonyBeetle}
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
					if ctx.Match.GetPower(blocker, false) >= 4000 {
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
					if ctx.Match.GetPower(blocker, false) >= 4000 {
						blockers = append(blockers, blocker)
					}
				}

				event.Blockers = blockers

			})

		}

	}, fx.Creature)

}
