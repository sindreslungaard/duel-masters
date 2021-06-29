package dm04

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// GregoriaPrincessOfWar ...
func GregoriaPrincessOfWar(c *match.Card) {

	c.Name = "Gregoria, Princess of War"
	c.Power = 5000
	c.Civ = civ.Darkness
	c.Family = family.DarkLord
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.CardMoved); ok {

			if event.CardID == card.ID && event.To == match.GRAVEYARD {

				getDemonCommands(card, ctx).Map(func(x *match.Card) {
					x.RemoveConditionBySource(card.ID)
				})
			}
		}

		if card.Zone != match.BATTLEZONE {
			return
		}

		getDemonCommands(card, ctx).Map(func(x *match.Card) {
			x.AddUniqueSourceCondition(cnd.PowerAttacker, 2000, card.ID)
			x.AddUniqueSourceCondition(cnd.Blocker, true, card.ID)
		})
	})

}

func getDemonCommands(card *match.Card, ctx *match.Context) fx.CardCollection {

	demonCommands := fx.FindFilter(
		card.Player,
		match.BATTLEZONE,
		func(x *match.Card) bool { return x.Family == family.DemonCommand },
	)

	demonCommands = append(demonCommands,

		fx.FindFilter(
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			func(x *match.Card) bool { return x.Family == family.DemonCommand },
		)...,
	)

	return demonCommands
}
