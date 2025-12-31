package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

func KanesillTheExplorer(c *match.Card) {

	c.Name = "Kanesill, the Explorer"
	c.Power = 4000
	c.Civ = civ.Light
	c.Family = []string{family.Gladiator}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Blocker(), fx.CantAttackPlayers)
}

func TelitolTheExplorer(c *match.Card) {

	c.Name = "Telitol, the Explorer"
	c.Power = 3000
	c.Civ = civ.Light
	c.Family = []string{family.Gladiator}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Blocker(), fx.CantAttackPlayers, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		shields, err := card.Player.Container(match.SHIELDZONE)

		if err != nil {
			return
		}

		ids := make([]string, 0)

		for _, s := range shields {
			ids = append(ids, s.ImageID)
		}

		ctx.Match.ShowCards(
			card.Player,
			fmt.Sprintf("%s's effect: your shields:", card.Name),
			ids,
		)
	}))
}
