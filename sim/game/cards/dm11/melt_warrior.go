package dm11

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// HysteriaLizard ...
func HysteriaLizard(c *match.Card) {

	c.Name = "Hysteria Lizard"
	c.Power = 3000
	c.Civs = []string{civ.Fire}
	c.Family = []string{family.MeltWarrior}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.ForceAttack, fx.PowerAttacker3000)

}

// BonfireLizard ...
func BonfireLizard(c *match.Card) {

	c.Name = "Bonfire Lizard"
	c.Power = 4000
	c.Civs = []string{civ.Fire}
	c.Family = []string{family.MeltWarrior}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.WaveStriker, fx.WhileWaveStriker(fx.When(fx.Summoned, fx.DestroyUpToXOpBlockers(2))))

}

// LockdownLizard ...
func LockdownLizard(c *match.Card) {

	c.Name = "Lockdown Lizard"
	c.Power = 3000
	c.Civs = []string{civ.Fire}
	c.Family = []string{family.MeltWarrior}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.When(fx.InTheBattlezone, func(card *match.Card, ctx *match.Context) {
		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {

			if card.Zone != match.BATTLEZONE {
				exit()
				return
			}

			// A persistent effect runs before every card handler in the same
			// context, and the flow is cancelled before any of them, so the
			// block does not depend on handler order. fx.Creature resolves a
			// tap ability in a scheduled callback, which a cancelled context
			// never reaches.
			//
			// Persistent effects do not see each other's cancellations, so a
			// future card that reacted to a tap ability with a side effect
			// rather than a refusal would still act. Nothing does today.
			event, ok := ctx2.Event.(*match.TapAbility)

			if !ok {
				return
			}

			user, err := ctx2.Match.CurrentPlayer().Player.GetCard(event.CardID, match.BATTLEZONE)

			if err != nil {
				return
			}

			// "Players can't" is both of them, its own controller included.
			ctx2.Match.WarnPlayer(user.Player, fmt.Sprintf("%s can't use tap abilities because of %s", user.Name, card.Name))
			ctx2.InterruptFlow()
		})
	}))

}
