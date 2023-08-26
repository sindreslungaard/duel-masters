package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"sort"
)

// HypersquidWalter ...
func HypersquidWalter(c *match.Card) {

	c.Name = "Hypersquid Walter"
	c.Power = 1000
	c.Civ = civ.Water
	c.Family = []string{family.CyberLord}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.When(fx.Attacking, func(card *match.Card, ctx *match.Context) {

		ctx.ScheduleAfter(func() {
			card.Player.DrawCards(1)
		})

	}))

}

// Corile ...
func Corile(c *match.Card) {

	c.Name = "Corile"
	c.Power = 2000
	c.Civ = civ.Water
	c.Family = []string{family.CyberLord}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		// NOTE:
		// When moving an evolution card, the attached cards usually follow
		// but since we want to move it to the front of the deck
		// and let the opponent choose the order, we need to do some stuff manually
		// hence all the extra stuff below

		var toMove fx.CardCollection = []*match.Card{}

		// find attached cards recursively in case cards evolved from other evolution cards
		var addRecursive func(*match.Card)
		addRecursive = func(x *match.Card) {
			toMove = append(toMove, x)

			for _, attached := range x.Attachments() {
				addRecursive(attached)
			}
		}

		fx.Select(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			"Corile: Move 1 of your opponent's creatures from their battlezone to the top of their deck",
			1,
			1,
			false,
		).Map(func(x *match.Card) {
			// temporarily move this card to the hidden zone
			x.Player.MoveCard(x.ID, match.BATTLEZONE, match.HIDDENZONE)

			// finds all attached cards recursively and adds references of them to toMove[]
			addRecursive(x)
		})

		// let the opponent select which card to add at the very top if more than 1
		if len(toMove) > 1 {

			ctx.Match.Wait(card.Player, "Waiting for your opponent to make an action...")

			fx.SelectFilter(
				ctx.Match.Opponent(card.Player),
				ctx.Match,
				ctx.Match.Opponent(card.Player),
				match.HIDDENZONE,
				"Corile: These cards will be moved to the top of your deck. Select which one you want to be at the very top.",
				1,
				1,
				false,
				func(x *match.Card) bool {
					for _, y := range toMove {
						if x == y {
							return true
						}
					}
					return false
				},
			).Map(func(x *match.Card) {
				// sort toMove so that this card is added to the very top of the deck
				sort.SliceStable(toMove, func(i, j int) bool {
					return toMove[i] != x
				})
			})

			ctx.Match.EndWait(card.Player)
		}

		toMove.Map(func(x *match.Card) {
			ctx.Match.MoveCardToFront(x, match.DECK, card)
		})

	}))

}
