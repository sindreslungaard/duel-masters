package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// Bombersaur ...
func Bombersaur(c *match.Card) {

	c.Name = "Bombersaur"
	c.Power = 5000
	c.Civ = civ.Fire
	c.Family = family.RockBeast
	c.ManaCost = 5
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

					toSelect := 2

					if len(manazone) < toSelect {
						toSelect = len(manazone)
					}

					if toSelect < 1 {
						continue
					}

					ctx.Match.Wait(ctx.Match.Opponent(p), "Waiting for your opponent to make an action")
					ctx.Match.NewAction(p, manazone, toSelect, toSelect, fmt.Sprintf("Bombersaur: Select %v card(s) from your manazone that will be sent to your graveyard", toSelect), false)

					for {

						action := <-p.Action

						if len(action.Cards) != toSelect || !match.AssertCardsIn(manazone, action.Cards...) {
							ctx.Match.DefaultActionWarning(p)
							continue
						}

						for _, id := range action.Cards {

							p.MoveCard(id, match.MANAZONE, match.GRAVEYARD)

							ctx.Match.Chat("Server", fmt.Sprintf("Bombersaur destroyed %v of %s's mana", toSelect, p.Username()))

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

// Meteosaur ...
func Meteosaur(c *match.Card) {

	c.Name = "Meteosaur"
	c.Power = 2000
	c.Civ = civ.Fire
	c.Family = family.RockBeast
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.CardMoved); ok {

			if event.CardID == card.ID && event.To == match.BATTLEZONE {

				creatures := match.Filter(
					card.Player,
					ctx.Match,
					ctx.Match.Opponent(card.Player),
					match.BATTLEZONE,
					"Meteosaur: Select 1 of your opponent's creatures with power 2000 or less and destroy it",
					1,
					1,
					true,
					func(x *match.Card) bool { return x.Power <= 2000 },
				)

				for _, creature := range creatures {
					ctx.Match.Destroy(creature, card)
				}

			}
		}

	})

}
