package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// TagtappTheRetaliator ...
func TagtappTheRetaliator(c *match.Card) {

	c.Name = "Tagtapp, the Retaliator"
	c.Power = 3000
	c.Civs = []string{civ.Fire, civ.Nature}
	c.Family = []string{family.SpiritQuartz}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Fire, civ.Nature}

	c.PowerModifier = func(m *match.Match, attacking bool) int {
		if c.Zone != match.BATTLEZONE {
			return 0
		}

		// A multicolored card in the opponent's mana zone counts as a water card
		// whenever water is one of its civilizations.
		return len(fx.FindFilter(
			m.Opponent(c.Player),
			match.MANAZONE,
			func(x *match.Card) bool { return x.HasCiv(civ.Water) },
		)) * 1000
	}

	c.Use(fx.Creature, fx.PutIntoManaZoneTapped, fx.PowerBreakerTiers(6000, 0))

}

// TanzanyteTheAwakener ...
func TanzanyteTheAwakener(c *match.Card) {

	c.Name = "Tanzanyte, the Awakener"
	c.Power = 9000
	c.Civs = []string{civ.Water, civ.Darkness}
	c.Family = []string{family.SpiritQuartz}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Water, civ.Darkness}

	c.TapAbility = func(card *match.Card, ctx *match.Context) {
		fx.SelectFilter(
			card.Player,
			ctx.Match,
			card.Player,
			match.GRAVEYARD,
			fmt.Sprintf("%s's effect: Choose a creature in your graveyard. Every creature with that name returns to your hand.", card.Name),
			1,
			1,
			false,
			func(x *match.Card) bool { return x.HasCondition(cnd.Creature) },
			false,
		).Map(func(chosen *match.Card) {
			// Snapshot the graveyard before moving anything out of it.
			matching := fx.FindFilter(
				card.Player,
				match.GRAVEYARD,
				func(x *match.Card) bool { return x.HasCondition(cnd.Creature) && x.Name == chosen.Name },
			)

			for _, creature := range matching {
				card.Player.MoveCard(creature.ID, match.GRAVEYARD, match.HAND, card.ID)
			}

			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s returned %d copies of %s from the graveyard to %s's hand", card.Name, len(matching), chosen.Name, card.Player.Username()))
		})
	}

	c.Use(fx.Creature, fx.PutIntoManaZoneTapped, fx.Doublebreaker, fx.TapAbility)

}

// SoderlightTheColdBlade ...
func SoderlightTheColdBlade(c *match.Card) {

	c.Name = "Soderlight, the Cold Blade"
	c.Power = 4000
	c.Civs = []string{civ.Water, civ.Darkness}
	c.Family = []string{family.SpiritQuartz}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Water, civ.Darkness}

	c.Use(fx.Creature, fx.PutIntoManaZoneTapped, fx.CantBeBlocked, fx.SilentSkill(fx.OpponentChoosesAndDestroysCreature))

}

// DeklowazTheTerminator ...
func DeklowazTheTerminator(c *match.Card) {

	c.Name = "Deklowaz, the Terminator"
	c.Power = 5000
	c.Civs = []string{civ.Darkness, civ.Fire}
	c.Family = []string{family.SpiritQuartz}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Darkness, civ.Fire}

	c.TapAbility = func(card *match.Card, ctx *match.Context) {

		fx.DestroyAllCreaturesXPowerOrLess(3000, match.DestroyedByMiscAbility)(card, ctx)

		opponent := ctx.Match.Opponent(card.Player)

		// Snapshot the hand before anything leaves it, and read power from the
		// snapshot too: a card in hand has no battle zone modifiers, but the
		// engine still owns that calculation.
		hand := fx.Find(opponent, match.HAND)

		if len(hand) < 1 {
			return
		}

		revealed := make([]string, 0, len(hand))
		for _, x := range hand {
			revealed = append(revealed, x.ImageID)
		}

		ctx.Match.ShowCards(card.Player, fmt.Sprintf("%s's effect: your opponent's hand", card.Name), revealed)

		discarded := 0
		for _, x := range hand {
			if !x.HasCondition(cnd.Creature) || ctx.Match.GetPower(x, false) > 3000 {
				continue
			}

			moved, err := opponent.MoveCard(x.ID, match.HAND, match.GRAVEYARD, card.ID)

			if err != nil || moved.Zone != match.GRAVEYARD {
				continue
			}

			discarded++
		}

		if discarded > 0 {
			ctx.Match.ReportActionInChat(opponent, fmt.Sprintf("%s discarded %d creature(s) with power 3000 or less because of %s", opponent.Username(), discarded, card.Name))
		}
	}

	c.Use(fx.Creature, fx.PutIntoManaZoneTapped, fx.TapAbility)

}
