package dm11

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// SquawkingLunatron ...
func SquawkingLunatron(c *match.Card) {

	c.Name = "Squawking Lunatron"
	c.Power = 4000
	c.Civs = []string{civ.Water}
	c.Family = []string{family.CyberMoon}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.SilentSkill(func(card *match.Card, ctx *match.Context) {
		fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.MANAZONE,
			fmt.Sprintf("%s's effect: Return up to 3 cards from your mana zone to your hand.", card.Name),
			1,
			3,
			true,
		).Map(func(x *match.Card) {
			moved, err := card.Player.MoveCard(x.ID, match.MANAZONE, match.HAND, card.ID)

			if err != nil || moved.Zone != match.HAND {
				return
			}

			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was returned from %s's mana zone to their hand by %s", x.Name, card.Player.Username(), card.Name))
		})
	}))

}

// WarpedLunatron ...
func WarpedLunatron(c *match.Card) {

	c.Name = "Warped Lunatron"
	c.Power = 6000
	c.Civs = []string{civ.Water}
	c.Family = []string{family.CyberMoon}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Doublebreaker, fx.When(fx.InTheBattlezone, func(card *match.Card, ctx *match.Context) {
		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {

			players := []*match.Player{card.Player, ctx2.Match.Opponent(card.Player)}

			if card.Zone != match.BATTLEZONE {
				for _, player := range players {
					fx.Find(player, match.BATTLEZONE).Map(func(x *match.Card) {
						x.RemoveSpecificConditionBySource(cnd.DoesntUntap, card.ID)
					})
				}

				exit()
				return
			}

			// "Creatures in the battle zone" is everybody's. Re-applied on
			// every event because persistent effects run before any card
			// handler, so the condition is always in place before fx.Creature
			// would untap anything.
			for _, player := range players {
				fx.Find(player, match.BATTLEZONE).Map(func(x *match.Card) {
					x.AddUniqueSourceCondition(cnd.DoesntUntap, true, card.ID)
				})
			}

			if _, ok := ctx2.Event.(*match.UntapStep); !ok {
				return
			}

			// The trade is offered to whoever the turn belongs to, not to this
			// creature's controller. Scheduled so it happens after the untap
			// step has finished refusing to untap anything.
			ctx2.ScheduleAfter(func() { warpedLunatronTrade(card, ctx2) })
		})
	}))

}

// warpedLunatronTrade lets the player whose turn it is tap mana to untap
// creatures, two cards of mana for each creature.
func warpedLunatronTrade(card *match.Card, ctx *match.Context) {
	player := ctx.Match.CurrentPlayer().Player

	tappedCreatures := fx.FindFilter(player, match.BATTLEZONE, func(x *match.Card) bool { return x.Tapped })
	untappedMana := fx.FindFilter(player, match.MANAZONE, func(x *match.Card) bool { return !x.Tapped })

	// Nothing to buy, or nothing to buy it with.
	if len(tappedCreatures) < 1 || len(untappedMana) < 2 {
		return
	}

	// Never offer more mana than there are creatures to untap with it.
	affordable := len(untappedMana)
	if affordable > len(tappedCreatures)*2 {
		affordable = len(tappedCreatures) * 2
	}

	paid := fx.Select(
		player,
		ctx.Match,
		player,
		match.MANAZONE,
		fmt.Sprintf("%s: You may tap any number of cards in your mana zone. For every 2 you tap, untap one of your creatures.", card.Name),
		1,
		affordable,
		true,
	)

	untaps := len(paid) / 2

	if untaps < 1 {
		return
	}

	for _, mana := range paid {
		player.TapCard(mana)
	}

	fx.SelectFilter(
		player,
		ctx.Match,
		player,
		match.BATTLEZONE,
		fmt.Sprintf("%s: Choose %d of your creatures to untap.", card.Name, untaps),
		untaps,
		untaps,
		false,
		func(x *match.Card) bool { return x.Tapped },
		false,
	).Map(func(creature *match.Card) {
		creature.Tapped = false
		ctx.Match.ReportActionInChat(player, fmt.Sprintf("%s untapped %s by tapping mana", player.Username(), creature.Name))
	})
}
