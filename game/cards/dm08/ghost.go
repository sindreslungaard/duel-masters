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
		fx.WhenAll([]func(*match.Card, *match.Context) bool{
			fx.AnotherOwnCreatureSummoned,
			func(card *match.Card, ctx *match.Context) bool {
				event, ok := ctx.Event.(*match.CardMoved)
				if !ok {
					return false
				}

				summonedCard, _ := card.Player.GetCard(event.CardID, match.BATTLEZONE)

				return summonedCard != nil && summonedCard.SharesAFamily(append(family.Dragons, family.Dragonoid))
			},
		},
			func(card *match.Card, ctx *match.Context) {
				creatures := make([]*match.Card, 0)

				creatures = append(creatures, fx.Find(card.Player, match.BATTLEZONE)...)

				creatures = append(creatures, fx.Find(ctx.Match.Opponent(card.Player), match.BATTLEZONE)...)

				if len(creatures) > 0 {
					slices.SortFunc(creatures, func(a *match.Card, b *match.Card) int {
						if ctx.Match.GetPower(a, false) < ctx.Match.GetPower(b, false) {
							return -1
						} else if ctx.Match.GetPower(a, false) > ctx.Match.GetPower(b, false) {
							return 1
						}
						return 0
					})

					minPower := ctx.Match.GetPower(creatures[0], false)
					cardsWithMinPower := make([]*match.Card, 0)

					for _, curr := range creatures {
						if ctx.Match.GetPower(curr, false) == minPower {
							cardsWithMinPower = append(cardsWithMinPower, curr)
						}
					}

					if len(cardsWithMinPower) == 1 {
						ctx.Match.Destroy(cardsWithMinPower[0], card, match.DestroyedByMiscAbility)
					} else if len(cardsWithMinPower) > 1 {
						// you choose among the tied creatures
						creaturesToDestroy := fx.SelectFromCollection(
							card.Player,
							ctx.Match,
							cardsWithMinPower,
							match.BATTLEZONE,
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
