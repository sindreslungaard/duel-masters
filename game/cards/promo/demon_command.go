package promo

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// OlgateNightmareSamurai ...
func OlgateNightmareSamurai(c *match.Card) {

	c.Name = "Olgate, Nightmare Samurai"
	c.Power = 6000
	c.Civ = civ.Darkness
	c.Family = []string{family.DemonCommand}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Doublebreaker, fx.When(fx.AnotherOwnCreatureDestroyed, fx.MayUntapSelf))

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
