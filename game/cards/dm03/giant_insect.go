package dm03

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// Gigamantis ...
func Gigamantis(c *match.Card) {

	c.Name = "Gigamantis"
	c.Power = 5000
	c.Civ = civ.Nature
	c.Family = family.GiantInsect
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.Evolution, func(card *match.Card, ctx *match.Context) {
		if event, ok := ctx.Event.(*match.CardMoved); ok {

			if event.CardID == card.ID && event.From == match.BATTLEZONE && event.To == match.GRAVEYARD && card.HasCondition(cnd.Creature) {

				card.Player.MoveCard(card.ID, match.BATTLEZONE, match.MANAZONE)
				ctx.Match.Chat("Server", fmt.Sprintf("Gigamantis: %s was moved to mana zone instead of graveyard", card.Name))
			}
		}
	})

}

// SniperMosquito ...
func SniperMosquito(c *match.Card) {

	c.Name = "Sniper Mosquito"
	c.Power = 2000
	c.Civ = civ.Nature
	c.Family = family.GiantInsect
	c.ManaCost = 1
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.When(fx.AttackConfirmed, func(card *match.Card, ctx *match.Context) {

		fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.MANAZONE,
			"Select a card from your mana zone to be returned to your hand",
			1,
			1,
			false,
		).Map(func(x *match.Card) {
			ctx.Match.MoveCard(x, match.HAND, card)
		})
	}))
}

// SwordButterfly ...
func SwordButterfly(c *match.Card) {

	c.Name = "Sword Butterfly"
	c.Power = 2000
	c.Civ = civ.Nature
	c.Family = family.GiantInsect
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.PowerAttacker3000);
}
