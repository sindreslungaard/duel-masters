package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// GajirabuteVileCenturion ...
func GajirabuteVileCenturion(c *match.Card) {
	c.Name = "Gajirabute, Vile Centurion"
	c.Power = 3000
	c.Civ = civ.Darkness
	c.Family = []string{family.DemonCommand}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		opponent := ctx.Match.Opponent(card.Player)
		fx.SelectBackside(
			card.Player,
			ctx.Match,
			opponent,
			match.SHIELDZONE,
			fmt.Sprintf("%s's effect: Choose one of your opponent's shields and put it into their graveyard.", card.Name),
			1,
			1,
			false,
		).Map(func(shield *match.Card) {
			moved, err := opponent.MoveCard(shield.ID, match.SHIELDZONE, match.GRAVEYARD, card.ID)
			if err == nil && moved.Zone == match.GRAVEYARD {
				ctx.Match.ReportActionInChat(opponent, fmt.Sprintf("A shield was put into %s's graveyard by %s", opponent.Username(), card.Name))
			}
		})
	}))
}
