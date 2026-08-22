package promo

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// StarCryDragon ...
func StarCryDragon(c *match.Card) {

	c.Name = "Star-Cry Dragon"
	c.Power = 8000
	c.Civs = []string{civ.Fire}
	c.Family = []string{family.ArmoredDragon}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.Doublebreaker, func(card *match.Card, ctx *match.Context) {

		if card.Zone != match.BATTLEZONE {
			return
		}

		if event, ok := ctx.Event.(*match.GetPowerEvent); ok && event.Card.Player == card.Player && event.Card.ID != card.ID {

			if event.Card.HasFamily(family.ArmoredDragon) {
				event.Power += 3000
			}
		}

	})
}

// SuperDragonMachineDolzark ...
func SuperDragonMachineDolzark(c *match.Card) {

	c.Name = "Super Dragon Machine Dolzark"
	c.Power = 7000
	c.Civs = []string{civ.Fire, civ.Nature}
	c.Family = []string{family.ArmoredDragon, family.EarthDragon}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Fire, civ.Nature}

	c.Use(fx.Creature, fx.PutIntoManaZoneTapped, fx.Doublebreaker, func(card *match.Card, ctx *match.Context) {

		if card.Zone != match.BATTLEZONE {
			return
		}

		event, ok := ctx.Event.(*match.AttackConfirmed)
		if !ok {
			return
		}

		attacker, err := card.Player.GetCard(event.CardID, match.BATTLEZONE)
		if err != nil || attacker.ID == card.ID || attacker.Player != card.Player || !attacker.SharesAFamily(family.Dragons) {
			return
		}

		fx.MayPutOpCreatureXPowerOrLessIntoMana(5000)(card, ctx)
	})
}

// VelyrikaDragon ...
func VelyrikaDragon(c *match.Card) {

	c.Name = "Velyrika Dragon"
	c.Power = 7000
	c.Civs = []string{civ.Fire}
	c.Family = []string{family.ArmoredDragon}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.Doublebreaker, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		fx.SearchDeckTakeCards(
			card,
			ctx,
			1,
			func(x *match.Card) bool { return x.HasFamily(family.ArmoredDragon) },
			"Armored Dragon",
		)

	}))
}
