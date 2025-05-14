package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

func BazagazealDragon(c *match.Card) {
	c.Name = "Bazagazeal Dragon"
	c.Power = 8000
	c.Civ = civ.Fire
	c.Family = []string{family.ArmoredDragon}
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.Doublebreaker, fx.SpeedAttacker, fx.AttackUntapped, fx.When(fx.EndOfMyTurn, func(card *match.Card, ctx *match.Context) {
		if card.Zone != match.BATTLEZONE {
			return
		}

		ctx.ScheduleAfter(func() {
			card.Player.MoveCard(card.ID, match.BATTLEZONE, match.HAND, card.ID)
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was returned to the %s's hand", c.Name, c.Player.Username()))
		})
	}))

}

func BolmeteusSteelDragon(c *match.Card) {
	c.Name = "Bolmeteus Steel Dragon"
	c.Power = 7000
	c.Civ = civ.Fire
	c.Family = []string{family.ArmoredDragon}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.Doublebreaker, fx.When(fx.BreakShield, func(card *match.Card, ctx *match.Context) {

		event, _ := ctx.Event.(*match.BreakShieldEvent)
		if event.Source != card {
			return
		}

		ctx.InterruptFlow()
		for _, shield := range event.Cards {
			moved, err := ctx.Match.Opponent(card.Player).MoveCard(shield.ID, match.SHIELDZONE, match.GRAVEYARD, card.ID)
			if err == nil {
				ctx.Match.ReportActionInChat(ctx.Match.Opponent(card.Player), fmt.Sprintf("%s was moved to %s's graveyard instead of hand by %s", moved.Name, ctx.Match.Opponent(card.Player).Username(), card.Name))
			}
		}

	}))
}
