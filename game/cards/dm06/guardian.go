package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func LuGilaSilverRiftGuardian(c *match.Card) {

	c.Name = "Lu Gila, Silver Rift Guardian"
	c.Power = 4000
	c.Civ = civ.Light
	c.Family = []string{family.Guardian}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Blocker, fx.CantAttackPlayers, func(card *match.Card, ctx *match.Context) {

		if card.Zone != match.BATTLEZONE && card.Zone != match.HIDDENZONE {
			return
		}

		if event, ok := ctx.Event.(*match.CardMoved); ok {
			if event.To != match.BATTLEZONE {
				return
			}

			playedCard, err := ctx.Match.Player1.Player.GetCard(event.CardID, match.BATTLEZONE)
			if err != nil {
				playedCard, err = ctx.Match.Player2.Player.GetCard(event.CardID, match.BATTLEZONE)
				if err != nil {
					return
				}
			}

			if !playedCard.HasCondition(cnd.Evolution) {
				return
			}

			// If Lu Gila is the card evolved, its evolution also becomes tapped according to duel master
			// rulings (https://duelmasters.fandom.com/wiki/Lu_Gila,_Silver_Rift_Guardian/Rulings).
			at := playedCard.Attachments()
			if card.Zone == match.HIDDENZONE && (len(at) == 0 || at[len(at)-1] != card) {
				return
			}

			playedCard.Tapped = true
		}

	})

}
