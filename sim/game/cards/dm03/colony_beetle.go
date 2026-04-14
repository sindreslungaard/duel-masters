package dm03

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// PouchShell ...
func PouchShell(c *match.Card) {

	c.Name = "Pouch Shell"
	c.Power = 1000
	c.Civ = civ.Nature
	c.Family = []string{family.ColonyBeetle}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		fx.SelectFilter(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			fmt.Sprintf("%s: Select an evolution card from your opponent's battle zone and send the top card to their graveyard", card.Name),
			0,
			1,
			true,
			func(x *match.Card) bool { return x.HasCondition(cnd.Evolution) },
			false,
		).Map(func(x *match.Card) {
			tapped := x.Tapped
			baseCards := x.Attachments()
			x.ClearAttachments()

			ctx.Match.Destroy(x, card, match.DestroyedByMiscAbility)

			for _, baseCard := range baseCards {
				baseCard.Player.MoveCard(baseCard.ID, match.HIDDENZONE, match.BATTLEZONE, card.ID)
				if tapped {
					baseCard.Tapped = true
				}
			}

			ctx.Match.BroadcastState()
		})

	}))

}
