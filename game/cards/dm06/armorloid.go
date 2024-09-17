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

	c.Use(fx.Creature, fx.When(fx.InTheBattlezone, func(card *match.Card, ctx *match.Context) {

		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
			// Remove the conditions and exit the persistent effect when the card is no longer in the battle zone
			if card.Zone != match.BATTLEZONE {
				card.RemoveConditionBySource(card.ID)
				exit()
				return
			}

			// Whenever a card is moved while the persistent effect is active, remove the conditions
			// and add them back only if there are no other fire creatures in the battle zone
			if _, ok := ctx2.Event.(*match.CardMoved); ok {
				card.RemoveConditionBySource(card.ID)

				fireCreatures := fx.FindFilter(
					c.Player,
					match.BATTLEZONE,
					func(card *match.Card) bool { return card.Civ == civ.Fire },
				)

				if len(fireCreatures) == 1 {
					card.AddUniqueSourceCondition(cnd.PowerAttacker, 3000, c.ID)
					card.AddUniqueSourceCondition(cnd.DoubleBreaker, true, c.ID)
				}
			}
		})
	}))

}
