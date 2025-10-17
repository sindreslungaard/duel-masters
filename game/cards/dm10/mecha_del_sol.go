package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// RyudmilaChannelerOfSuns ...
func RyudmilaChannelerOfSuns(c *match.Card) {

	c.Name = "Ryudmila, Channeler of Suns"
	c.Power = 2000
	c.Civ = civ.Light
	c.Family = []string{family.MechaDelSol}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.ModifyPowers(func(event *match.GetPowerEvent) {
		if c.Zone == match.BATTLEZONE && event.Card == c {
			event.Power += len(fx.FindFilter(
				c.Player,
				match.BATTLEZONE,
				func(x *match.Card) bool {
					return !x.Tapped && x.ID != c.ID
				},
			)) * 2000
		}
	}), fx.When(fx.WouldBeDestroyed, fx.ShuffleCardIntoDeck))

}
