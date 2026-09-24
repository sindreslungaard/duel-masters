package dm12

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// AuraPegasusAvatarOfLife ...
func AuraPegasusAvatarOfLife(c *match.Card) {

	c.Name = "Aura Pegasus, Avatar of Life"
	c.Power = 12000
	c.Civs = []string{civ.Light, civ.Nature}
	c.Family = []string{family.Pegasus}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Light, civ.Nature}

	c.Use(
		fx.Creature,
		fx.PutIntoManaZoneTapped,
		fx.Triplebreaker,
		fx.VortexEvolution(
			family.HornedBeast, func(x *match.Card) bool { return x.HasFamily(family.HornedBeast) },
			family.AngelCommand, func(x *match.Card) bool { return x.HasFamily(family.AngelCommand) },
		),
		fx.When(fx.Attacking, auraPegasusAvatarOfLifeRevealTopCard),
		fx.When(fx.LeftBattlezone, auraPegasusAvatarOfLifeRevealTopCard),
	)

}

// auraPegasusAvatarOfLifeRevealTopCard implements "reveal the top card of
// your deck. If it's a non-evolution creature, put it into the battle zone.
// Otherwise, put it into your hand." cnd.Evolution is rebuilt for every
// creature at every untap step regardless of zone, so it is already correct
// for a card that has never left the deck.
func auraPegasusAvatarOfLifeRevealTopCard(card *match.Card, ctx *match.Context) {

	revealed := card.Player.PeekDeck(1)

	if len(revealed) < 1 {
		return
	}

	top := revealed[0]

	if top.HasCondition(cnd.Creature) && !top.HasCondition(cnd.Evolution) {
		fx.ForcePutCreatureIntoBZ(ctx, top, match.DECK, card)
		return
	}

	moved, err := card.Player.MoveCard(top.ID, match.DECK, match.HAND, card.ID)

	if err == nil {
		ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s put %s into their hand from the top of their deck by %s's effect", card.Player.Username(), moved.Name, card.Name))
	}

}
