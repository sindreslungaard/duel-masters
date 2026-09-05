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

	// Both of Meloppe's clauses reduce to the same rule regardless of which
	// side controls it: while Meloppe is in play, a card's own owner always
	// picks which of their own face-down shields is chosen, instead of
	// whoever the effect's default chooser would otherwise be.
	//
	// Breaking shields always has the attacker pick which of the defender's
	// face-down shields come off (see SelectAndReturnShields), so that case is
	// handled through SelectShields.Chooser. Every other "choose one of a
	// player's shields" effect (looking at one, discarding one, and so on)
	// goes through ChooseShieldEvent instead (see fx.ResolveShieldChooser), so
	// this same effect intercepts both.
	c.Use(fx.Creature, fx.When(fx.InTheBattlezone, func(card *match.Card, ctx *match.Context) {
		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
			if card.Zone != match.BATTLEZONE {
				exit()
				return
			}

			switch event := ctx2.Event.(type) {
			case *match.SelectShields:
				event.Chooser = ctx2.Match.Opponent(event.Attacker.Player)
			case *match.ChooseShieldEvent:
				event.Chooser = event.ShieldOwner
			}
		})
	}))

}
