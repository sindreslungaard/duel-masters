package fx

import (
	"duel-masters/game/match"
	"duel-masters/tests/scenario"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// Armored Cannon Balbaro evolves from a Human, and the scenario deck is made
	// of Immortal Baron, Vorg, which is one.
	evolutionEventEvolutionUID = "24353d06-89ef-4867-9513-485750d01e10"
	evolutionEventBaseUID      = "af3bc221-1cc2-4f58-83ea-2673ac2c66c5"
	evolutionEventSrc          = "evolution_event_test_setup"
)

// evolutionEventWatcher records every EvolutionEvent the match dispatches. It is
// attached to a card the watcher's owner keeps in hand, because card handlers
// run in every zone.
type evolutionEventWatcher struct {
	events []*match.EvolutionEvent
}

func (w *evolutionEventWatcher) attach(card *match.Card) {
	card.Use(func(_ *match.Card, ctx *match.Context) {
		if event, ok := ctx.Event.(*match.EvolutionEvent); ok {
			w.events = append(w.events, event)
		}
	})
}

// evolverOf resolves the player an EvolutionEvent names, the way a card handler
// reading MatchPlayerID has to. Which of the two players starts is not fixed, so
// a test must never assume the number itself.
func evolverOf(scn *scenario.TestScenario, event *match.EvolutionEvent) *match.Player {
	if event.MatchPlayerID == 1 {
		return scn.Match.Player1.Player
	}

	return scn.Match.Player2.Player
}

// putInBattlezone spawns a card straight into play.
func putInBattlezone(t *testing.T, scn *scenario.TestScenario, player *match.Player, uid string) *match.Card {
	t.Helper()

	_, err := player.SpawnCard(uid, match.HAND)
	require.NoError(t, err)

	card, err := scn.FindCard(player, match.HAND, uid)
	require.NoError(t, err)

	moved, err := player.MoveCard(card.ID, match.HAND, match.BATTLEZONE, evolutionEventSrc)
	require.NoError(t, err)
	require.Equal(t, match.BATTLEZONE, moved.Zone)

	return moved
}

// summonWithOwnMana pays for a card with copies of itself and plays it.
func summonWithOwnMana(t *testing.T, scn *scenario.TestScenario, player *match.PlayerReference, uid string) *match.Card {
	t.Helper()

	card, err := player.Player.SpawnCard(uid, match.HAND)
	require.NoError(t, err)

	for range card.ManaCost {
		_, err := player.Player.SpawnCard(uid, match.MANAZONE)
		require.NoError(t, err)
	}

	require.NoError(t, scn.ActionPlayCard(player, card.ID))
	require.NoError(t, scn.WaitForEventLoop())

	return card
}

func TestEvolutionEventCarriesTheWholePile(t *testing.T) {
	scn := scenario.New()
	player := scn.Match.CurrentPlayer()

	watcher := &evolutionEventWatcher{}
	observer, err := player.Player.SpawnCard(evolutionEventBaseUID, match.HAND)
	require.NoError(t, err)
	watcher.attach(observer)

	base := putInBattlezone(t, scn, player.Player, evolutionEventBaseUID)
	evolution := summonWithOwnMana(t, scn, player, evolutionEventEvolutionUID)

	require.Len(t, watcher.events, 1)
	assert.Equal(t, evolution.ID, watcher.events[0].CardID)
	assert.Equal(t, base.ID, watcher.events[0].BaseID)
	assert.Equal(t, player.Player, evolverOf(scn, watcher.events[0]))
}

func TestEvolutionEventFiresForEitherPlayer(t *testing.T) {
	scn := scenario.New()
	player := scn.Match.CurrentPlayer()
	opponent := scn.Match.PlayerRef(scn.Match.Opponent(player.Player))

	watcher := &evolutionEventWatcher{}
	observer, err := player.Player.SpawnCard(evolutionEventBaseUID, match.HAND)
	require.NoError(t, err)
	watcher.attach(observer)

	putInBattlezone(t, scn, opponent.Player, evolutionEventBaseUID)

	require.NoError(t, scn.ActionEndTurn(player))
	evolution := summonWithOwnMana(t, scn, opponent, evolutionEventEvolutionUID)

	require.Len(t, watcher.events, 1, "the opponent evolving is still an evolution")
	assert.Equal(t, evolution.ID, watcher.events[0].CardID)
	assert.Equal(t, opponent.Player, evolverOf(scn, watcher.events[0]), "and the event says whose it was")
}

func TestEvolutionEventFiresOnlyWhenSomethingEvolves(t *testing.T) {
	scn := scenario.New()
	player := scn.Match.CurrentPlayer()

	watcher := &evolutionEventWatcher{}
	observer, err := player.Player.SpawnCard(evolutionEventBaseUID, match.HAND)
	require.NoError(t, err)
	watcher.attach(observer)

	// An ordinary summon, and a creature moved straight into play, are neither
	// of them evolutions.
	summonWithOwnMana(t, scn, player, evolutionEventBaseUID)
	putInBattlezone(t, scn, player.Player, evolutionEventBaseUID)

	assert.Empty(t, watcher.events)
}

func TestEvolutionEventIsDispatchedWithTheBaseAlreadyHidden(t *testing.T) {
	scn := scenario.New()
	player := scn.Match.CurrentPlayer()

	watcher := &evolutionEventWatcher{}
	observer, err := player.Player.SpawnCard(evolutionEventBaseUID, match.HAND)
	require.NoError(t, err)

	base := putInBattlezone(t, scn, player.Player, evolutionEventBaseUID)

	// Recorded at dispatch time rather than afterwards, so the zones are the
	// ones handlers actually see. A handler resolves the ids on the event
	// against the player it names; the test already holds the cards.
	var baseZone, evolutionZone string
	var evolution *match.Card

	observer.Use(func(_ *match.Card, ctx *match.Context) {
		if event, ok := ctx.Event.(*match.EvolutionEvent); ok {
			watcher.events = append(watcher.events, event)
			baseZone = base.Zone

			if evolution != nil {
				evolutionZone = evolution.Zone
			}
		}
	})

	evolution, err = player.Player.SpawnCard(evolutionEventEvolutionUID, match.HAND)
	require.NoError(t, err)

	for range evolution.ManaCost {
		_, err := player.Player.SpawnCard(evolutionEventEvolutionUID, match.MANAZONE)
		require.NoError(t, err)
	}

	require.NoError(t, scn.ActionPlayCard(player, evolution.ID))
	require.NoError(t, scn.WaitForEventLoop())

	require.Len(t, watcher.events, 1)
	assert.Equal(t, base.ID, watcher.events[0].BaseID)
	assert.Equal(t, match.HIDDENZONE, baseZone, "the base is under the evolution already")
	assert.NotEqual(t, match.BATTLEZONE, evolutionZone, "and the evolution has not landed yet")
}

// TestEvolutionEntersTappedWhenItsBaseIsTapped guards against a regression
// where an evolution creature's tap state was set while resolving
// CardPlayedEvent, only for the unconditional Tapped = false every
// Player.MoveCard applies on arrival to silently wipe it out once the
// evolution creature itself moved from hand to the battle zone.
func TestEvolutionEntersTappedWhenItsBaseIsTapped(t *testing.T) {
	scn := scenario.New()
	player := scn.Match.CurrentPlayer()

	base := putInBattlezone(t, scn, player.Player, evolutionEventBaseUID)
	base.Tapped = true

	evolution := summonWithOwnMana(t, scn, player, evolutionEventEvolutionUID)

	require.Equal(t, match.BATTLEZONE, evolution.Zone)
	assert.True(t, evolution.Tapped, "it took the tap state of the base it evolved from")
}

// TestEvolutionEntersUntappedWhenItsBaseIsUntapped is the other half of the
// rule: landing on an untapped base leaves the evolution creature untapped.
func TestEvolutionEntersUntappedWhenItsBaseIsUntapped(t *testing.T) {
	scn := scenario.New()
	player := scn.Match.CurrentPlayer()

	putInBattlezone(t, scn, player.Player, evolutionEventBaseUID)

	evolution := summonWithOwnMana(t, scn, player, evolutionEventEvolutionUID)

	require.Equal(t, match.BATTLEZONE, evolution.Zone)
	assert.False(t, evolution.Tapped)
}
