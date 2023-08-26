package fx

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"fmt"
)

/*
Summoning an evolution creature works just like summoning a regular creature except you can only summon an
evolution creature when you have the correct type of creature (whether it needs a specific race, civilization,
or ability) already in the battle zone. The evolution creature tells you what the creature it evolves from needs to have.

Once you summon an evolution creature, place it on top of the creature it "evolves" from.
Leave that creature underneath the evolution creature, but ignore it any effects, power or other attributes it might have.
Only the evolution creature's name, abilities, civilization, and power matter.
Evolution creatures don't get summoning sickness even if the creature they evolved from had summoning sickness,
so you can attack immediately after you put them into the battle zone.
There's no limit to the number of evolution creatures you can put on top of each other.

An evolution creature is put into the battle zone in the same tap position as the creature it is put on.
If it is put on a creature that is untapped, it is untapped. If it is put on a tapped creature, it enters the battle zone tapped.

If an evolution creature is moved from the battle zone to anywhere else, then the whole pile moves,
not just the evolution creature on top. As long as the evolution creature is in play,
it and the card it evolved from count as only one card. As soon as the pile ends up somewhere other than the battle zone,
they are separate cards again. So if a card effect tells you put your evolution creature into the mana zone,
you get that many separate cards in your mana zone.

However, when a card underneath the evolution creature is removed, the evolution creature itself remains in the battle zone.
If only the top card of the evolution creature is removed (by an ability such as Dias Zeta, the Temporal Suppressor or Royal Durian),
the player controlling the evolution creature can only have the same amount of creatures under it equal to the number of creatures
required for the evolution, and the rest into your graveyard. For example, if the top card of Zeek Cavalie, Lord of Dragon Spirits is removed,
one of the creatures under it stays in the battle zone and the rest goes into your graveyard. Super Infinite Evolution is exempt from
this rule since it has no designated number of evolution bait. Evolution creatures and linked Gods can be counted as 1 creature.
Also, removing the top card now counts as targeting the creature itself, and will trigger effects that involve the creature being
targeted/leaving the field, etc.
*/

// Evolution has default behaviour for evolution cards according to the rules commented above
func Evolution(card *match.Card, ctx *match.Context) {

	if _, ok := ctx.Event.(*match.UntapStep); ok {

		card.AddCondition(cnd.Evolution, true, card.ID)
	}

	if event, ok := ctx.Event.(*match.PlayCardEvent); ok {

		if event.CardID != card.ID {
			return
		}

		ctx.ScheduleAfter(func() {
			card.RemoveCondition(cnd.SummoningSickness)
		})

	}

	if event, ok := ctx.Event.(*match.CardPlayedEvent); ok {

		if event.CardID != card.ID {
			return
		}

		creatures := match.Filter(
			card.Player,
			ctx.Match,
			card.Player,
			match.BATTLEZONE,
			fmt.Sprintf("Choose 1 %s to evolve %s from", card.Family, card.Name),
			1,
			1,
			false,
			func(x *match.Card) bool { return x.SharesAFamily(card.Family) },
		)

		if len(creatures) < 1 {
			ctx.InterruptFlow()
			ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("There are no cards for %s to evolve from in your battle zone", card.Name))
			return
		}

		creature := creatures[0]

		card.ClearAttachments()
		card.Tapped = creature.Tapped
		card.Player.MoveCard(creature.ID, match.BATTLEZONE, match.HIDDENZONE)
		card.Attach(creature)
		card.AddCondition(cnd.Evolution, true, card.ID)
	}

	// Card moved
	if event, ok := ctx.Event.(*match.CardMoved); ok {

		if event.CardID != card.ID || event.To == match.BATTLEZONE || event.To == match.HIDDENZONE {
			return
		}

		for _, creature := range card.Attachments() {
			card.Player.MoveCard(creature.ID, match.HIDDENZONE, event.To)
			ctx.Match.Chat("Server", fmt.Sprintf("%s was sent to the %s together with %s", creature.Name, event.To, card.Name))
		}

		card.ClearAttachments()

	}

}
