package fx

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
)

// waveStrikerThreshold is the printed "2 or more other creatures".
const waveStrikerThreshold = 2

// WaveStriker implements the keyword itself: "Wave striker (While 2 or more
// other creatures in the battle zone have "wave striker," this creature has its
// Wave striker ability.)"
//
// The keyword only marks the creature. What the ability does, and whether it is
// switched on right now, is the job of WhileWaveStriker and its siblings.
//
// Like the other printed keywords this is rebuilt every untap step and in every
// zone, so an effect can ask whether a card in hand is a wave striker.
func WaveStriker(card *match.Card, ctx *match.Context) {
	if _, ok := ctx.Event.(*match.UntapStep); ok {
		card.AddUniqueSourceCondition(cnd.WaveStriker, true, card.ID)
	}
}

// WaveStrikersInBattleZone counts the wave strikers a player has in play.
func WaveStrikersInBattleZone(player *match.Player) int {
	return len(FindFilter(
		player,
		match.BATTLEZONE,
		func(x *match.Card) bool { return x.HasCondition(cnd.WaveStriker) },
	))
}

// WaveStrikerActive reports whether the card's wave striker ability is switched
// on: the threshold counts *other* creatures, and the printed text says "in the
// battle zone" rather than "you have", so both players' creatures count towards
// it.
//
// A card outside the battle zone has no ability to switch on.
func WaveStrikerActive(card *match.Card, m *match.Match) bool {
	if card.Zone != match.BATTLEZONE {
		return false
	}

	others := WaveStrikersInBattleZone(card.Player) + WaveStrikersInBattleZone(m.Opponent(card.Player))

	if card.HasCondition(cnd.WaveStriker) {
		others--
	}

	return others >= waveStrikerThreshold
}

// WaveStrikerReady is WaveStrikerActive in the shape fx.When expects.
func WaveStrikerReady(card *match.Card, ctx *match.Context) bool {
	return WaveStrikerActive(card, ctx.Match)
}

// WhileWaveStriker runs the handler only while the ability is switched on.
//
// For a triggered ability this is checked as the trigger fires, which is what
// the printed wording means: a wave striker summoned as the third one sees
// itself in the battle zone and counts the two that were already there.
func WhileWaveStriker(h match.HandlerFunc) match.HandlerFunc {
	return When(WaveStrikerReady, h)
}

// WaveStrikerPower is a PowerModifier that adds the bonus while the ability is
// switched on. Power is read constantly, so this is re-evaluated rather than
// stored, and it needs no cleanup when the board changes underneath it.
func WaveStrikerPower(card *match.Card, bonus int) func(m *match.Match, attacking bool) int {
	return func(m *match.Match, attacking bool) int {
		if !WaveStrikerActive(card, m) {
			return 0
		}

		return bonus
	}
}

// WaveStrikerGrant keeps a condition on the creature itself in step with
// whether the ability is switched on.
//
// It is re-evaluated on every event rather than installed once, because the
// count it depends on changes whenever a creature enters or leaves the battle
// zone, and the ability has to follow it in both directions.
func WaveStrikerGrant(condition string, val interface{}) match.HandlerFunc {
	return func(card *match.Card, ctx *match.Context) {
		if WaveStrikerActive(card, ctx.Match) {
			card.AddUniqueSourceCondition(condition, val, card.ID)
			return
		}

		card.RemoveSpecificConditionBySource(condition, card.ID)
	}
}

// WaveStrikerGrantToOwnCreatures does the same for every creature its
// controller has in the battle zone.
//
// Cards that have left the battle zone are swept too: a creature bounced to
// hand keeps its conditions until the end of the turn, and it must not carry a
// stale bonus back into play.
func WaveStrikerGrantToOwnCreatures(condition string, val interface{}) match.HandlerFunc {
	return func(card *match.Card, ctx *match.Context) {
		active := WaveStrikerActive(card, ctx.Match)

		zones := []string{match.BATTLEZONE, match.HAND, match.GRAVEYARD, match.MANAZONE, match.SHIELDZONE, match.HIDDENZONE}

		FindMultiple(card.Player, zones).Map(func(x *match.Card) {
			if active && x.Zone == match.BATTLEZONE {
				x.AddUniqueSourceCondition(condition, val, card.ID)
				return
			}

			x.RemoveSpecificConditionBySource(condition, card.ID)
		})
	}
}
