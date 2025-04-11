package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
	"slices"
)

// ScreamSlicerShadowOfFear ...
func ScreamSlicerShadowOfFear(c *match.Card) {

	c.Name = "Scream Slicer, Shadow of Fear"
	c.Power = 4000
	c.Civ = civ.Darkness
	c.Family = []string{family.Ghost}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature,
		fx.When(
			fx.AnotherOwnDragonoidOrDragonSummoned,
			func(card *match.Card, ctx *match.Context) {
				allCreatures := make([]*match.Card, 0)

				myCreatures := fx.Find(card.Player, match.BATTLEZONE)

				oppCreatures := fx.Find(ctx.Match.Opponent(card.Player), match.BATTLEZONE)

				allCreatures = append(allCreatures, myCreatures...)
				allCreatures = append(allCreatures, oppCreatures...)

				if len(allCreatures) > 0 {
					slices.SortFunc(allCreatures, func(a *match.Card, b *match.Card) int {
						if ctx.Match.GetPower(a, false) < ctx.Match.GetPower(b, false) {
							return -1
						} else if ctx.Match.GetPower(a, false) > ctx.Match.GetPower(b, false) {
							return 1
						}
						return 0
					})

					minPower := ctx.Match.GetPower(allCreatures[0], false)
					myCardsWithMinPower := make([]*match.Card, 0)
					oppCardsWithMinPower := make([]*match.Card, 0)

					for _, curr := range myCreatures {
						if ctx.Match.GetPower(curr, false) == minPower {
							myCardsWithMinPower = append(myCardsWithMinPower, curr)
						}
					}

					for _, curr := range oppCreatures {
						if ctx.Match.GetPower(curr, false) == minPower {
							oppCardsWithMinPower = append(oppCardsWithMinPower, curr)
						}
					}

					if len(myCardsWithMinPower) == 1 && len(oppCardsWithMinPower) == 0 {
						ctx.Match.Destroy(myCardsWithMinPower[0], card, match.DestroyedByMiscAbility)
					} else if len(oppCardsWithMinPower) == 1 && len(myCardsWithMinPower) == 0 {
						ctx.Match.Destroy(oppCardsWithMinPower[0], card, match.DestroyedByMiscAbility)
					} else {

						cardsMap := make(map[string][]*match.Card)

						cardsMap["Your creatures"] = myCardsWithMinPower
						cardsMap["Opponent's creatures"] = oppCardsWithMinPower

						creaturesToDestroy := fx.SelectMultipart(
							card.Player,
							ctx.Match,
							cardsMap,
							fmt.Sprintf("%s's effect: Choose 1 of the following creatures to destroy.", card.Name),
							1,
							1,
							false,
						)

						if len(creaturesToDestroy) > 0 {
							ctx.Match.Destroy(creaturesToDestroy[0], card, match.DestroyedByMiscAbility)
						}
					}
				}
			}))
}
