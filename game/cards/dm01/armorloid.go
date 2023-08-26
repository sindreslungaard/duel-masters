package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// ArmoredWalkerUrherion ...
func ArmoredWalkerUrherion(c *match.Card) {

	c.Name = "Armored Walker Urherion"
	c.Power = 3000
	c.Civ = civ.Fire
	c.Family = []string{family.Armorloid}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature)

	c.PowerModifier = func(m *match.Match, attacking bool) int {

		power := 0

		if attacking && match.ContainerHas(c.Player, match.BATTLEZONE, func(x *match.Card) bool { return x.HasFamily(family.Human) }) {
			power += 2000
		}

		return power
	}

}

// RothusTheTraveler ...
func RothusTheTraveler(c *match.Card) {

	c.Name = "Rothus, the Traveler"
	c.Power = 4000
	c.Civ = civ.Fire
	c.Family = []string{family.Armorloid}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.CardMoved); ok {

			if event.CardID == card.ID && event.To == match.BATTLEZONE {

				creatures := match.Search(card.Player, ctx.Match, card.Player, match.BATTLEZONE, "Rothus, the Traveler: Select 1 creature from your battlezone that will be sent to your graveyard", 1, 1, false)

				for _, creature := range creatures {
					ctx.Match.Destroy(creature, card, match.DestroyedByMiscAbility)
				}

				ctx.Match.Wait(card.Player, "Waiting for your opponent to make an action")
				defer ctx.Match.EndWait(card.Player)

				opponentCreatures := match.Search(ctx.Match.Opponent(card.Player), ctx.Match, ctx.Match.Opponent(card.Player), match.BATTLEZONE, "Rothus, the Traveler: Select 1 creature from your battlezone that will be sent to your graveyard", 1, 1, false)

				for _, creature := range opponentCreatures {
					ctx.Match.Destroy(creature, card, match.DestroyedByMiscAbility)
				}

			}

		}

	})

}
