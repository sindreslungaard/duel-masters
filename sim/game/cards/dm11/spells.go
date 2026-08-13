package dm11

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// SolarTrap ...
func SolarTrap(c *match.Card) {

	c.Name = "Solar Trap"
	c.Civs = []string{civ.Light}
	c.ManaCost = 1
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.When(fx.SpellCast, fx.TapOpCreature))

}

// TenTonCrunch ...
func TenTonCrunch(c *match.Card) {

	c.Name = "Ten-Ton Crunch"
	c.Civs = []string{civ.Fire}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, fx.DestroyOpCreatureXPowerOrLess(3000, false, match.DestroyedBySpell)))

}

// MorbidMedicine ...
func MorbidMedicine(c *match.Card) {

	c.Name = "Morbid Medicine"
	c.Civs = []string{civ.Darkness}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Spell, fx.When(fx.SpellCast, fx.ReturnXCreaturesFromGraveToHand(2)))

}

// EmergencyTyphoon ...
func EmergencyTyphoon(c *match.Card) {

	c.Name = "Emergency Typhoon"
	c.Civs = []string{civ.Water}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		fx.DrawUpTo2(card, ctx)

		// The discard is not conditional on the draw: the printed text says to
		// draw up to 2 and then discard, so drawing none still costs a card.
		fx.DiscardOwnXCards(1)(card, ctx)
	}))

}

// RainbowGate ...
func RainbowGate(c *match.Card) {

	c.Name = "Rainbow Gate"
	c.Civs = []string{civ.Nature}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		fx.SearchDeckTakeCards(
			card,
			ctx,
			1,
			func(x *match.Card) bool { return x.HasCondition(cnd.Creature) && x.IsMulticolored() },
			"multi-colored creature",
		)
	}))

}

// MiraculousSnare ...
func MiraculousSnare(c *match.Card) {

	c.Name = "Miraculous Snare"
	c.Civs = []string{civ.Light, civ.Water}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light, civ.Water}

	c.Use(fx.Spell, fx.PutIntoManaZoneTapped, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		nonEvolution := func(x *match.Card) bool { return !x.HasCondition(cnd.Evolution) }

		// Either side of the board is fair game, so the caster's own creatures
		// are offered alongside the opponent's.
		cards := map[string][]*match.Card{
			"Your creatures":            fx.FindFilter(card.Player, match.BATTLEZONE, nonEvolution),
			"Your opponent's creatures": fx.FindFilter(ctx.Match.Opponent(card.Player), match.BATTLEZONE, nonEvolution),
		}

		fx.SelectMultipart(
			card.Player,
			ctx.Match,
			cards,
			fmt.Sprintf("%s's effect: Choose a non-evolution creature in the battle zone and add it to its owner's shields face down.", card.Name),
			1,
			1,
			false,
		).Map(func(creature *match.Card) {
			moved, err := creature.Player.MoveCard(creature.ID, match.BATTLEZONE, match.SHIELDZONE, card.ID)

			if err != nil || moved.Zone != match.SHIELDZONE {
				return
			}

			ctx.Match.ReportActionInChat(creature.Player, fmt.Sprintf("%s was added to %s's shields by %s", creature.Name, creature.Player.Username(), card.Name))
		})
	}))

}

// HideAndSeek ...
func HideAndSeek(c *match.Card) {

	c.Name = "Hide and Seek"
	c.Civs = []string{civ.Water, civ.Darkness}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Water, civ.Darkness}

	c.Use(fx.Spell, fx.PutIntoManaZoneTapped, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		opponent := ctx.Match.Opponent(card.Player)

		fx.SelectFilter(
			card.Player,
			ctx.Match,
			opponent,
			match.BATTLEZONE,
			fmt.Sprintf("%s's effect: Choose one of your opponent's non-evolution creatures and return it to their hand.", card.Name),
			1,
			1,
			false,
			func(x *match.Card) bool { return !x.HasCondition(cnd.Evolution) },
			false,
		).Map(func(creature *match.Card) {
			moved, err := opponent.MoveCard(creature.ID, match.BATTLEZONE, match.HAND, card.ID)

			if err != nil || moved.Zone != match.HAND {
				return
			}

			ctx.Match.ReportActionInChat(opponent, fmt.Sprintf("%s was returned to %s's hand by %s", creature.Name, opponent.Username(), card.Name))
		})

		// The discard is a sentence of its own rather than a consequence of the
		// return, so it happens even when there was no creature to choose.
		fx.OpponentDiscardsRandomCard(card, ctx)
	}))

}

// ReapAndSow ...
func ReapAndSow(c *match.Card) {

	c.Name = "Reap and Sow"
	c.Civs = []string{civ.Fire, civ.Nature}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Fire, civ.Nature}

	c.Use(fx.Spell, fx.PutIntoManaZoneTapped, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		fx.ManaBurn(1, false)(card, ctx)
		fx.Draw1ToMana(card, ctx)
	}))

}

// RiseAndShine ...
func RiseAndShine(c *match.Card) {

	c.Name = "Rise and Shine"
	c.Civs = []string{civ.Light, civ.Water}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light, civ.Water}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.PutIntoManaZoneTapped, fx.When(fx.SpellCast,
		fx.RevealTopXTake1ReorderRestOnBottom(4, func(x *match.Card) bool {
			return x.HasCondition(cnd.Blocker)
		}, "card that has \"blocker\"")))

}

// RouletteOfRuin ...
func RouletteOfRuin(c *match.Card) {

	c.Name = "Roulette of Ruin"
	c.Civs = []string{civ.Darkness}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, fx.ChooseANumberAndDiscardByCost))

}
