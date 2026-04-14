package dm03

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// GarkagoDragon ...
func GarkagoDragon(c *match.Card) {

	c.Name = "Garkago Dragon"
	c.Power = 6000
	c.Civ = civ.Fire
	c.Family = []string{family.ArmoredDragon}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Fire}

	// Get +1000 power for each other fire card in your battle zone
	c.PowerModifier = func(m *match.Match, attacking bool) int {
		return (getFireCardsInYourBattleZone(c) - 1) * 1000 //-1 to exclude self
	}

	c.Use(fx.Creature, fx.Doublebreaker, fx.AttackUntapped)
}

// BoltailDragon ...
func BoltailDragon(c *match.Card) {

	c.Name = "Boltail Dragon"
	c.Power = 9000
	c.Civ = civ.Fire
	c.Family = []string{family.ArmoredDragon}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.Doublebreaker)
}

// UberdragonJabaha ...
func UberdragonJabaha(c *match.Card) {

	c.Name = "Uberdragon Jabaha"
	c.Power = 11000
	c.Civ = civ.Fire
	c.Family = []string{family.ArmoredDragon}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.Doublebreaker, fx.Evolution, func(card *match.Card, ctx *match.Context) {

		if card.Zone != match.BATTLEZONE {
			return
		}

		if event, ok := ctx.Event.(*match.AttackCreature); ok {
			uberDragonJabahaSpecial(card, ctx, event.CardID)
		}

		if event, ok := ctx.Event.(*match.AttackPlayer); ok {
			uberDragonJabahaSpecial(card, ctx, event.CardID)
		}
	})
}

func uberDragonJabahaSpecial(card *match.Card, ctx *match.Context, cardID string) {

	p := ctx.Match.CurrentPlayer()

	creature, err := p.Player.GetCard(cardID, match.BATTLEZONE)

	if err != nil {
		return
	}

	if creature.Civ != civ.Fire || creature.ID == card.ID {
		return
	}

	if ctx.Match.IsPlayerTurn(card.Player) {
		creature.RemoveConditionBySource(card.ID)
		creature.AddCondition(cnd.PowerAttacker, 2000, card.ID)
	}
}

// Return the number of water creatures in your battle zone
func getFireCardsInYourBattleZone(card *match.Card) int {

	battleZone, err := card.Player.Container(match.BATTLEZONE)

	if err != nil {
		return 0
	}

	count := 0

	for _, battleZoneCard := range battleZone {
		if battleZoneCard.Civ == civ.Fire {
			count++
		}
	}

	return count
}
