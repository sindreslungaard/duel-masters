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

	c.Use(fx.Creature, fx.When(fx.AttackConfirmed, func(card *match.Card, ctx *match.Context) {
		fx.SelectFilter(
			card.Player,
			ctx.Match,
			card.Player,
			match.GRAVEYARD,
			fmt.Sprintf("%s: You may return 1 light spell from your graveyard to your hand", card.Name),
			1,
			1,
			true,
			func(x *match.Card) bool { return x.HasCondition(cnd.Spell) && x.Civ == civ.Light },
			false,
		).Map(func(x *match.Card) {
			x.Player.MoveCard(x.ID, match.GRAVEYARD, match.HAND, card.ID)
			ctx.Match.ReportActionInChat(x.Player, fmt.Sprintf("%s retrieved %s from the graveyard to their hand", x.Player.Username(), x.Name))
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
