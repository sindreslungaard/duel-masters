package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// BurstShot ...
func BurstShot(c *match.Card) {

	c.Name = "Burst Shot"
	c.Civ = civ.Fire
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, fx.ShieldTrigger, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			opponent := ctx.Match.Opponent(card.Player)

			myCreatures, err := card.Player.Container(match.BATTLEZONE)
			if err != nil {
				return
			}

			opponentCreatures, err := opponent.Container(match.BATTLEZONE)
			if err != nil {
				return
			}

			toDestroy := make([]*match.Card, 0)

			for _, creature := range myCreatures {
				if ctx.Match.GetPower(creature, false) <= 2000 {
					toDestroy = append(toDestroy, creature)
				}
			}

			for _, creature := range opponentCreatures {
				if ctx.Match.GetPower(creature, false) <= 2000 {
					toDestroy = append(toDestroy, creature)
				}
			}

			for _, creature := range toDestroy {
				ctx.Match.Destroy(creature, card, match.DestroyedBySpell)
			}

		}

	})

}

// LogicCube ...
func LogicCube(c *match.Card) {

	c.Name = "Logic Cube"
	c.Civ = civ.Light
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		creatures := fx.SelectFilterFullList(
			card.Player,
			ctx.Match,
			card.Player,
			match.DECK,
			"Select 1 spell from your deck that will be shown to your opponent and sent to your hand",
			1,
			1,
			true,
			func(x *match.Card) bool { return x.HasCondition(cnd.Spell) },
			true,
		)

		for _, creature := range creatures {

			card.Player.MoveCard(creature.ID, match.DECK, match.HAND, card.ID)
			ctx.Match.Chat("Server", fmt.Sprintf("%s retrieved %s from the deck to their hand", card.Player.Username(), creature.Name))

		}

		card.Player.ShuffleDeck()

	}))

}

// ThoughtProbe ...
func ThoughtProbe(c *match.Card) {

	c.Name = "Thought Probe"
	c.Civ = civ.Water
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		creatures := fx.Find(
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
		)

		if len(creatures) >= 3 {
			card.Player.DrawCards(3)
		}

	}))

}

// CriticalBlade ...
func CriticalBlade(c *match.Card) {

	c.Name = "Critical Blade"
	c.Civ = civ.Darkness
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		fx.SelectFilter(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			"Critical Blade: Select 1 of your opponent's blockers that will be destroyed",
			1,
			1,
			false,
			func(x *match.Card) bool { return x.HasCondition(cnd.Blocker) },
		).Map(func(x *match.Card) {
			ctx.Match.Destroy(x, card, match.DestroyedBySpell)
		})

	}))

}

// ReconOperation ...
func ReconOperation(c *match.Card) {

	c.Name = "Recon Operation"
	c.Civ = civ.Water
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		shields := fx.SelectBackside(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.SHIELDZONE,
			"Recon Operation: Select 3 of your opponent's shields that will be shown to you",
			1,
			3,
			true,
		)

		ids := make([]string, 0)

		for _, s := range shields {
			ids = append(ids, s.ImageID)
		}

		ctx.Match.ShowCards(
			card.Player,
			"Your opponent's shields:",
			ids,
		)

	}))

}

// ManaCrisis ...
func ManaCrisis(c *match.Card) {

	c.Name = "Mana Crisis"
	c.Civ = civ.Nature
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		fx.Select(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.MANAZONE,
			"Mana Crisis: Select 1 card from your opponent's manazone and put it in their graveyard",
			1,
			1,
			false,
		).Map(func(x *match.Card) {
			x.Player.MoveCard(x.ID, match.MANAZONE, match.GRAVEYARD, card.ID)
			ctx.Match.Chat("Server", fmt.Sprintf("%s was put into %s's graveyard from their manazone by %s", x.Name, x.Player.Username(), card.Name))
		})

	}))

}

// RumbleGate ...
func RumbleGate(c *match.Card) {

	c.Name = "Rumble Gate"
	c.Civ = civ.Fire
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		fx.Find(
			card.Player,
			match.BATTLEZONE,
		).Map(func(x *match.Card) {
			x.AddCondition(cnd.PowerAmplifier, 1000, card.ID)
			x.AddCondition(cnd.AttackUntapped, nil, card.ID)
			ctx.Match.Chat("Server", fmt.Sprintf("%s can attack untapped creatures and was given +1000 power until the end of this turn by %s", x.Name, card.Name))
		})

	}))

}

// LostSoul ...
func LostSoul(c *match.Card) {

	c.Name = "Lost Soul"
	c.Civ = civ.Darkness
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		fx.Find(
			ctx.Match.Opponent(card.Player),
			match.HAND,
		).Map(func(x *match.Card) {
			ctx.Match.Opponent(card.Player).MoveCard(x.ID, match.HAND, match.GRAVEYARD, card.ID)
		})

		ctx.Match.Chat("Server", fmt.Sprintf("%s's hand was discarded by %s", ctx.Match.Opponent(card.Player).Username(), card.Name))

	}))

}

// RainbowStone ...
func RainbowStone(c *match.Card) {

	c.Name = "Rainbow Stone"
	c.Civ = civ.Nature
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.DECK,
			"Rainbow Stone: Put a card from your deck into your manazone",
			1,
			1,
			true,
		).Map(func(x *match.Card) {
			x.Player.MoveCard(x.ID, match.DECK, match.MANAZONE, card.ID)
			ctx.Match.Chat("Server", fmt.Sprintf("%s put %s in their manazone from their deck", x.Player.Username(), x.Name))
		})

		card.Player.ShuffleDeck()

	}))

}

// DiamondCutter ...
func DiamondCutter(c *match.Card) {

	c.Name = "Diamond Cutter"
	c.Civ = civ.Light
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {

			fx.Find(
				card.Player,
				match.BATTLEZONE,
			).Map(func(c *match.Card) {

				if c.HasCondition(cnd.SummoningSickness) {
					c.RemoveCondition(cnd.SummoningSickness)
					c.AddCondition(cnd.CantAttackCreatures, nil, card.ID)
				}

				c.RemoveCondition(cnd.CantAttackPlayers)

				if event, ok := ctx2.Event.(*match.AttackCreature); ok {

					// Is this event for me or someone else?
					if event.CardID != c.ID || !c.HasCondition(cnd.CantAttackCreatures) {
						return
					}

					ctx2.Match.WarnPlayer(c.Player, fmt.Sprintf("%s can't attack creatures", c.Name))

					ctx2.InterruptFlow()

				}
			})

			// remove persistent effect when turn ends
			_, ok := ctx2.Event.(*match.EndStep)
			if ok {
				exit()
			}

		})

	}))

}
