package dm12

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// ExtremeCrawler ...
func ExtremeCrawler(c *match.Card) {

	c.Name = "Extreme Crawler"
	c.Power = 7000
	c.Civs = []string{civ.Water}
	c.Family = []string{family.EarthEater}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Doublebreaker, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		// Snapshotted before anything moves, and "other" spares the crawler
		// that is doing the clearing.
		others := fx.FindFilter(card.Player, match.BATTLEZONE, func(x *match.Card) bool { return x.ID != card.ID })

		returned := 0
		for _, creature := range others {
			moved, err := card.Player.MoveCard(creature.ID, match.BATTLEZONE, match.HAND, card.ID)

			if err == nil && moved.Zone == match.HAND {
				returned++
			}
		}

		if returned > 0 {
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s returned %d of %s's creatures to their hand", card.Name, returned, card.Player.Username()))
		}
	}))

}

// TyphoonCrawler ...
func TyphoonCrawler(c *match.Card) {

	c.Name = "Typhoon Crawler"
	c.Power = 5000
	c.Civs = []string{civ.Water}
	c.Family = []string{family.EarthEater}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {
		event, ok := ctx.Event.(*match.AttackCreature)

		if !ok || card.Zone != match.BATTLEZONE {
			return
		}

		attacker, err := ctx.Match.Opponent(card.Player).GetCard(event.CardID, match.BATTLEZONE)

		if err != nil || !(attacker.HasCiv(civ.Fire) || attacker.HasCiv(civ.Nature)) {
			return
		}

		// fx.Creature fills the attackable list while handling this same event
		// and only prompts in a scheduled callback. This card belongs to the
		// non-active player, so its handlers run after the attacker's and the
		// list can still be narrowed here, the same way Bodacious Giant widens
		// it. Narrowing from a scheduled callback would be too late, because
		// the prompt is scheduled first.
		remaining := make([]*match.Card, 0, len(event.AttackableCreatures))

		for _, creature := range event.AttackableCreatures {
			if creature.ID != card.ID {
				remaining = append(remaining, creature)
			}
		}

		event.AttackableCreatures = remaining
	})

}

// TropicCrawler ...
func TropicCrawler(c *match.Card) {

	c.Name = "Tropic Crawler"
	c.Power = 3000
	c.Civs = []string{civ.Water}
	c.Family = []string{family.EarthEater}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Blocker(), fx.CantAttackCreatures, fx.CantAttackPlayers,
		fx.When(fx.Blocks, fx.OpponentChoosesOwnCreatureToHand))

}
