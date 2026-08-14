package dm11

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// MiraculousTruce ...
func MiraculousTruce(c *match.Card) {

	c.Name = "Miraculous Truce"
	c.Civs = []string{civ.Light, civ.Nature}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light, civ.Nature}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.PutIntoManaZoneTapped, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		civilizations := []string{civ.Light, civ.Water, civ.Darkness, civ.Fire, civ.Nature}

		chosen := fx.MultipleChoiceQuestion(
			card.Player,
			ctx.Match,
			fmt.Sprintf("%s's effect: Choose a civilization. Its creatures can't attack you until the start of your next turn.", card.Name),
			civilizations,
		)

		if chosen < 0 || chosen >= len(civilizations) {
			return
		}

		barred := civilizations[chosen]
		ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s: %s creatures can't attack %s until the start of their next turn", card.Name, barred, card.Player.Username()))

		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {

			// "Until the start of your next turn". The caster's own turn has
			// already begun if this was cast on their turn, so the next one to
			// arrive is always the turn the effect should end on. A shield
			// trigger cast on the opponent's turn lands on the same rule.
			if _, ok := ctx2.Event.(*match.BeginTurnStep); ok && ctx2.Match.IsPlayerTurn(card.Player) {
				exit()
				return
			}

			event, ok := ctx2.Event.(*match.AttackPlayer)

			if !ok {
				return
			}

			attacker, err := ctx2.Match.Opponent(card.Player).GetCard(event.CardID, match.BATTLEZONE)

			// Only attacks against the caster are stopped, so an attacker that
			// does not belong to their opponent is none of this effect's
			// business.
			if err != nil || !attacker.HasCiv(barred) {
				return
			}

			ctx2.Match.WarnPlayer(attacker.Player, fmt.Sprintf("%s can't attack %s because of %s", attacker.Name, card.Player.Username(), card.Name))
			ctx2.InterruptFlow()
		})
	}))

}

// MiraculousMeltdown ...
func MiraculousMeltdown(c *match.Card) {

	c.Name = "Miraculous Meltdown"
	c.Civs = []string{civ.Darkness, civ.Fire}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Darkness, civ.Fire}

	c.Use(
		fx.Spell,
		fx.PutIntoManaZoneTapped,
		func(card *match.Card, ctx *match.Context) {

			// "You can cast this spell only if your opponent has more shields
			// than you do." Cancelling the play event stops the cast before any
			// mana is paid, because fx.Spell only takes payment in a callback
			// scheduled after this traversal.
			event, ok := ctx.Event.(*match.PlayCardEvent)

			if !ok || event.CardID != card.ID {
				return
			}

			if shieldCount(card.Player) < shieldCount(ctx.Match.Opponent(card.Player)) {
				return
			}

			ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("%s can only be cast while your opponent has more shields than you do", card.Name))
			ctx.InterruptFlow()
		},
		fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
			opponent := ctx.Match.Opponent(card.Player)

			keep := shieldCount(card.Player)

			shields, err := opponent.Container(match.SHIELDZONE)

			if err != nil || len(shields) <= keep {
				return
			}

			kept := fx.SelectBackside(
				opponent,
				ctx.Match,
				opponent,
				match.SHIELDZONE,
				fmt.Sprintf("%s: Choose %d of your shields to keep. The rest go to your hand.", card.Name, keep),
				keep,
				keep,
				false,
			)

			keptIDs := make(map[string]bool, len(kept))
			for _, shield := range kept {
				keptIDs[shield.ID] = true
			}

			// Snapshotted before anything moves, so the walk is not disturbed
			// by the shield zone shrinking underneath it.
			losing := make([]*match.Card, 0, len(shields))
			for _, shield := range shields {
				if !keptIDs[shield.ID] {
					losing = append(losing, shield)
				}
			}

			triggers := make([]*match.Card, 0, len(losing))
			for _, shield := range losing {
				moved, err := opponent.MoveCard(shield.ID, match.SHIELDZONE, match.HAND, card.ID)

				if err != nil || moved.Zone != match.HAND {
					continue
				}

				if moved.HasCondition(cnd.ShieldTrigger) {
					triggers = append(triggers, moved)
				}
			}

			ctx.Match.ReportActionInChat(opponent, fmt.Sprintf("%s put %d of their shields into their hand because of %s", opponent.Username(), len(losing), card.Name))

			// The reminder text is explicit that these shields still offer
			// their triggers even though they were not broken.
			ctx.Match.ResolveShieldTriggers(triggers, card)
		}),
	)

}

// MiraculousRebirth ...
func MiraculousRebirth(c *match.Card) {

	c.Name = "Miraculous Rebirth"
	c.Civs = []string{civ.Fire, civ.Nature}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Fire, civ.Nature}

	c.Use(fx.Spell, fx.PutIntoManaZoneTapped, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		fx.SelectFilter(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			fmt.Sprintf("%s's effect: Destroy one of your opponent's creatures that has power 5000 or less.", card.Name),
			1,
			1,
			false,
			func(creature *match.Card) bool { return ctx.Match.GetPower(creature, false) <= 5000 },
			false,
		).Map(func(victim *match.Card) {
			// Read before the destruction, because a card in the graveyard is
			// no longer the creature that was in play.
			cost := victim.ManaCost

			ctx.Match.Destroy(victim, card, match.DestroyedBySpell)

			// "When your opponent puts that creature into his graveyard" is a
			// condition, not a certainty: a replacement effect can send it
			// somewhere else, and then there is nothing to match.
			if victim.Zone != match.GRAVEYARD {
				return
			}

			fx.SearchDeckPutCreatureIntoBZ(
				card,
				ctx,
				func(x *match.Card) bool { return x.ManaCost == cost },
				fmt.Sprintf("a creature that costs %d", cost),
			)
		})
	}))

}

// MiraculousPlague ...
func MiraculousPlague(c *match.Card) {

	c.Name = "Miraculous Plague"
	c.Civs = []string{civ.Water, civ.Darkness}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Water, civ.Darkness}

	c.Use(fx.Spell, fx.PutIntoManaZoneTapped, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		opponent := ctx.Match.Opponent(card.Player)

		creatures := fx.Select(
			card.Player,
			ctx.Match,
			opponent,
			match.BATTLEZONE,
			fmt.Sprintf("%s's effect: Choose 2 of your opponent's creatures. They keep one and the other is destroyed.", card.Name),
			2,
			2,
			false,
		)

		fx.OpponentKeepsOneOfTheseAndLosesTheOther(
			card,
			ctx,
			creatures,
			match.BATTLEZONE,
			fmt.Sprintf("%s: Choose one of these creatures to return to your hand. The other one is destroyed.", card.Name),
			func(lost *match.Card) {
				ctx.Match.Destroy(lost, card, match.DestroyedBySpell)
			},
		)

		mana := fx.Select(
			card.Player,
			ctx.Match,
			opponent,
			match.MANAZONE,
			fmt.Sprintf("%s's effect: Choose 2 cards in your opponent's mana zone. They keep one and the other is discarded.", card.Name),
			2,
			2,
			false,
		)

		fx.OpponentKeepsOneOfTheseAndLosesTheOther(
			card,
			ctx,
			mana,
			match.MANAZONE,
			fmt.Sprintf("%s: Choose one of these cards to return to your hand. The other one goes to your graveyard.", card.Name),
			func(lost *match.Card) {
				moved, err := opponent.MoveCard(lost.ID, match.MANAZONE, match.GRAVEYARD, card.ID)

				if err == nil && moved.Zone == match.GRAVEYARD {
					ctx.Match.ReportActionInChat(opponent, fmt.Sprintf("%s was put into %s's graveyard by %s", lost.Name, opponent.Username(), card.Name))
				}
			},
		)
	}))

}

// shieldCount is how many shields a player is holding.
func shieldCount(player *match.Player) int {
	shields, err := player.Container(match.SHIELDZONE)

	if err != nil {
		return 0
	}

	return len(shields)
}
