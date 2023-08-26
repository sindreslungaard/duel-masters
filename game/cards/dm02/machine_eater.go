package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// EngineerKipo ...
func EngineerKipo(c *match.Card) {

	c.Name = "Engineer Kipo"
	c.Power = 2000
	c.Civ = civ.Fire
	c.Family = []string{family.MachineEater}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.CreatureDestroyed); ok {

			if event.Card == card {

				players := make([]*match.Player, 0)
				players = append(players, card.Player)
				players = append(players, ctx.Match.Opponent(card.Player))

				for _, p := range players {

					manazone, err := p.Container(match.MANAZONE)

					if err != nil {
						continue
					}

					toSelect := 1

					if len(manazone) < toSelect {
						toSelect = len(manazone)
					}

					if toSelect < 1 {
						continue
					}

					ctx.Match.Wait(ctx.Match.Opponent(p), "Waiting for your opponent to make an action")
					ctx.Match.NewAction(p, manazone, toSelect, toSelect, fmt.Sprintf("Engineer Kipo: Select %v card(s) from your manazone that will be sent to your graveyard", toSelect), false)

					for {

						action := <-p.Action

						if len(action.Cards) != toSelect || !match.AssertCardsIn(manazone, action.Cards...) {
							ctx.Match.DefaultActionWarning(p)
							continue
						}

						for _, id := range action.Cards {

							p.MoveCard(id, match.MANAZONE, match.GRAVEYARD)

							ctx.Match.Chat("Server", fmt.Sprintf("Engineer Kipo destroyed %v of %s's mana", toSelect, p.Username()))

						}

						break

					}

					ctx.Match.EndWait(ctx.Match.Opponent(p))
					ctx.Match.CloseAction(p)

				}

			}

		}

	})

}
