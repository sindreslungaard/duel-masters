package dm11

import (
	"duel-masters/game/civ"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// LiveAndBreathe ...
func LiveAndBreathe(c *match.Card) {

	c.Name = "Live and Breathe"
	c.Civs = []string{civ.Light, civ.Nature}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light, civ.Nature}

	c.Use(fx.Spell, fx.PutIntoManaZoneTapped, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		// The fetched creature arrives in the battle zone as well, which would
		// look like another summon and search again forever. Only arrivals
		// from hand count as summoning, and the guard makes that certain even
		// if some other card moves a creature mid-search.
		searching := false

		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {

			if _, ok := ctx2.Event.(*match.EndOfTurnStep); ok {
				exit()
				return
			}

			event, ok := ctx2.Event.(*match.CardMoved)

			if !ok || searching || event.To != match.BATTLEZONE || event.From != match.HAND {
				return
			}

			summoned, err := card.Player.GetCard(event.CardID, match.BATTLEZONE)

			if err != nil {
				return
			}

			searching = true
			defer func() { searching = false }()

			fx.SearchDeckPutCreatureIntoBZ(
				card,
				ctx2,
				func(x *match.Card) bool { return x.Name == summoned.Name },
				fmt.Sprintf("another %s", summoned.Name),
			)
		})
	}))

}

// SlashAndBurn ...
func SlashAndBurn(c *match.Card) {

	c.Name = "Slash and Burn"
	c.Civs = []string{civ.Darkness, civ.Fire}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness, civ.Fire}

	c.Use(fx.Spell, fx.PutIntoManaZoneTapped, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {

			if _, ok := ctx2.Event.(*match.EndOfTurnStep); ok {
				exit()
				return
			}

			event, ok := ctx2.Event.(*match.CreatureDestroyed)

			if !ok || event.Card.Player == card.Player {
				return
			}

			ctx2.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s: %s was destroyed, so its controller loses a card from their mana zone and a shield", card.Name, event.Card.Name))

			fx.OpponentChoosesManaBurn(card, ctx2)
			fx.OpponentChoosesOwnShieldToGraveyard(card, ctx2)
		})
	}))

}
