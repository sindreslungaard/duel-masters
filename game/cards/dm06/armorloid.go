package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func ValiantWarriorExorious(c *match.Card) {

	c.Name = "Valiant Warrior Exorious"
	c.Power = 4000
	c.Civ = civ.Fire
	c.Family = []string{family.Armorloid}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.AttackUntapped, fx.PowerAttacker3000)

}

func AutomatedWeaponmasterMachai(c *match.Card) {

	c.Name = "Automated Weaponmaster Machai"
	c.Power = 4000
	c.Civ = civ.Fire
	c.Family = []string{family.Armorloid}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.ForceAttack)

}

func ArmoredScoutGestuchar(c *match.Card) {

	c.Name = "Armored Scout Gestuchar"
	c.Power = 4000
	c.Civ = civ.Fire
	c.Family = []string{family.Armorloid}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Fire}

	c.PowerModifier = func(m *match.Match, attacking bool) int {

		if (len(fx.FindFilter(c.Player, match.BATTLEZONE, func(card *match.Card) bool { return c.Civ == civ.Fire })) == 1) && attacking {
			return 3000
		}

		return 0

	}

	c.Use(fx.Creature, fx.When(fx.AttackingPlayer, func(card *match.Card, ctx *match.Context) {

		if len(fx.FindFilter(c.Player, match.BATTLEZONE, func(card *match.Card) bool { return c.Civ == civ.Fire })) == 1 {
			card.AddCondition(cnd.DoubleBreaker, nil, card.ID)
		}

	}))

}
