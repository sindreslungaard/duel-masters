package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// TwitchHornTheAgressor ...
func TwitchHornTheAgressor(c *match.Card) {

	c.Name = "Twitch Horn, the Aggressor"
	c.Power = 2000
	c.Civ = civ.Nature
	c.Family = []string{family.HornedBeast}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Nature}

	c.PowerModifier = func(m *match.Match, attacking bool) int {
		if attacking {
			return len(fx.FindFilter(
				c.Player,
				match.MANAZONE,
				func(x *match.Card) bool {
					return x.Tapped
				},
			)) * 2000
		}

		return 0
	}

	c.Use(fx.Creature)

}
