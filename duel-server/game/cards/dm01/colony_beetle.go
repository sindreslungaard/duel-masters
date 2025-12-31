package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// DomeShell ...
func DomeShell(c *match.Card) {

	c.Name = "Dome Shell"
	c.Power = 3000
	c.Civ = civ.Nature
	c.Family = []string{family.ColonyBeetle}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.PowerAttacker2000)

}

// StormShell ...
func StormShell(c *match.Card) {

	c.Name = "Storm Shell"
	c.Power = 2000
	c.Civ = civ.Nature
	c.Family = []string{family.ColonyBeetle}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		opponent := ctx.Match.Opponent(card.Player)

		battlezone, err := opponent.Container(match.BATTLEZONE)

		if err != nil {
			return
		}

		if len(battlezone) < 1 {
			return
		}

		selectedCards := fx.Select(
			opponent,
			ctx.Match,
			opponent,
			match.BATTLEZONE,
			fmt.Sprintf("%s: Select 1 card from your battlezone that will be sent to your manazone", card.Name),
			1,
			1,
			false,
		)

		movedCard, err := opponent.MoveCard(selectedCards[0].ID, match.BATTLEZONE, match.MANAZONE, card.ID)

		if err != nil {
			return
		}

		ctx.Match.ReportActionInChat(ctx.Match.Opponent(card.Player), fmt.Sprintf("%s was moved from %s's battlezone to their manazone", movedCard.Name, opponent.Username()))

	}))

}

// TowerShell ...
func TowerShell(c *match.Card) {

	c.Name = "Tower Shell"
	c.Power = 5000
	c.Civ = civ.Nature
	c.Family = []string{family.ColonyBeetle}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.CantBeBlockedByPowerUpTo4000)

}
