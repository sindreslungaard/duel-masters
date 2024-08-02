package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"strings"
)

func CoccoLupia(c *match.Card) {

	c.Name = "Cocco Lupia"
	c.Power = 1000
	c.Civ = civ.Fire
	c.Family = []string{family.FireBird}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		event, ok := ctx.Event.(*match.PlayCardEvent)

		if !ok || !ctx.Match.IsPlayerTurn(card.Player) || card.Zone != match.BATTLEZONE {
			return
		}

		creature, err := card.Player.GetCard(event.CardID, match.HAND)

		if err != nil {
			return
		}

		ok = false
		for _, f := range creature.Family {
			if f == family.Dragonoid {
				return
			}

			if strings.Contains(strings.ToLower(f), "dragon") {
				ok = true
				break
			}
		}

		if !ok {
			return
		}

		manaCost := creature.ManaCost

		for _, condition := range creature.Conditions() {
			if condition.ID == cnd.ReducedCost {
				manaCost -= condition.Val.(int)
				if manaCost < 1 {
					manaCost = 1
				}
			}

			if condition.ID == cnd.IncreasedCost {
				manaCost += condition.Val.(int)
			}
		}

		if manaCost <= 2 {
			return
		}

		subtraction := 2
		if manaCost == 3 {
			subtraction = 1
		}

		creature.AddUniqueSourceCondition(cnd.ReducedCost, subtraction, card.ID)
	})

}
