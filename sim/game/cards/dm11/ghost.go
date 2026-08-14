package dm11

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// GazerEyesShadowOfSecrets ...
func GazerEyesShadowOfSecrets(c *match.Card) {

	c.Name = "Gazer Eyes, Shadow of Secrets"
	c.Power = 3000
	c.Civs = []string{civ.Darkness}
	c.Family = []string{family.Ghost}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.SilentSkill(func(card *match.Card, ctx *match.Context) {
		opponent := ctx.Match.Opponent(card.Player)

		// Selecting from the opponent's hand is what "look at" amounts to: the
		// cards are listed for the chooser and nobody else.
		fx.Select(
			card.Player,
			ctx.Match,
			opponent,
			match.HAND,
			fmt.Sprintf("%s's effect: Look at your opponent's hand and choose a card. They discard it.", card.Name),
			1,
			1,
			false,
		).Map(func(chosen *match.Card) {
			moved, err := opponent.MoveCard(chosen.ID, match.HAND, match.GRAVEYARD, card.ID)

			if err != nil || moved.Zone != match.GRAVEYARD {
				return
			}

			ctx.Match.ReportActionInChat(opponent, fmt.Sprintf("%s was discarded from %s's hand by %s", chosen.Name, opponent.Username(), card.Name))
		})
	}))

}
