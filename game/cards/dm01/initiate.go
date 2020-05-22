package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// FreiVizierOfAir ...
func FreiVizierOfAir(c *match.Card) {

	c.Name = "Frei, Vizier of Air"
	c.Power = 3000
	c.Civ = civ.Light
	c.Family = family.Initiate
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Untap)

}

// IereVizierOfBullets ...
func IereVizierOfBullets(c *match.Card) {

	c.Name = "Iere, Vizier of Bullets"
	c.Power = 3000
	c.Civ = civ.Light
	c.Family = family.Initiate
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature)

}

// LokVizierOfHunting ...
func LokVizierOfHunting(c *match.Card) {

	c.Name = "Lok, Vizier of Hunting"
	c.Power = 4000
	c.Civ = civ.Light
	c.Family = family.Initiate
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature)

}

// MieleVizierOfLightning ...
func MieleVizierOfLightning(c *match.Card) {

	c.Name = "Miele, Vizier of Lightning"
	c.Power = 1000
	c.Civ = civ.Light
	c.Family = family.Initiate
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.CardMoved); ok {

			if event.CardID == card.ID && event.To == match.BATTLEZONE {

				creatures := match.Search(card.Player, ctx.Match, ctx.Match.Opponent(card.Player), match.BATTLEZONE, "Miele, Vizier of Lightning: Select 1 of your opponent's creature and tap it. Close to not tap any creatures.", 1, 1, true)

				for _, creature := range creatures {
					creature.Tapped = true
				}

			}

		}

	})

}
