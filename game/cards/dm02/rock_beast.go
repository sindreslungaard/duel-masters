package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// Galsaur ...
func Galsaur(c *match.Card) {

	c.Name = "Galsaur"
	c.Power = 4000
	c.Civ = civ.Fire
	c.Family = []string{family.RockBeast}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Fire}

	c.PowerModifier = func(m *match.Match, attacking bool) int {

		if (len(fx.Find(c.Player, match.BATTLEZONE)) == 1) && attacking {
			return 4000
		}

		return 0

	}

	c.Use(fx.Creature, fx.When(fx.AttackingPlayer, func(card *match.Card, ctx *match.Context) {

		if len(fx.Find(card.Player, match.BATTLEZONE)) == 1 {
			card.AddCondition(cnd.DoubleBreaker, nil, card.ID)
		}

	}))

}

// Bombersaur ...
func Bombersaur(c *match.Card) {

	c.Name = "Bombersaur"
	c.Power = 5000
	c.Civ = civ.Fire
	c.Family = []string{family.RockBeast}
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
