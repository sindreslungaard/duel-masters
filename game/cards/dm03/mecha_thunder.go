package dm03

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// RaVuSeekerOfLightning ...
func RaVuSeekerOfLightning(c *match.Card) {

	c.Name = "Ra Vu, Seeker of Lightning"
	c.Power = 4000
	c.Civ = civ.Light
	c.Family = []string{family.MechaThunder}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.When(fx.Attacking, func(card *match.Card, ctx *match.Context) {

		ctx.ScheduleAfter(func() {
			spells := match.Filter(card.Player, ctx.Match, card.Player, match.GRAVEYARD, "You may select 1 spell from your graveyard that will be sent to your hand", 0, 1, false, func(x *match.Card) bool { return x.HasCondition(cnd.Spell) })

			for _, spell := range spells {

				card.Player.MoveCard(spell.ID, match.GRAVEYARD, match.HAND)
				ctx.Match.Chat("Server", fmt.Sprintf("%s retrieved %s from the graveyard to their hand", spell.Player.Username(), spell.Name))

			}
		})
	}))

}

// UrPaleSeekerOfSunlight ...
func UrPaleSeekerOfSunlight(c *match.Card) {

	c.Name = "Ur Pale, Seeker of Sunlight"
	c.Power = 2500
	c.Civ = civ.Light
	c.Family = []string{family.MechaThunder}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light}

	c.PowerModifier = func(m *match.Match, attacking bool) int {

		if match.ContainerHas(c.Player, match.MANAZONE, func(x *match.Card) bool { return x.Civ != civ.Light }) {
			return 0
		}

		return 2000
	}

	c.Use(fx.Creature)

}
