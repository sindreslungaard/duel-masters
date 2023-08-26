package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// AstrocometDragon ...
func AstrocometDragon(c *match.Card) {

	c.Name = "Astrocomet Dragon"
	c.Power = 6000
	c.Civ = civ.Fire
	c.Family = []string{family.ArmoredDragon}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.PowerAttacker4000, fx.Doublebreaker)

}

// BolshackDragon ...
func BolshackDragon(c *match.Card) {

	c.Name = "Bolshack Dragon"
	c.Power = 6000
	c.Civ = civ.Fire
	c.Family = []string{family.ArmoredDragon}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.Doublebreaker, func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.AttackPlayer); ok {
			if event.CardID != card.ID {
				return
			}
			bolshackSpecial(card)
		}

		if event, ok := ctx.Event.(*match.AttackCreature); ok {
			if event.CardID != card.ID {
				return
			}
			bolshackSpecial(card)
		}

	})

}

// +1000 power attacker for each fire civilization card in the graveyard
func bolshackSpecial(card *match.Card) {

	card.RemoveConditionBySource(card.ID + "-custom")

	graveyard, err := card.Player.Container(match.GRAVEYARD)

	if err != nil {
		return
	}

	power := 0

	for _, graveyardCard := range graveyard {
		if graveyardCard.Civ == civ.Fire {
			power += 1000
		}
	}

	card.AddCondition(cnd.PowerAttacker, power, card.ID+"-custom")

}
