package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// BubbleScarab ...
func BubbleScarab(c *match.Card) {

	c.Name = "Bubble Scarab"
	c.Civ = civ.Nature
	c.Power = 4000
	c.Family = []string{family.GiantInsect}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Nature}

	var myAttackedCreature *match.Card = nil

	c.Use(fx.Creature, fx.When(func(card *match.Card, ctx *match.Context) bool {
		if event, ok := ctx.Event.(*match.SelectBlockers); ok && card.Zone == match.BATTLEZONE {
			if event.AttackedCardID != "" {
				attackedCreature, _ := card.Player.GetCard(event.AttackedCardID, match.BATTLEZONE)

				if attackedCreature != nil {
					myAttackedCreature = attackedCreature
					return true
				}
			}
		}

		return false
	}, func(card *match.Card, ctx *match.Context) {
		if myAttackedCreature != nil {
			fx.Select(
				card.Player,
				ctx.Match,
				card.Player,
				match.HAND,
				fmt.Sprintf("%s's effect: You may discard a card from your hand. If you do, %s gets +3000 Power until the end of the turn.", card.Name, myAttackedCreature.Name),
				1,
				1,
				true,
			).Map(func(x *match.Card) {
				_, err := card.Player.MoveCard(x.ID, match.HAND, match.GRAVEYARD, card.ID)

				if err == nil {
					ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s discards %s to give %s +3000 Power until the end of the turn.", card.Player.Username(), x.Name, myAttackedCreature.Name))
					myAttackedCreature.AddUniqueSourceCondition(cnd.PowerAmplifier, 3000, card.ID)
					myAttackedCreature = nil
				}
			})
		}
	}))

}
