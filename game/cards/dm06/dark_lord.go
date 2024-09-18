package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

func SchukaDukeOfAmnesia(c *match.Card) {

	c.Name = "Schuka, Duke of Amnesia"
	c.Power = 5000
	c.Civ = civ.Darkness
	c.Family = []string{family.DarkLord}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.When(fx.Destroyed, func(card *match.Card, ctx *match.Context) {

		for _, c := range append(fx.Find(card.Player, match.HAND), fx.Find(ctx.Match.Opponent(card.Player), match.HAND)...) {
			c.Player.MoveCard(c.ID, match.HAND, match.GRAVEYARD, card.ID)
			ctx.Match.ReportActionInChat(c.Player, fmt.Sprintf("%s was was discarded from %s's hand by %s", c.Name, c.Player.Username(), card.Name))
		}
	}))
}
