package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
	"math/rand"
)

// AuraBlast ...
func AuraBlast(c *match.Card) {

	c.Name = "Aura Blast"
	c.Civ = civ.Nature
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Spell, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			card.Player.MapContainer(match.BATTLEZONE, func(creature *match.Card) {
				creature.AddCondition(cnd.PowerAttacker, 2000, card.ID)
				ctx.Match.Chat("Server", fmt.Sprintf("%s was given power attacker +2000", creature.Name))
			})

		}

	})

}

// BrainSerum ...
func BrainSerum(c *match.Card) {

	c.Name = "Brain Serum"
	c.Civ = civ.Water
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Spell, fx.ShieldTrigger, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			card.Player.DrawCards(2)

		}

	})

}

// BurningPower ...
func BurningPower(c *match.Card) {

	c.Name = "Burning Power"
	c.Civ = civ.Fire
	c.ManaCost = 1
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			creatures := match.Search(card.Player, ctx.Match, card.Player, match.BATTLEZONE, "Select 1 creature from your battlezone that will gain \"Power Attacker +2000\"", 1, 1, false)

			for _, creature := range creatures {

				creature.AddCondition(cnd.PowerAttacker, 2000, card.ID)
				ctx.Match.Chat("Server", fmt.Sprintf("%s was given power attacker +2000", creature.Name))

			}

		}

	})

}

// ChaosStrike ...
func ChaosStrike(c *match.Card) {

	c.Name = "Chaos Strike"
	c.Civ = civ.Fire
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			creatures := match.Search(card.Player, ctx.Match, ctx.Match.Opponent(card.Player), match.BATTLEZONE, "Select 1 of your opponent's creatures that will be tapped", 1, 1, false)

			for _, creature := range creatures {

				creature.Tapped = true
				ctx.Match.Chat("Server", fmt.Sprintf("%s was tapped by %s", creature.Name, card.Name))

			}

		}

	})

}

// CreepingPlague ...
func CreepingPlague(c *match.Card) {

	c.Name = "Creeping Plague"
	c.Civ = civ.Darkness
	c.ManaCost = 1
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Spell, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {
			card.AddCondition(cnd.Active, nil, card.ID)
		}

		if event, ok := ctx.Event.(*match.Battle); ok {

			if !event.Blocked || !card.HasCondition(cnd.Active) {
				return
			}

			event.Attacker.AddCondition(cnd.Slayer, nil, card.ID)
			ctx.Match.Chat("Server", fmt.Sprintf("%s was given \"Slayer\" by %s", event.Attacker.Name, card.Name))

		}

	})

}

// CrimsonHammer ...
func CrimsonHammer(c *match.Card) {

	c.Name = "Crimson Hammer"
	c.Civ = civ.Fire
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			creatures := match.Filter(card.Player, ctx.Match, ctx.Match.Opponent(card.Player), match.BATTLEZONE, "Select 1 of your opponent's creatures that will be tapped", 1, 1, false, func(x *match.Card) bool { return x.Power <= 2000 })

			for _, creature := range creatures {

				ctx.Match.Destroy(creature, card)

			}

		}

	})

}

// CrystalMemory ...
func CrystalMemory(c *match.Card) {

	c.Name = "Crystal Memory"
	c.Civ = civ.Water
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Spell, fx.ShieldTrigger, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			selectedCards := match.Search(card.Player, ctx.Match, card.Player, match.DECK, "Select 1 card from your deck that will be sent to your hand", 1, 1, false)

			for _, selectedCard := range selectedCards {

				card.Player.MoveCard(selectedCard.ID, match.DECK, match.HAND)

			}

			card.Player.ShuffleDeck()

			ctx.Match.Chat("Server", card.Player.Username()+" retrieved a card from their deck")

		}

	})

}

// DarkReversal ...
func DarkReversal(c *match.Card) {

	c.Name = "Dark Reversal"
	c.Civ = civ.Darkness
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Spell, fx.ShieldTrigger, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			creatures := match.Filter(card.Player, ctx.Match, card.Player, match.GRAVEYARD, "Select 1 creature from your graveyard that will be sent to your hand", 1, 1, false, func(x *match.Card) bool { return x.HasCondition(cnd.Creature) })

			for _, creature := range creatures {

				card.Player.MoveCard(creature.ID, match.GRAVEYARD, match.HAND)
				ctx.Match.Chat("Server", fmt.Sprintf("%s retrieved %s from their graveyard", card.Player.Username(), creature.Name))

			}

		}

	})

}

// DeathSmoke ...
func DeathSmoke(c *match.Card) {

	c.Name = "Death Smoke"
	c.Civ = civ.Darkness
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Spell, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			creatures := match.Filter(card.Player, ctx.Match, ctx.Match.Opponent(card.Player), match.BATTLEZONE, "Destroy one of your opponent's untapped creatures", 1, 1, false, func(x *match.Card) bool { return x.Tapped == false })

			for _, creature := range creatures {

				ctx.Match.Destroy(creature, card)

			}

		}

	})

}

// DimensionGate ...
func DimensionGate(c *match.Card) {

	c.Name = "Dimension Gate"
	c.Civ = civ.Nature
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Spell, fx.ShieldTrigger, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			creatures := match.Filter(card.Player, ctx.Match, card.Player, match.DECK, "Select 1 creature from your deck that will be shown to your opponent and sent to your hand", 1, 1, false, func(x *match.Card) bool { return x.HasCondition(cnd.Creature) })

			for _, creature := range creatures {

				card.Player.MoveCard(creature.ID, match.DECK, match.HAND)
				ctx.Match.Chat("Server", fmt.Sprintf("%s retrieved %s from the deck to their hand", card.Player.Username(), creature.Name))

			}

			card.Player.ShuffleDeck()

		}

	})

}

// GhostTouch ...
func GhostTouch(c *match.Card) {

	c.Name = "Ghost Touch"
	c.Civ = civ.Darkness
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Spell, fx.ShieldTrigger, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			hand, err := ctx.Match.Opponent(card.Player).Container(match.HAND)

			if err != nil {
				return
			}

			if len(hand) < 1 {
				return
			}

			discardedCard, err := ctx.Match.Opponent(card.Player).MoveCard(hand[rand.Intn(len(hand))].ID, match.HAND, match.GRAVEYARD)
			if err == nil {
				ctx.Match.Chat("Server", fmt.Sprintf("%s was discarded from %s's hand", discardedCard.Name, discardedCard.Player.Username()))
			}

		}

	})

}

// HolyAwe ...
func HolyAwe(c *match.Card) {

	c.Name = "Holy Awe"
	c.Civ = civ.Light
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.ShieldTrigger, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			creatures, err := ctx.Match.Opponent(card.Player).Container(match.BATTLEZONE)

			if err != nil {
				return
			}

			for _, creature := range creatures {
				creature.Tapped = true
				ctx.Match.Chat("Server", fmt.Sprintf("%s was tapped by Holy Awe", creature.Name))
			}

		}

	})

}

// LaserWing ...
func LaserWing(c *match.Card) {

	c.Name = "Laser Wing"
	c.Civ = civ.Light
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			creatures := match.Search(card.Player, ctx.Match, card.Player, match.BATTLEZONE, "Select up to 2 creatures that can't be blocked this turn", 1, 2, false)

			for _, creature := range creatures {

				creature.AddCondition(cnd.CantBeBlocked, nil, card.ID)
				ctx.Match.Chat("Server", creature.Name+" can't be blocked this turn")

			}

		}

	})

}

// MagmaGazer ...
func MagmaGazer(c *match.Card) {

	c.Name = "Magma Gazer"
	c.Civ = civ.Fire
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			creatures := match.Search(card.Player, ctx.Match, card.Player, match.BATTLEZONE, "Select 1 creature from your battlezone that will gain \"Power Attacker +4000\" and \"Double breaker\"", 1, 1, false)

			for _, creature := range creatures {

				creature.AddCondition(cnd.PowerAttacker, 4000, card.ID)
				creature.AddCondition(cnd.DoubleBreaker, nil, card.ID)
				ctx.Match.Chat("Server", fmt.Sprintf("%s was given power attacker +4000 and double breaker until the end of the turn", creature.Name))

			}

		}

	})

}

// MoonlightFlash ...
func MoonlightFlash(c *match.Card) {

	c.Name = "Moonlight Flash"
	c.Civ = civ.Light
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			creatures := match.Search(card.Player, ctx.Match, ctx.Match.Opponent(card.Player), match.BATTLEZONE, "Select up to 2 of your opponents creatures that will be tapped", 1, 2, false)

			for _, creature := range creatures {

				creature.Tapped = true
				ctx.Match.Chat("Server", creature.Name+" was tapped")

			}

		}

	})

}

// NaturalSnare ...
func NaturalSnare(c *match.Card) {

	c.Name = "Natural Snare"
	c.Civ = civ.Nature
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Spell, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			creatures := match.Search(card.Player, ctx.Match, ctx.Match.Opponent(card.Player), match.BATTLEZONE, "Select 1 of your opponent's creatures and put it in their manazone", 1, 1, false)

			for _, creature := range creatures {

				creature.Player.MoveCard(creature.ID, match.BATTLEZONE, match.MANAZONE)
				creature.Tapped = false
				ctx.Match.Chat("Server", fmt.Sprintf("%s was moved to %s's manazone", creature.Name, creature.Player.Username()))

			}

		}

	})

}

// PangeasSong ...
func PangeasSong(c *match.Card) {

	c.Name = "Pangea's Song"
	c.Civ = civ.Nature
	c.ManaCost = 1
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Spell, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			creatures := match.Search(card.Player, ctx.Match, card.Player, match.BATTLEZONE, "Select 1 of your creatures and put it in your manazone", 1, 1, false)

			for _, creature := range creatures {

				creature.Player.MoveCard(creature.ID, match.BATTLEZONE, match.MANAZONE)
				creature.Tapped = false
				ctx.Match.Chat("Server", fmt.Sprintf("%s was moved to %s's manazone", creature.Name, creature.Player.Username()))

			}

		}

	})

}

// SolarRay ...
func SolarRay(c *match.Card) {

	c.Name = "SolarRay"
	c.Civ = civ.Light
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.ShieldTrigger, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			creatures := match.Search(card.Player, ctx.Match, ctx.Match.Opponent(card.Player), match.BATTLEZONE, "Select 1 of your opponents creatures that will be tapped", 1, 1, false)

			for _, creature := range creatures {

				creature.Tapped = true
				ctx.Match.Chat("Server", creature.Name+" was tapped")

			}

		}

	})

}
