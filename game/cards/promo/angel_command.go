package promo

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// AmnisHolyElemental ...
func AmnisHolyElemental(c *match.Card) {

	c.Name = "Amnis, Holy Elemental"
	c.Power = 5000
	c.Civ = civ.Light
	c.Family = []string{family.AngelCommand}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Blocker, func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.AttackCreature); ok && !ctx.Match.IsPlayerTurn(card.Player) {
			attacker, err := ctx.Match.Opponent(card.Player).GetCard(event.CardID, match.BATTLEZONE)
			if err != nil {
				return
			}

			if attacker.Civ != civ.Darkness {
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

			if attacker.Civ != civ.Darkness {
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

			if event.Context == match.DestroyedInBattle && event.Source.Civ == civ.Darkness {
				ctx.InterruptFlow()
			}

		}

	})

}
