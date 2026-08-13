package dm11

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// SpinningTerrorTheWretched ...
func SpinningTerrorTheWretched(c *match.Card) {

	c.Name = "Spinning Terror, the Wretched"
	c.Power = 1000
	c.Civs = []string{civ.Darkness}
	c.Family = []string{family.DevilMask}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Darkness}

	c.PowerModifier = func(m *match.Match, attacking bool) int {
		if c.Zone != match.BATTLEZONE {
			return 0
		}

		// Recounted on every read rather than cached, because creatures tap and
		// untap constantly and the bonus has to follow.
		return len(fx.FindFilter(
			m.Opponent(c.Player),
			match.BATTLEZONE,
			func(x *match.Card) bool { return x.Tapped },
		)) * 2000
	}

	c.Use(fx.Creature)

}

// EvilIncarnate ...
func EvilIncarnate(c *match.Card) {

	c.Name = "Evil Incarnate"
	c.Power = 11000
	c.Civs = []string{civ.Darkness}
	c.Family = []string{family.DevilMask}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Evolution, fx.Doublebreaker, func(card *match.Card, ctx *match.Context) {

		if _, ok := ctx.Event.(*match.StartOfTurnStep); !ok || card.Zone != match.BATTLEZONE {
			return
		}

		// "Each player's turn" includes its controller's, and the player who
		// has to give something up is whoever the turn belongs to. Scheduled so
		// the rest of the start of turn has resolved first.
		ctx.ScheduleAfter(func() {
			if card.Zone != match.BATTLEZONE {
				return
			}

			fx.PlayerChoosesAndDestroysOwnCreature(ctx.Match.CurrentPlayer().Player, card, ctx)
		})
	})

}
