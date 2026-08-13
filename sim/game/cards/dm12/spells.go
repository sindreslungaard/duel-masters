package dm12

import (
	"duel-masters/game/civ"
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
