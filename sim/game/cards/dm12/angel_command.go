package dm12

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// ValkyerStarstormElemental ...
func ValkyerStarstormElemental(c *match.Card) {

	c.Name = "Valkyer, Starstorm Elemental"
	c.Power = 7000
	c.Civs = []string{civ.Light}
	c.Family = []string{family.AngelCommand}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Blocker(), fx.CantAttackPlayers)

}

// KilstineNebulaElemental ...
func KilstineNebulaElemental(c *match.Card) {

	c.Name = "Kilstine, Nebula Elemental"
	c.Power = 5000
	c.Civs = []string{civ.Light}
	c.Family = []string{family.AngelCommand}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.WaveStriker, func(card *match.Card, ctx *match.Context) {
		active := fx.WaveStrikerActive(card, ctx.Match)

		// Swept across every zone when inactive, so a creature that left the
		// battle zone cannot carry a stale bonus back into it.
		fx.FindMultiple(card.Player, []string{match.BATTLEZONE, match.HAND, match.GRAVEYARD, match.MANAZONE, match.SHIELDZONE, match.HIDDENZONE}).Map(func(x *match.Card) {
			if !active || x.Zone != match.BATTLEZONE || x.ID == card.ID {
				x.RemoveSpecificConditionBySource(cnd.PowerAmplifier, card.ID)
				x.RemoveSpecificConditionBySource(cnd.DoubleBreaker, card.ID)
				x.RemoveSpecificConditionBySource(cnd.Blocker, card.ID)
				return
			}

			x.AddUniqueSourceCondition(cnd.PowerAmplifier, 5000, card.ID)
			x.AddUniqueSourceCondition(cnd.DoubleBreaker, true, card.ID)

			// ForceBlocker both marks the creature and offers it for the attack
			// being resolved, which a granted blocker needs in order to block.
			fx.ForceBlocker(x, ctx, card.ID)
		})
	})

}

// UlarusPunishmentElemental ...
func UlarusPunishmentElemental(c *match.Card) {

	c.Name = "Ularus, Punishment Elemental"
	c.Power = 4500
	c.Civs = []string{civ.Light}
	c.Family = []string{family.AngelCommand}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		// Ularus is in the battle zone by now, so it counts towards its own
		// allowance.
		allowance := len(fx.Find(card.Player, match.BATTLEZONE))

		if allowance < 1 {
			return
		}

		opponent := ctx.Match.Opponent(card.Player)
		faceDown := func(x *match.Card) bool { return !x.ShieldFaceUp }

		// "A shield" names no owner, so either shield zone is fair game. They
		// are offered backside only: a shield the caster has not turned face up
		// yet must stay hidden.
		shields := map[string][]*match.Card{
			"Your shields":            fx.FindFilter(card.Player, match.SHIELDZONE, faceDown),
			"Your opponent's shields": fx.FindFilter(opponent, match.SHIELDZONE, faceDown),
		}

		turned := fx.SelectMultipartBackside(
			card.Player,
			ctx.Match,
			shields,
			fmt.Sprintf("%s's effect: You may turn up to %d shields face up, one for each creature you have in the battle zone.", card.Name, allowance),
			1,
			allowance,
			true,
		)

		turned.Map(func(shield *match.Card) {
			shield.ShieldFaceUp = true
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("A shield of %s was turned face up by %s", shield.Player.Username(), card.Name))
		})

		if len(turned) > 0 {
			ctx.Match.BroadcastState()
		}
	}))

}
