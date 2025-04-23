package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// MishaChannelerOfSuns ...
func MishaChannelerOfSuns(c *match.Card) {

	c.Name = "Misha, Channeler of Suns"
	c.Power = 5000
	c.Civ = civ.Light
	c.Family = []string{family.MechaDelSol}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.CantBeAttackedIf(func(attacker *match.Card) bool {
		return attacker.SharesAFamily(family.Dragons)
	}))
}

// SashaChannelerOfSuns ...
func SashaChannelerOfSuns(c *match.Card) {

	c.Name = "Sasha, Channeler of Suns"
	c.Power = 9500
	c.Civ = civ.Light
	c.Family = []string{family.MechaDelSol}
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Doublebreaker, fx.DragonBlocker(), func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.Battle); ok {
			if event.Attacker == card && event.Defender.SharesAFamily(family.Dragons) {
				event.AttackerPower += 6000
			} else if event.Defender == card && event.Attacker.SharesAFamily(family.Dragons) {
				event.DefenderPower += 6000
			}
		}

	})

}
