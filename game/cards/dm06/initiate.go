package dm06

import (
	"fmt"

	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func CrazeValkyrieTheDrastic(c *match.Card) {

	c.Name = "Craze Valkyrie, the Drastic"
	c.Power = 7500
	c.Civ = civ.Light
	c.Family = []string{family.Initiate}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Evolution, fx.Doublebreaker, func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.CardMoved); ok {

			if event.CardID == card.ID && event.To == match.BATTLEZONE {

				creatures := match.Search(card.Player, ctx.Match, ctx.Match.Opponent(card.Player), match.BATTLEZONE, "Craze Valkyrie, the Drastic: Choose up to 2 of your opponent's creature and tap them. Close to not tap any creatures.", 1, 2, true)

				for _, creature := range creatures {
					creature.Tapped = true

					ctx.Match.ReportActionInChat(ctx.Match.Opponent(card.Player), fmt.Sprintf("%s was tapped by %s's %s", creature.Name, card.Player.Username(), card.Name))
				}

			}

		}

	})

}

func BallasVizierOfElectrons(c *match.Card) {

	c.Name = "Ballas, Vizier of Electrons"
	c.Power = 2000
	c.Civ = civ.Light
	c.Family = []string{family.Initiate}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature)
}

func ChekiculVizierOfEndurance(c *match.Card) {

	c.Name = "Chekicul, Vizier of Endurance"
	c.Power = 1000
	c.Civ = civ.Light
	c.Family = []string{family.Initiate}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Blocker, fx.CantAttackCreatures, fx.CantAttackPlayers, func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.Battle); ok {
			if !event.Blocked || event.Defender != card {
				return
			}

			ctx.InterruptFlow()

			event.Attacker.Tapped = true
			event.Defender.Tapped = true

		}
	})
}

func ChenTregVizierOfBlades(c *match.Card) {

	c.Name = "Chen Treg, Vizier of Blades"
	c.Power = 2000
	c.Civ = civ.Light
	c.Family = []string{family.Initiate}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}
	c.TapAbility = func(card *match.Card, ctx *match.Context) {

		creatures := match.Search(card.Player, ctx.Match, ctx.Match.Opponent(card.Player), match.BATTLEZONE, "Chen Treg, Vizier of Blades: Select 1 of your opponent's creature and tap it.", 1, 1, false)

		for _, creature := range creatures {
			creature.Tapped = true
		}

	}

	c.Use(fx.Creature, fx.TapAbility)
}
