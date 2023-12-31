package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

func TankMutant(c *match.Card) {

	c.Name = "Tank Mutant"
	c.Power = 6000
	c.Civ = civ.Darkness
	c.Family = []string{family.Hedrian}
	c.ManaCost = 9
	c.ManaRequirement = []string{civ.Darkness}
	c.TapAbility = true

	c.Use(fx.Creature, fx.When(fx.TapAbility, func(card *match.Card, ctx *match.Context) {
		opponentCreatures := match.Search(ctx.Match.Opponent(card.Player), ctx.Match, ctx.Match.Opponent(card.Player), match.BATTLEZONE, "Tank Mutant: Select 1 creature from your battlezone that will be sent to your graveyard", 1, 1, false)

		for _, creature := range opponentCreatures {
			ctx.Match.Destroy(creature, card, match.DestroyedByMiscAbility)
		}
		card.Tapped = true
	}))

}

func BazookaMutant(c *match.Card) {

	c.Name = "Bazooka Mutant"
	c.Power = 8000
	c.Civ = civ.Darkness
	c.Family = []string{family.Hedrian}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.CantAttackPlayers, func(card *match.Card, ctx *match.Context) {

		if card.Zone != match.BATTLEZONE {
			return
		}

		if event, ok := ctx.Event.(*match.AttackCreature); ok {

			if !ok || event.CardID != card.ID {
				return
			}

			for _, creature := range event.AttackableCreatures {

				if !creature.HasCondition(cnd.Blocker) {

					ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("%s can't attack creatures other than blockers", card.Name))
					ctx.InterruptFlow()
					return
				}

			}

		}
	})
}
