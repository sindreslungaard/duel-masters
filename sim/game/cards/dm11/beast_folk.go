package dm11

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// QuillspikeRumbler ...
func QuillspikeRumbler(c *match.Card) {

	c.Name = "Quillspike Rumbler"
	c.Power = 3000
	c.Civs = []string{civ.Nature}
	c.Family = []string{family.BeastFolk}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {
		event, ok := ctx.Event.(*match.AttackConfirmed)

		if !ok || event.CardID != card.ID || !event.Creature || card.Zone != match.BATTLEZONE {
			return
		}

		// Added before the battle captures the attacker's power, and cleared
		// along with the card's other conditions at the end of the turn.
		card.AddCondition(cnd.PowerAmplifier, 3000, card.ID)
		ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s gets +3000 power until the end of the turn", card.Name))
	})

}
