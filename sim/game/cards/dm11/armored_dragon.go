package dm11

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
	"strings"
)

// HeavyweightDragon ...
func HeavyweightDragon(c *match.Card) {

	c.Name = "Heavyweight Dragon"
	c.Power = 9000
	c.Civs = []string{civ.Fire}
	c.Family = []string{family.ArmoredDragon}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Fire}

	c.TapAbility = func(card *match.Card, ctx *match.Context) {

		// Effective power is not something a player can read off the board once
		// power modifiers are involved, so a selection that does not make the
		// weight limit is explained and offered again rather than wasting the
		// ability. Cancelling is how the player gives up.
		for {
			ownPower := ctx.Match.GetPower(card, false)

			chosen := fx.SelectFilter(
				card.Player,
				ctx.Match,
				ctx.Match.Opponent(card.Player),
				match.BATTLEZONE,
				fmt.Sprintf("%s's effect: Choose up to 2 of your opponent's tapped creatures. They are destroyed only if their total power is less than %d.", card.Name, ownPower),
				1,
				2,
				true,
				func(creature *match.Card) bool { return creature.Tapped },
				false,
			)

			// Empty means the player cancelled, or that there was nothing
			// tapped to offer in the first place.
			if len(chosen) < 1 {
				return
			}

			// The whole selection is weighed as one: it is their combined power
			// that has to come in under this creature's, so either both die or
			// neither does.
			total := 0
			breakdown := make([]string, 0, len(chosen))

			for _, creature := range chosen {
				power := ctx.Match.GetPower(creature, false)
				total += power
				breakdown = append(breakdown, fmt.Sprintf("%s (%d)", creature.Name, power))
			}

			if total < ownPower {
				for _, creature := range chosen {
					ctx.Match.Destroy(creature, card, match.DestroyedByMiscAbility)
				}

				return
			}

			ctx.Match.WarnPlayer(card.Player, fmt.Sprintf(
				"%s: %s add up to %d total power, which is not less than %s's %d. Choose again, or close the selection to give up.",
				card.Name,
				strings.Join(breakdown, " + "),
				total,
				card.Name,
				ownPower,
			))
		}
	}

	c.Use(fx.Creature, fx.Doublebreaker, fx.TapAbility)

}
