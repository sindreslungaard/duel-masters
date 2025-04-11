package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// FreiVizierOfAir ...
func FreiVizierOfAir(c *match.Card) {

	c.Name = "Frei, Vizier of Air"
	c.Power = 3000
	c.Civ = civ.Light
	c.Family = []string{family.Initiate}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.When(fx.EndOfMyTurnCreatureBZ, fx.MayUntapSelf))

}

// IereVizierOfBullets ...
func IereVizierOfBullets(c *match.Card) {

	c.Name = "Iere, Vizier of Bullets"
	c.Power = 3000
	c.Civ = civ.Light
	c.Family = []string{family.Initiate}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature)

}

// LokVizierOfHunting ...
func LokVizierOfHunting(c *match.Card) {

	c.Name = "Lok, Vizier of Hunting"
	c.Power = 4000
	c.Civ = civ.Light
	c.Family = []string{family.Initiate}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature)

}

// MieleVizierOfLightning ...
func MieleVizierOfLightning(c *match.Card) {

	c.Name = "Miele, Vizier of Lightning"
	c.Power = 1000
	c.Civ = civ.Light
	c.Family = []string{family.Initiate}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.When(fx.Summoned, fx.MayTapOpCreature))

}

// ToelVizierOfHope ...
func ToelVizierOfHope(c *match.Card) {

	c.Name = "Toel, Vizier of Hope"
	c.Power = 2000
	c.Civ = civ.Light
	c.Family = []string{family.Initiate}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		if _, ok := ctx.Event.(*match.EndOfTurnStep); ok {

			if !card.Player.HasCard(match.BATTLEZONE, card.ID) ||
				!ctx.Match.IsPlayerTurn(card.Player) {
				return
			}

			if fx.BinaryQuestion(card.Player, ctx.Match, fmt.Sprintf("Do you want to untap all your creatures in the battlezone? (%s's effect)", card.Name)) {
				creatures, err := card.Player.Container(match.BATTLEZONE)

				if err != nil {
					return
				}

				madeChanges := false

				for _, creature := range creatures {

					if creature.Tapped {
						creature.Tapped = false
						madeChanges = true
					}
				}

				if madeChanges {
					ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s untapped all creatures in %s's battlezone", card.Name, card.Player.Username()))
				}
			}

		}

	})

}
