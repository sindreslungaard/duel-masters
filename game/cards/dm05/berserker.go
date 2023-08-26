package dm05

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// BallusDogfightEnforcerQ ...
func BallusDogfightEnforcerQ(c *match.Card) {

	c.Name = "Ballus, Dogfight Enforcer Q"
	c.Power = 3000
	c.Civ = civ.Light
	c.Family = []string{family.Berserker}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Survivor, fx.When(fx.EndOfMyTurn, func(card *match.Card, ctx *match.Context) {

		if card.Zone != match.BATTLEZONE {
			return
		}

		fx.FindFilter(
			card.Player,
			match.BATTLEZONE,
			func(x *match.Card) bool { return x.HasCondition(cnd.Survivor) },
		).Map(func(x *match.Card) {
			x.Tapped = false
			ctx.Match.Chat("Server", fmt.Sprintf("%s was untapped by %s's survivor ability", x.Name, card.Name))
		})

	}))

}

// KulusSoulshineEnforcer ...
func KulusSoulshineEnforcer(c *match.Card) {

	c.Name = "Kulus, Soulshine Enforcer"
	c.Power = 3500
	c.Civ = civ.Light
	c.Family = []string{family.Berserker}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		if len(fx.Find(card.Player, match.MANAZONE)) < len(fx.Find(ctx.Match.Opponent(card.Player), match.MANAZONE)) {

			cards := card.Player.PeekDeck(1)

			for _, toMove := range cards {

				card.Player.MoveCard(toMove.ID, match.DECK, match.MANAZONE)
				ctx.Match.Chat("Server", fmt.Sprintf("%s put %s into the manazone from the top of their deck", card.Player.Username(), toMove.Name))

			}
		}
	}))

}
