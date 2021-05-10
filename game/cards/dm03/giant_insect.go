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
