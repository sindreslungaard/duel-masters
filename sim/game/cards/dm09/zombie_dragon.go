package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// NecrodragonIzoristVhal ...
func NecrodragonIzoristVhal(c *match.Card) {

	c.Name = "Necrodragon Izorist Vhal"
	c.Power = 0
	c.Civs = []string{civ.Darkness}
	c.Family = []string{family.ZombieDragon}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Darkness}

	addPower := 0

	c.Use(fx.Creature, fx.When(fx.InTheBattlezone, func(card *match.Card, ctx *match.Context) {
		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
			if card.Zone != match.BATTLEZONE {
				c.Power = 0
				c.RemoveConditionBySource(card.ID)
				exit()
				return
			}

			// fx.Creature clears every creature's conditions (cnd.Creature included)
			// on EndOfTurnStep and only rebuilds cnd.Creature during UntapStep's own
			// handler pass. BeginTurnStep, UntapManaEvent, and UntapStep itself all run
			// with that identity still missing match-wide (persistent effects run
			// before per-card handlers in the same event), so recomputing here would
			// see an empty graveyard and destroy this creature every turn transition.
			// Skip those and keep the last known-good power until StartOfTurnStep,
			// once conditions are rebuilt. See fx.SilentSkill for the same gotcha.
			switch ctx2.Event.(type) {
			case *match.BeginTurnStep, *match.UntapManaEvent, *match.UntapStep:
				return
			}

			c.Power = 0
			c.RemoveConditionBySource(card.ID)

			addPower = len(fx.FindFilter(
				card.Player,
				match.GRAVEYARD,
				func(x *match.Card) bool {
					return x.HasCondition(cnd.Creature) && x.HasCiv(civ.Darkness)
				})) * 2000

			if addPower == 0 {
				exit()
				ctx2.Match.Destroy(card, card, match.DestroyedByMiscAbility)
				return
			}

			c.Power += addPower

			if c.Power >= 6000 {
				c.AddUniqueSourceCondition(cnd.DoubleBreaker, true, card.ID)
			}
		})
	}))

}
