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
