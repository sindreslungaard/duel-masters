package dm12

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// Meloppe ...
func Meloppe(c *match.Card) {

	c.Name = "Meloppe"
	c.Power = 1000
	c.Civs = []string{civ.Water}
	c.Family = []string{family.CyberLord}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Water}

	// Breaking shields always has the attacker pick which of the defender's
	// face-down shields come off (see SelectAndReturnShields), so both of
	// Meloppe's clauses reduce to the same rule regardless of which side
	// controls it: while Meloppe is in play, the shield owner picks their own
	// shields instead of the attacker picking for them.
	c.Use(fx.Creature, fx.When(fx.InTheBattlezone, func(card *match.Card, ctx *match.Context) {
		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
			if card.Zone != match.BATTLEZONE {
				exit()
				return
			}

			event, ok := ctx2.Event.(*match.SelectShields)
			if !ok {
				return
			}

			event.Chooser = ctx2.Match.Opponent(event.Attacker.Player)
		})
	}))

}
