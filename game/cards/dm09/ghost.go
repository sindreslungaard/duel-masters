package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// IceVaporShadowOfAnguish ...
func IceVaporShadowOfAnguish(c *match.Card) {

	c.Name = "Ice Vapor, Shadow of Anguish"
	c.Power = 1000
	c.Civ = civ.Darkness
	c.Family = []string{family.Ghost}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.When(fx.OppSpellCast, func(card *match.Card, ctx *match.Context) {

		if card.Zone != match.BATTLEZONE {
			return
		}

		//@TODO see if another event must be fired, so we know when the actual spell effect
		// has been resolved (something like AFTER OppSpellCast) !!
		// in SpellCast handler
		ctx.ScheduleAfter(func() {
			fx.Select(
				ctx.Match.Opponent(card.Player),
				ctx.Match,
				ctx.Match.Opponent(card.Player),
				match.HAND,
				fmt.Sprintf("%s's effect: Choose a card from your hand and discard it.", card.Name),
				1,
				1,
				false,
			).Map(func(x *match.Card) {
				ctx.Match.Opponent(card.Player).MoveCard(x.ID, match.HAND, match.GRAVEYARD, card.ID)
				ctx.Match.BroadcastState()
				ctx.Match.ReportActionInChat(ctx.Match.Opponent(card.Player), fmt.Sprintf("%s was discarded from %s's hand by %s's effect.", x.Name, ctx.Match.Opponent(card.Player).Username(), card.Name))
			})

			fx.Select(
				ctx.Match.Opponent(card.Player),
				ctx.Match,
				ctx.Match.Opponent(card.Player),
				match.MANAZONE,
				fmt.Sprintf("%s's effect: Choose a card from your mana zone and put it into your graveyard.", card.Name),
				1,
				1,
				false,
			).Map(func(x *match.Card) {
				ctx.Match.Opponent(card.Player).MoveCard(x.ID, match.MANAZONE, match.GRAVEYARD, card.ID)
				ctx.Match.BroadcastState()
				ctx.Match.ReportActionInChat(ctx.Match.Opponent(card.Player), fmt.Sprintf("%s was put into the graveyard from %s's mana zone by %s's effect.", x.Name, ctx.Match.Opponent(card.Player).Username(), card.Name))
			})
		})

	}))

}
