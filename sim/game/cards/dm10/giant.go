package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// BodaciousGiant ...
func BodaciousGiant(c *match.Card) {

	c.Name = "Bodacious Giant"
	c.Power = 12000
	c.Civs = []string{civ.Nature}
	c.Family = []string{family.Giant}
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Nature}

	attackedThisTurn := false

	taunting := func(card *match.Card, ctx *match.Context) bool {
		return card.Zone == match.BATTLEZONE &&
			card.Tapped &&
			!ctx.Match.IsPlayerTurn(card.Player) &&
			!attackedThisTurn
	}

	c.Use(
		fx.Creature,
		fx.Doublebreaker,
		func(card *match.Card, ctx *match.Context) {

			// Reset at the start of every turn so "that turn" is scoped correctly.
			if _, ok := ctx.Event.(*match.UntapStep); ok {
				attackedThisTurn = false
				return
			}

			// Block is dispatched with the chosen target once an attack on this
			// creature is confirmed, whether or not it ends up blocked.
			if event, ok := ctx.Event.(*match.Block); ok && event.AttackedCardID == card.ID {
				attackedThisTurn = true
				return
			}

			if !taunting(card, ctx) {
				return
			}

			opponent := ctx.Match.Opponent(card.Player)

			switch event := ctx.Event.(type) {

			case *match.EndTurnEvent:
				// Only creatures that could actually attack this one are held back.
				for _, creature := range fx.Find(opponent, match.BATTLEZONE) {
					if ctx.Cancelled() {
						// Another creature already blocked the end of the turn;
						// warning again for each of them would spam the player.
						return
					}

					if !bodaciousGiantCanBeAttackedBy(card, creature) {
						continue
					}

					fx.ForceAttack(creature, ctx)
				}

			case *match.AttackCreature:
				// fx.Creature fills AttackableCreatures while handling this event
				// and only prompts in a scheduled callback. This card belongs to
				// the non-active player, so its handlers run after the attacker's
				// and the list can still be narrowed here.
				//
				// With several taunting creatures the first one to narrow the list
				// wins, because the later ones no longer find themselves in it.
				if !bodaciousGiantIsAttackable(card, event.AttackableCreatures) {
					return
				}

				event.AttackableCreatures = []*match.Card{card}

			case *match.AttackPlayer:
				attacker, err := opponent.GetCard(event.CardID, match.BATTLEZONE)
				if err != nil || !bodaciousGiantCanBeAttackedBy(card, attacker) {
					return
				}

				ctx.Match.WarnPlayer(attacker.Player, fmt.Sprintf("%s must attack %s while it is able to", attacker.Name, card.Name))
				ctx.InterruptFlow()

			}

		},
	)

}

func bodaciousGiantIsAttackable(card *match.Card, attackable []*match.Card) bool {
	for _, candidate := range attackable {
		if candidate == card {
			return true
		}
	}

	return false
}

// bodaciousGiantCanBeAttackedBy mirrors the engine's own rule for which
// creatures an attacker may be pointed at, so "if able" cannot drift from what
// the attack prompt would actually offer.
func bodaciousGiantCanBeAttackedBy(card *match.Card, attacker *match.Card) bool {
	if attacker.HasCondition(cnd.CantAttackCreatures) {
		return false
	}

	return card.HasCondition(cnd.TreatedAsTapped) ||
		(!card.HasCondition(cnd.CantBeAttacked) && (card.Tapped || attacker.HasCondition(cnd.AttackUntapped)))
}
