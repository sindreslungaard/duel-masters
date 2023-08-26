package dm05

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// KipChippotto ...
func KipChippotto(c *match.Card) {

	c.Name = "Kip Chippotto"
	c.Power = 1000
	c.Civ = civ.Fire
	c.Family = []string{family.FireBird}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		if card.Zone != match.BATTLEZONE {
			return
		}

		if event, ok := ctx.Event.(*match.CreatureDestroyed); ok &&
			event.Card.ID != card.ID &&
			event.Card.Player == card.Player &&
			event.Card.HasFamily(family.ArmoredDragon) {

			fx.SelectFilter(
				card.Player,
				ctx.Match,
				card.Player,
				match.BATTLEZONE,
				fmt.Sprintf("%s: You may destroy this creature instead.", card.Name),
				1,
				1,
				true,
				func(c *match.Card) bool { return card.ID == c.ID },
			).Map(func(_ *match.Card) {

				ctx.InterruptFlow()

				ctx.Match.Destroy(card, card, match.DestroyedByMiscAbility)
				ctx.Match.Chat("Server", fmt.Sprintf("%s was destroyed instead of %s", card.Name, event.Card.Name))
			})

		}

	})

}
