package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
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

	c.Use(
		fx.Creature,
		fx.ModifyPowers(func(event *match.GetPowerEvent) {
			if c.Zone != match.BATTLEZONE || event.Card != c {
				return
			}

			event.Power += len(fx.FindFilter(
				c.Player,
				match.BATTLEZONE,
				func(card *match.Card) bool {
					return card.ID != c.ID &&
						card.HasCondition(cnd.Creature) &&
						!card.Tapped
				},
			)) * 2000
		}),
		fx.When(fx.WouldBeDestroyed, fx.ShuffleIntoDeckInsteadOfDestruction),
	)
}

// BerochikaChannelerOfSuns ...
func BerochikaChannelerOfSuns(c *match.Card) {

	c.Name = "Berochika, Channeler of Suns"
	c.Power = 5000
	c.Civ = civ.Light
	c.Family = []string{family.MechaDelSol}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		myShields, _ := card.Player.Container(match.SHIELDZONE)

		if len(myShields) >= 5 {
			fx.TopCardToShield(card, ctx)
		}
	}))

}
