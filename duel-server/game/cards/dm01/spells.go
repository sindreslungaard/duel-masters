package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// AuraBlast ...
func AuraBlast(c *match.Card) {

	c.Name = "Aura Blast"
	c.Civ = civ.Nature
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		fx.Find(
			card.Player,
			match.BATTLEZONE,
		).Map(func(x *match.Card) {
			x.AddCondition(cnd.PowerAttacker, 2000, card.ID)
			ctx.Match.ReportActionInChat(x.Player, fmt.Sprintf("%s was given \"Power Attacker +2000\"", x.Name))
		})
	}))

}

// BrainSerum ...
func BrainSerum(c *match.Card) {

	c.Name = "Brain Serum"
	c.Civ = civ.Water
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, fx.DrawUpTo2))

}

// BurningPower ...
func BurningPower(c *match.Card) {

	c.Name = "Burning Power"
	c.Civ = civ.Fire
	c.ManaCost = 1
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.BATTLEZONE,
			fmt.Sprintf("%s: Select 1 creature from your battlezone that will gain \"Power Attacker +2000\"", card.Name),
			1,
			1,
			false,
		).Map(func(x *match.Card) {
			x.AddCondition(cnd.PowerAttacker, 2000, card.ID)
			ctx.Match.ReportActionInChat(x.Player, fmt.Sprintf("%s was given power attacker +2000", x.Name))
		})
	}))

}

// ChaosStrike ...
func ChaosStrike(c *match.Card) {

	c.Name = "Chaos Strike"
	c.Civ = civ.Fire
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		fx.SelectFilter(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			fmt.Sprintf("%s: Select 1 of your opponent's untapped creatures in the battle zone. Your creatures can attack it this turn as though it were tapped.", card.Name),
			1,
			1,
			false,
			func(x *match.Card) bool {
				return !x.Tapped
			},
			false,
		).Map(func(x *match.Card) {
			x.AddUniqueSourceCondition(cnd.TreatedAsTapped, nil, card.ID)
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was given \"Treated as Tapped\" condition by %s. Now your creatures can attack it as it were tapped during this turn.", x.Name, card.Name))
			ctx.Match.ReportActionInChat(x.Player, fmt.Sprintf("%s will be attackable by your opponent as it were tapped during this turn", x.Name))
		})
	}))

}

// CreepingPlague ...
func CreepingPlague(c *match.Card) {

	c.Name = "Creeping Plague"
	c.Civ = civ.Darkness
	c.ManaCost = 1
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {

			if event, ok := ctx2.Event.(*match.Battle); ok {

				if !event.Blocked || event.Attacker.Player != card.Player {
					return
				}

				event.Attacker.AddCondition(cnd.Slayer, nil, card.ID)
				ctx2.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was given \"Slayer\" by %s", event.Attacker.Name, card.Name))

			}

			// remove persistent effect when turn ends
			_, ok := ctx2.Event.(*match.EndOfTurnStep)
			if ok {
				exit()
			}

		})

	}))

}

// CrimsonHammer ...
func CrimsonHammer(c *match.Card) {

	c.Name = "Crimson Hammer"
	c.Civ = civ.Fire
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, fx.When(fx.SpellCast, fx.DestroyBySpellOpCreature2000OrLess))

}

// CrystalMemory ...
func CrystalMemory(c *match.Card) {

	c.Name = "Crystal Memory"
	c.Civ = civ.Water
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		selectedCards := fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.DECK,
			fmt.Sprintf("%s effect: Select 1 card from your deck that will be sent to your hand", card.Name),
			1,
			1,
			true,
		)

		for _, selectedCard := range selectedCards {
			card.Player.MoveCard(selectedCard.ID, match.DECK, match.HAND, card.ID)
			ctx.Match.ReportActionInChat(card.Player, card.Player.Username()+" retrieved a card from their deck")
		}

		fx.ShuffleDeck(card, ctx, false)

	}))
}

// DarkReversal ...
func DarkReversal(c *match.Card) {

	c.Name = "Dark Reversal"
	c.Civ = civ.Darkness
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		fx.SelectFilter(
			card.Player,
			ctx.Match,
			card.Player,
			match.GRAVEYARD,
			fmt.Sprintf("%s: Select 1 creature from your graveyard that will be sent to your hand", card.Name),
			1,
			1,
			false,
			func(x *match.Card) bool { return x.HasCondition(cnd.Creature) },
			false,
		).Map(func(x *match.Card) {
			x.Player.MoveCard(x.ID, match.GRAVEYARD, match.HAND, card.ID)
			ctx.Match.ReportActionInChat(x.Player, fmt.Sprintf("%s retrieved %s from their graveyard", x.Player.Username(), x.Name))
		})
	}))

}

// DeathSmoke ...
func DeathSmoke(c *match.Card) {

	c.Name = "Death Smoke"
	c.Civ = civ.Darkness
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		fx.SelectFilter(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			fmt.Sprintf("%s: Destroy one of your opponent's untapped creatures", card.Name),
			1,
			1,
			false,
			func(x *match.Card) bool { return x.Tapped == false },
			false,
		).Map(func(x *match.Card) {
			ctx.Match.Destroy(x, card, match.DestroyedBySpell)
		})
	}))

}

// DimensionGate ...
func DimensionGate(c *match.Card) {

	c.Name = "Dimension Gate"
	c.Civ = civ.Nature
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, fx.SearchDeckTake1Creature))

}

// GhostTouch ...
func GhostTouch(c *match.Card) {

	c.Name = "Ghost Touch"
	c.Civ = civ.Darkness
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, fx.OpponentDiscardsRandomCard))

}

// HolyAwe ...
func HolyAwe(c *match.Card) {

	c.Name = "Holy Awe"
	c.Civ = civ.Light
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		fx.Find(
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
		).Map(func(x *match.Card) {
			x.Tapped = true
			ctx.Match.ReportActionInChat(x.Player, fmt.Sprintf("%s was tapped by Holy Awe", x.Name))
		})
	}))

}

// LaserWing ...
func LaserWing(c *match.Card) {

	c.Name = "Laser Wing"
	c.Civ = civ.Light
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.BATTLEZONE,
			fmt.Sprintf("%s: Select up to 2 creatures that can't be blocked this turn", card.Name),
			1,
			2,
			true,
		).Map(func(x *match.Card) {
			x.AddCondition(cnd.CantBeBlocked, nil, card.ID)
			ctx.Match.ReportActionInChat(x.Player, x.Name+" can't be blocked this turn")
		})
	}))

}

// MagmaGazer ...
func MagmaGazer(c *match.Card) {

	c.Name = "Magma Gazer"
	c.Civ = civ.Fire
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.BATTLEZONE,
			"Select 1 creature from your battlezone that will gain \"Power Attacker +4000\" and \"Double breaker\"",
			1,
			1,
			false,
		).Map(func(x *match.Card) {
			x.AddCondition(cnd.PowerAttacker, 4000, card.ID)
			x.AddCondition(cnd.DoubleBreaker, nil, card.ID)
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was given \"Power Attacker +4000\" and \"Double Breaker\" until the end of the turn", x.Name))
		})
	}))

}

// MoonlightFlash ...
func MoonlightFlash(c *match.Card) {

	c.Name = "Moonlight Flash"
	c.Civ = civ.Light
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		fx.Select(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			fmt.Sprintf("%s: Select up to 2 of your opponents creatures that will be tapped", card.Name),
			1,
			2,
			true,
		).Map(func(x *match.Card) {
			x.Tapped = true
			ctx.Match.ReportActionInChat(x.Player, x.Name+" was tapped")
		})
	}))

}

// NaturalSnare ...
func NaturalSnare(c *match.Card) {

	c.Name = "Natural Snare"
	c.Civ = civ.Nature
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		fx.Select(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			fmt.Sprintf("%s: Select 1 of your opponent's creatures and put it in their manazone", card.Name),
			1,
			1,
			false,
		).Map(func(x *match.Card) {
			x.Player.MoveCard(x.ID, match.BATTLEZONE, match.MANAZONE, card.ID)
			x.Tapped = false
			ctx.Match.ReportActionInChat(x.Player, fmt.Sprintf("%s was moved to %s's manazone", x.Name, x.Player.Username()))
		})
	}))

}

// PangaeasSong ...
func PangaeasSong(c *match.Card) {

	c.Name = "Pangaea's Song"
	c.Civ = civ.Nature
	c.ManaCost = 1
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Spell, fx.When(fx.SpellCast, fx.PutOwnCreatureFromBZToMZ))
}

// SolarRay ...
func SolarRay(c *match.Card) {

	c.Name = "Solar Ray"
	c.Civ = civ.Light
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, fx.TapOpCreature))

}

// SonicWing ...
func SonicWing(c *match.Card) {

	c.Name = "Sonic Wing"
	c.Civ = civ.Light
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.When(fx.SpellCast, fx.GiveOwnCreatureCantBeBlocked))

}

// SpiralGate ...
func SpiralGate(c *match.Card) {

	c.Name = "Spiral Gate"
	c.Civ = civ.Water
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, fx.ReturnCreatureToOwnersHand))

}

// Teleportation ...
func Teleportation(c *match.Card) {

	c.Name = "Teleportation"
	c.Civ = civ.Water
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		cards := make(map[string][]*match.Card)

		myCards, err := card.Player.Container(match.BATTLEZONE)

		if err != nil {
			return
		}

		opponentCards, err := ctx.Match.Opponent(card.Player).Container(match.BATTLEZONE)

		if err != nil {
			return
		}

		if len(myCards) < 1 && len(opponentCards) < 1 {
			return
		}

		cards["Your creatures"] = myCards
		cards["Opponent's creatures"] = opponentCards

		fx.SelectMultipart(
			card.Player,
			ctx.Match,
			cards,
			fmt.Sprintf("%s: Select up to 2 creatures in the battlezone and return it to its owner's hand", card.Name),
			1,
			2,
			true,
		).Map(func(x *match.Card) {
			ref, err := c.Player.MoveCard(x.ID, match.BATTLEZONE, match.HAND, card.ID)

			if err != nil {
				ref, err := ctx.Match.Opponent(c.Player).MoveCard(x.ID, match.BATTLEZONE, match.HAND, card.ID)

				if err == nil {
					ctx.Match.ReportActionInChat(ctx.Match.Opponent(card.Player), fmt.Sprintf("%s was moved to %s's hand", ref.Name, ctx.Match.PlayerRef(ref.Player).Socket.User.Username))
				}
			} else {
				ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was moved to %s's hand", ref.Name, ctx.Match.PlayerRef(ref.Player).Socket.User.Username))
			}
		})
	}))

}

// TerrorPit ...
func TerrorPit(c *match.Card) {

	c.Name = "Terror Pit"
	c.Civ = civ.Darkness
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, fx.DestroyOpCreature))
}

// TornadoFlame ...
func TornadoFlame(c *match.Card) {

	c.Name = "Tornado Flame"
	c.Civ = civ.Fire
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		fx.SelectFilter(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			fmt.Sprintf("%s: Destroy one of your opponent's creatures that has power 4000 or less", card.Name),
			1,
			1,
			false,
			func(x *match.Card) bool { return ctx.Match.GetPower(x, false) <= 4000 },
			false,
		).Map(func(x *match.Card) {
			ctx.Match.Destroy(x, card, match.DestroyedBySpell)
		})
	}))

}

// UltimateForce ...
func UltimateForce(c *match.Card) {

	c.Name = "Ultimate Force"
	c.Civ = civ.Nature
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Spell, fx.When(fx.SpellCast, fx.Draw2ToMana))

}

// VirtualTripwire ...
func VirtualTripwire(c *match.Card) {

	c.Name = "Virtual Tripwire"
	c.Civ = civ.Water
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Spell, fx.When(fx.SpellCast, fx.TapOpCreature))

}
