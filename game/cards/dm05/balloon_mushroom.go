package dm05

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// BalloonshroomQ ...
func BalloonshroomQ(c *match.Card) {

	c.Name = "Balloonshroom Q"
	c.Power = 2000
	c.Civ = civ.Nature
	c.Family = []string{family.BalloonMushroom}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.Survivor, func(card *match.Card, ctx *match.Context) {

		if card.Zone != match.BATTLEZONE {
			return
		}

		if event, ok := ctx.Event.(*match.CreatureDestroyed); ok &&
			event.Card.Player == card.Player &&
			event.Card.HasCondition(cnd.Survivor) {

			ctx.InterruptFlow()

			card.Player.MoveCard(event.Card.ID, match.BATTLEZONE, match.MANAZONE)
			event.Card.Tapped = false
			ctx.Match.Chat("Server", fmt.Sprintf("%s was destroyed by %s and moved to the mana zone", event.Card.Name, event.Source.Name))

		}

	})
}
