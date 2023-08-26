package dm04

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// RimuelCloudbreakElemental ...
func RimuelCloudbreakElemental(c *match.Card) {

	c.Name = "Rimuel, Cloudbreak Elemental"
	c.Power = 6000
	c.Civ = civ.Light
	c.Family = []string{family.AngelCommand}
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Doublebreaker, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		lightMana := len(fx.FindFilter(
			card.Player,
			match.MANAZONE,
			func(x *match.Card) bool { return x.Civ == civ.Light && !x.Tapped },
		))

		nrCreaturesOpp := len(fx.FindFilter(
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			func(x *match.Card) bool { return !x.Tapped },
		))

		toSelect := lightMana
		if toSelect > nrCreaturesOpp {
			toSelect = nrCreaturesOpp
		}

		fx.Select(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			fmt.Sprintf("Rimuel, Cloudbreak Elemental: select %d of opponent's creatures and tap them", toSelect),
			toSelect,
			toSelect,
			false,
		).Map(func(x *match.Card) {
			x.Tapped = true
			ctx.Match.Chat("Server", fmt.Sprintf("%s's %s was tapped by Rimuel, Cloudbreak Elemental", x.Player.Username(), x.Name))
		})
	}))

}

// AlcadeiasLordOfSpirits ...
func AlcadeiasLordOfSpirits(c *match.Card) {

	c.Name = "Alcadeias, Lord of Spirits"
	c.Power = 12500
	c.Civ = civ.Light
	c.Family = []string{family.AngelCommand}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Evolution, fx.Doublebreaker, func(card *match.Card, ctx *match.Context) {

		if card.Zone != match.BATTLEZONE {
			return
		}

		if event, ok := ctx.Event.(*match.ShieldTriggerEvent); ok {

			if event.Card.Civ != civ.Light && event.Card.HasCondition(cnd.Spell) {
				ctx.InterruptFlow()
			}

		}

		if event, ok := ctx.Event.(*match.PlayCardEvent); ok {

			p := ctx.Match.CurrentPlayer()

			playedCard, err := p.Player.GetCard(event.CardID, match.HAND)

			if err != nil || !playedCard.HasCondition(cnd.Spell) {
				return
			}
			if playedCard.Civ != civ.Light {
				ctx.Match.WarnPlayer(ctx.Match.Opponent(card.Player), "Only light spells may be cast while Alcadeias, Lord of Spirits is in the battle zone")
				ctx.InterruptFlow()
			}
		}
	})

}

// AerisFlightElemental ...
func AerisFlightElemental(c *match.Card) {

	c.Name = "Aeris, Flight Elemental"
	c.Power = 9000
	c.Civ = civ.Light
	c.Family = []string{family.AngelCommand}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.CantAttackPlayers, func(card *match.Card, ctx *match.Context) {

		event, ok := ctx.Event.(*match.AttackCreature)

		if !ok || event.CardID != card.ID {
			return
		}

		fx.FindFilter(
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			func(x *match.Card) bool {
				return x.Civ == civ.Darkness && x.Tapped == false
			},
		).Map(func(x *match.Card) {
			// don't add if already in the list of attackable creatures
			for _, creature := range event.AttackableCreatures {
				if creature.ID == x.ID {
					return
				}
			}

			event.AttackableCreatures = append(event.AttackableCreatures, x)
		})

	})

}
