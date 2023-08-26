package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// BolzardDragon ...
func BolzardDragon(c *match.Card) {

	c.Name = "Bolzard Dragon"
	c.Power = 5000
	c.Civ = civ.Fire
	c.Family = []string{family.ArmoredDragon}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.AttackConfirmed); ok {
			if event.CardID != card.ID {
				return
			}
			manaCards := match.Search(
				card.Player,
				ctx.Match,
				ctx.Match.Opponent(card.Player),
				match.MANAZONE,
				"Select 1 card from your opponent's mana zone that will be sent to their graveyard",
				1,
				1,
				false,
			)

			for _, mana := range manaCards {
				mana.Player.MoveCard(mana.ID, match.MANAZONE, match.GRAVEYARD)
				ctx.Match.Chat("Server", fmt.Sprintf("%s was sent from %s's manazone to their graveyard by %s", mana.Name, mana.Player.Username(), card.Name))
			}
		}

	})

}
