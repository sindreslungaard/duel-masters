package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// RocketdiveSkyterror ...
func RocketdiveSkyterror(c *match.Card) {

	c.Name = "Rocketdive Skyterror"
	c.Power = 5000
	c.Civ = civ.Fire
	c.Family = []string{family.ArmoredWyvern}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.CantBeAttacked, fx.CantAttackPlayers, fx.PowerAttacker1000)
}

// TorpedoSkyterror ...
func TorpedoSkyterror(c *match.Card) {

	c.Name = "Torpedo Skyterror"
	c.Power = 4000
	c.Civ = civ.Fire
	c.Family = []string{family.ArmoredWyvern}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.When(fx.Attacking, func(card *match.Card, ctx *match.Context) {

		card.RemoveConditionBySource(card.ID + "-custom")

		otherTapped := fx.FindFilter(
			card.Player,
			match.BATTLEZONE,
			func(x *match.Card) bool {
				return x.Tapped && x.ID != card.ID
			},
		)

		power := len(otherTapped) * 2000

		card.AddUniqueSourceCondition(cnd.PowerAttacker, power, card.ID+"-custom")

	}))
}
