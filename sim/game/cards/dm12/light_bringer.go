package dm12

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// MizoyTheOracle ...
func MizoyTheOracle(c *match.Card) {

	c.Name = "Mizoy, the Oracle"
	c.Power = 2500
	c.Civs = []string{civ.Light}
	c.Family = []string{family.LightBringer}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		matching := func(x *match.Card) bool { return x.HasCiv(civ.Darkness) || x.HasCiv(civ.Fire) }

		// "In the battle zone" is both sides of it.
		cards := map[string][]*match.Card{
			"Your creatures":            fx.FindFilter(card.Player, match.BATTLEZONE, matching),
			"Your opponent's creatures": fx.FindFilter(ctx.Match.Opponent(card.Player), match.BATTLEZONE, matching),
		}

		fx.SelectMultipart(
			card.Player,
			ctx.Match,
			cards,
			fmt.Sprintf("%s's effect: You may choose a darkness or fire creature in the battle zone and tap it.", card.Name),
			1,
			1,
			true,
		).Map(func(chosen *match.Card) {
			if card.Player.TapCard(chosen) {
				ctx.Match.ReportActionInChat(chosen.Player, fmt.Sprintf("%s was tapped by %s", chosen.Name, card.Name))
			}
		})
	}))

}

// PharziTheOracle ...
func PharziTheOracle(c *match.Card) {

	c.Name = "Pharzi, the Oracle"
	c.Power = 1000
	c.Civs = []string{civ.Light}
	c.Family = []string{family.LightBringer}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Light}

	// Destroyed is the post-destruction observation, so the graveyard it looks
	// through already holds this creature. A spell is what it wants, though, so
	// it can never fish itself back out. Unselectable cards are left out of the
	// offer so a graveyard without a single spell opens no prompt at all.
	c.Use(fx.Creature, fx.When(fx.Destroyed, func(card *match.Card, ctx *match.Context) {
		fx.SelectFilter(
			card.Player,
			ctx.Match,
			card.Player,
			match.GRAVEYARD,
			fmt.Sprintf("%s's effect: You may return a spell from your graveyard to your hand.", card.Name),
			1,
			1,
			true,
			func(x *match.Card) bool { return x.HasCondition(cnd.Spell) },
			false,
		).Map(func(spell *match.Card) {
			moved, err := card.Player.MoveCard(spell.ID, match.GRAVEYARD, match.HAND, card.ID)

			if err == nil && moved.Zone == match.HAND {
				ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was returned to %s's hand by %s", spell.Name, card.Player.Username(), card.Name))
			}
		})
	}))

}
