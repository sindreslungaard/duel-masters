package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func TankMutant(c *match.Card) {

	c.Name = "Tank Mutant"
	c.Power = 6000
	c.Civ = civ.Darkness
	c.Family = []string{family.Hedrian}
	c.ManaCost = 9
	c.ManaRequirement = []string{civ.Darkness}
	c.TapAbility = func(card *match.Card, ctx *match.Context) {

		ctx.Match.Wait(card.Player, "Waiting for your opponent to make an action")
		defer ctx.Match.EndWait(card.Player)

		fx.Select(
			ctx.Match.Opponent(card.Player),
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			"Tank Mutant: Select 1 creature from your battlezone that will be sent to your graveyard",
			1,
			1,
			false,
		).Map(func(creature *match.Card) {
			ctx.Match.Destroy(creature, card, match.DestroyedByMiscAbility)
		})
	}

	c.Use(fx.Creature, fx.TapAbility)

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

		if event, ok := ctx.Event.(*match.AttackCreature); ok && event.CardID == card.ID {

			canAU := card.HasCondition(cnd.AttackUntapped)

			attackableBlockers := fx.FindFilter(
				ctx.Match.Opponent(card.Player),
				match.BATTLEZONE,
				func(x *match.Card) bool {
					return x.HasCondition(cnd.Blocker) && (canAU || x.Tapped)
				},
			)

			if len(attackableBlockers) == 0 {
				ctx.Match.WarnPlayer(ctx.Match.Opponent(card.Player), "No blockers to attack.")
				ctx.InterruptFlow()
			}

			event.AttackableCreatures = attackableBlockers
		}
	})
}
