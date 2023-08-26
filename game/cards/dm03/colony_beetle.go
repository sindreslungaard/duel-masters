package dm03

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// PouchShell ...
func PouchShell(c *match.Card) {

	c.Name = "Pouch Shell"
	c.Power = 1000
	c.Civ = civ.Nature
	c.Family = []string{family.ColonyBeetle}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		if match.AmISummoned(card, ctx) {

			fx.SelectFilter(
				card.Player,
				ctx.Match,
				ctx.Match.Opponent(card.Player),
				match.BATTLEZONE,
				"Pouch Shell: Select an evolution card from your opponent's battle zone and send the top card to their graveyard",
				0,
				1,
				true,
				func(x *match.Card) bool { return x.HasCondition(cnd.Evolution) },
			).Map(func(x *match.Card) {
				tapped := x.Tapped
				baseCard := x.Attachments()[0]
				x.ClearAttachments()
				ctx.Match.Destroy(x, card, match.DestroyedByMiscAbility)
				baseCard.Player.MoveCard(baseCard.ID, match.HIDDENZONE, match.BATTLEZONE)
				if tapped && !baseCard.Tapped {
					baseCard.Tapped = true
				}
			})
		}

	})

}
