package promo

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// OlgateNightmareSamurai ...
func OlgateNightmareSamurai(c *match.Card) {

	c.Name = "Olgate, Nightmare Samurai"
	c.Power = 6000
	c.Civ = civ.Darkness
	c.Family = []string{family.DemonCommand}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Doublebreaker, func(card *match.Card, ctx *match.Context) {

		if card.Zone != match.BATTLEZONE {
			return
		}

		if event, ok := ctx.Event.(*match.CardMoved); ok &&
			event.From == match.BATTLEZONE &&
			event.To == match.GRAVEYARD {

			// check if it was the card's player whose creature got destroyed
			var p *match.Player
			if event.MatchPlayerID == 1 {
				p = ctx.Match.Player1.Player
			} else {
				p = ctx.Match.Player2.Player
			}

			if card.Player == p && card.Tapped {

				if fx.BinaryQuestion(card.Player, ctx.Match, fmt.Sprintf("Do you want to untap %s?", card.Name)) {
					card.Tapped = false
				}

			}

		}

	})

}

// GiliamTheTormentor ...
func GiliamTheTormentor(c *match.Card) {

	c.Name = "Giliam, the Tormentor"
	c.Power = 5000
	c.Civ = civ.Darkness
	c.Family = []string{family.DemonCommand}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Blocker, func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.AttackCreature); ok && !ctx.Match.IsPlayerTurn(card.Player) {
			attacker, err := ctx.Match.Opponent(card.Player).GetCard(event.CardID, match.BATTLEZONE)
			if err != nil {
				return
			}

			if attacker.Civ != civ.Light {
				var newBlockersList []*match.Card
				for _, blocker := range event.Blockers {
					if blocker.ID != card.ID {
						newBlockersList = append(newBlockersList, blocker)
					}
				}
				event.Blockers = newBlockersList
			}
		}

		if event, ok := ctx.Event.(*match.AttackPlayer); ok && !ctx.Match.IsPlayerTurn(card.Player) {
			attacker, err := ctx.Match.Opponent(card.Player).GetCard(event.CardID, match.BATTLEZONE)
			if err != nil {
				return
			}

			if attacker.Civ != civ.Light {
				var newBlockersList []*match.Card
				for _, blocker := range event.Blockers {
					if blocker.ID != card.ID {
						newBlockersList = append(newBlockersList, blocker)
					}
				}
				event.Blockers = newBlockersList
			}
		}

		if event, ok := ctx.Event.(*match.CreatureDestroyed); ok && event.Card == card {

			if event.Context == match.DestroyedInBattle && event.Source.Civ == civ.Light {
				ctx.InterruptFlow()
			}

		}

	})

}
