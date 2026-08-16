package dm12

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// maxChosenPower is the ceiling Mechadragon's Breath puts on its own number.
const maxChosenPower = 6000

// EnigmaticCascade ...
func EnigmaticCascade(c *match.Card) {

	c.Name = "Enigmatic Cascade"
	c.Civs = []string{civ.Water}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		// A spell is still in hand while it resolves, so it has to be kept out
		// of its own offer.
		notItself := func(x *match.Card) bool { return x.ID != card.ID }

		others := fx.FindFilter(card.Player, match.HAND, notItself)

		if len(others) < 1 {
			return
		}

		// "Any number" includes none, so the whole thing is declinable.
		discarded := fx.SelectFilter(
			card.Player,
			ctx.Match,
			card.Player,
			match.HAND,
			fmt.Sprintf("%s's effect: Discard any number of cards from your hand. You will draw that many.", card.Name),
			1,
			len(others),
			true,
			notItself,
			false,
		)

		drawn := 0
		for _, x := range discarded {
			moved, err := card.Player.MoveCard(x.ID, match.HAND, match.GRAVEYARD, card.ID)

			if err != nil || moved.Zone != match.GRAVEYARD {
				continue
			}

			drawn++
		}

		if drawn < 1 {
			return
		}

		card.Player.DrawCards(drawn)
		ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s discarded %d card(s) and drew that many with %s", card.Player.Username(), drawn, card.Name))
	}))

}

// MechadragonsBreath ...
func MechadragonsBreath(c *match.Card) {

	c.Name = "Mechadragon's Breath"
	c.Civs = []string{civ.Fire}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		power := fx.SelectCount(
			card.Player,
			ctx.Match,
			fmt.Sprintf("%s's effect: Choose a number up to %d. Every creature with exactly that power is destroyed.", card.Name, maxChosenPower),
			0,
			maxChosenPower,
		)

		fx.DestroyAllCreaturesWithExactPower(power, match.DestroyedBySpell)(card, ctx)
	}))

}

// clonedBladeMaxPower is the ceiling Cloned Blade puts on what it can destroy.
const clonedBladeMaxPower = 3000

// clonedTargetsFor works out the selection bounds the "Cloned" cycle offers: one
// target is mandatory and every copy of the spell lying in either graveyard adds
// one more that the caster may take.
func clonedTargetsFor(card *match.Card, ctx *match.Context) (int, int) {
	return 1, 1 + fx.ClonedCopiesInGraveyards(card, ctx.Match)
}

// ClonedBlade ...
func ClonedBlade(c *match.Card) {

	c.Name = "Cloned Blade"
	c.Civs = []string{civ.Fire}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		min, max := clonedTargetsFor(card, ctx)

		fx.SelectFilter(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			fmt.Sprintf("%s's effect: Destroy one of your opponent's creatures that has power %d or less, plus one more for each %s in each graveyard.", card.Name, clonedBladeMaxPower, card.Name),
			min,
			max,
			false,
			func(x *match.Card) bool { return ctx.Match.GetPower(x, false) <= clonedBladeMaxPower },
			false,
		).Map(func(creature *match.Card) {
			ctx.Match.Destroy(creature, card, match.DestroyedBySpell)
		})
	}))

}

// ClonedDeflector ...
func ClonedDeflector(c *match.Card) {

	c.Name = "Cloned Deflector"
	c.Civs = []string{civ.Light}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		min, max := clonedTargetsFor(card, ctx)

		fx.Select(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			fmt.Sprintf("%s's effect: Tap one of your opponent's creatures, plus one more for each %s in each graveyard.", card.Name, card.Name),
			min,
			max,
			false,
		).Map(func(creature *match.Card) {
			fx.TapCreature(creature, ctx, card)
		})
	}))

}

// ClonedNightmare ...
func ClonedNightmare(c *match.Card) {

	c.Name = "Cloned Nightmare"
	c.Civs = []string{civ.Darkness}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		hand, err := ctx.Match.Opponent(card.Player).Container(match.HAND)

		if err != nil || len(hand) < 1 {
			return
		}

		_, max := clonedTargetsFor(card, ctx)

		// The cards are taken at random rather than chosen, so unlike the rest
		// of the cycle the choice on offer is how many to take, not which.
		discards := 1
		if max > 1 {
			discards += fx.SelectCount(
				card.Player,
				ctx.Match,
				fmt.Sprintf("%s's effect: One card is discarded at random from your opponent's hand. You may take up to %d more, one for each %s in each graveyard.", card.Name, max-1, card.Name),
				0,
				max-1,
			)
		}

		for range discards {
			fx.OpponentDiscardsRandomCard(card, ctx)
		}
	}))

}

// CosmicDarts ...
func CosmicDarts(c *match.Card) {

	c.Name = "Cosmic Darts"
	c.Civs = []string{civ.Light}
	c.ManaCost = 1
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		shields := fx.SelectBackside(
			ctx.Match.Opponent(card.Player),
			ctx.Match,
			card.Player,
			match.SHIELDZONE,
			fmt.Sprintf("%s: Choose one of your opponent's shields", card.Name),
			1,
			1,
			false,
		)

		if len(shields) < 1 {
			return
		}

		shield := shields[0]

		if !shield.HasCondition(cnd.Spell) {
			ctx.Match.ShowCards(card.Player, fmt.Sprintf("%s: Your shield", card.Name), []string{shield.ImageID})
			return
		}

		// Non-dismissible so the caster actually sees the revealed spell before
		// being asked whether to cast it, instead of the preview being raced by
		// the immediately-following prompt.
		ctx.Match.ShowCardsNonDismissible(card.Player, fmt.Sprintf("%s: Your shield", card.Name), []string{shield.ImageID})

		if !fx.BinaryQuestion(card.Player, ctx.Match, fmt.Sprintf("%s: You may cast %s immediately for no cost", card.Name, shield.Name)) {
			return
		}

		moved, err := card.Player.MoveCard(shield.ID, match.SHIELDZONE, match.HAND, card.ID)

		if err != nil || moved.Zone != match.HAND {
			return
		}

		ctx.Match.CastSpell(moved, true)
	}))

}

// ClonedSpiral ...
func ClonedSpiral(c *match.Card) {

	c.Name = "Cloned Spiral"
	c.Civs = []string{civ.Water}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		min, max := clonedTargetsFor(card, ctx)

		// "A creature in the battle zone" with no owner named, so the caster's
		// own creatures are offered alongside the opponent's.
		creatures := map[string][]*match.Card{
			"Your creatures":            fx.Find(card.Player, match.BATTLEZONE),
			"Your opponent's creatures": fx.Find(ctx.Match.Opponent(card.Player), match.BATTLEZONE),
		}

		fx.SelectMultipart(
			card.Player,
			ctx.Match,
			creatures,
			fmt.Sprintf("%s's effect: Return a creature in the battle zone to its owner's hand, plus one more for each %s in each graveyard.", card.Name, card.Name),
			min,
			max,
			false,
		).Map(func(creature *match.Card) {
			moved, err := creature.Player.MoveCard(creature.ID, match.BATTLEZONE, match.HAND, card.ID)

			if err != nil || moved.Zone != match.HAND {
				return
			}

			ctx.Match.ReportActionInChat(creature.Player, fmt.Sprintf("%s was returned to %s's hand by %s", creature.Name, creature.Player.Username(), card.Name))
		})
	}))

}
