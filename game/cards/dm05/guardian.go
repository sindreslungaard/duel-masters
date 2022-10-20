package dm05

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// SnorkLaShrineGuardian ...
func SnorkLaShrineGuardian(c *match.Card) {

	c.Name = "Snork La, Shrine Guardian"
	c.Power = 3000
	c.Civ = civ.Light
	c.Family = family.Guardian
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Blocker, fx.CantAttackPlayers, func(card *match.Card, ctx *match.Context) {

		if card.Zone != match.BATTLEZONE {
			return
		}

		event, ok := ctx.Event.(*match.CardMoved)

		if !ok {
			return
		}

		if event.From != match.MANAZONE && event.To != match.GRAVEYARD {
			return
		}

		x, err := card.Player.GetCard(event.CardID, match.GRAVEYARD)

		if err != nil {
			return
		}

		card.Player.MoveCard(x.ID, match.GRAVEYARD, match.MANAZONE)
		ctx.Match.Chat("Server", fmt.Sprintf("Snork La, Shrine Guardian prevented %s from being discarded from the manazone", x.Name))

	})

}
