package promo

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// BrigadeShellQ ...
func BrigadeShellQ(c *match.Card) {

	c.Name = "Brigade Shell Q"
	c.Power = 1000
	c.Civ = civ.Nature
	c.Family = []string{family.ColonyBeetle, family.Survivor}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.Survivor, func(card *match.Card, ctx *match.Context) {

		if !ctx.Match.IsPlayerTurn(card.Player) || card.Zone != match.BATTLEZONE {
			return
		}

		event, ok := ctx.Event.(*match.AttackConfirmed)

		if !ok {
			return
		}

		creature, err := card.Player.GetCard(event.CardID, match.BATTLEZONE)

		if err != nil {
			return
		}

		if creature.HasCondition(cnd.Survivor) {
			cards := card.Player.PeekDeck(1)
			if len(cards) < 1 {
				return
			}

			c := cards[0]

			if c.HasFamily(family.Survivor) {
				fx.Draw1(card, ctx)
				ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s drew %s when %s attacked due to %s's survivor ability", card.Player.Username(), c.Name, creature.Name, card.Name))
			} else {
				c.Player.MoveCard(c.ID, match.DECK, match.GRAVEYARD, card.ID)
				ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s put %s in graveyard from their deck when %s attacked due to %s's survivor ability", card.Player.Username(), c.Name, creature.Name, card.Name))
			}
		}

	})

}
