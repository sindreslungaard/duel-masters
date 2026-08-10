package tests

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// rawActionReceive matches a receive from a player's action channel, such as
// `action := <-p.Action` or `case <-player.Action:`.
var rawActionReceive = regexp.MustCompile(`<-\s*[\w.]*\bAction\b`)

// allowedRawActionReceives lists the files that are allowed to touch the action
// channel directly. Player.NextAction owns the receive, and the parser owns the
// send and the drain that keeps a malicious client from stacking up answers.
var allowedRawActionReceives = map[string]bool{
	filepath.Join("game", "match", "player.go"): true,
	filepath.Join("game", "match", "match.go"):  true,
}

// TestNoRawActionChannelReceives keeps the whole engine on Player.NextAction.
//
// Player.Dispose closes the action channel when a match is disposed, and a
// receive from a closed channel yields the zero PlayerAction immediately and
// forever. A prompt loop that receives without checking the second, "still
// open" value reads that as an invalid answer, retries, and spins a core for
// the lifetime of the process. NextAction reports closure so the caller can
// abandon its prompt instead.
func TestNoRawActionChannelReceives(t *testing.T) {
	root := ".."

	var offenders []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if info.Name() == "node_modules" || info.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}

		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		if allowedRawActionReceives[relative] {
			return nil
		}

		contents, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		for number, line := range strings.Split(string(contents), "\n") {
			if rawActionReceive.MatchString(line) {
				offenders = append(offenders, filepath.ToSlash(relative)+":"+strconv.Itoa(number+1)+": "+strings.TrimSpace(line))
			}
		}

		return nil
	})

	require.NoError(t, err)
	require.Empty(
		t,
		offenders,
		"receive player answers through Player.NextAction, which aborts when the match has been disposed. "+
			"Receiving from Player.Action directly spins forever once the channel is closed:\n%s",
		strings.Join(offenders, "\n"),
	)
}

// TestNoRecoverInEffectCode keeps the path from a prompt to the match event loop
// clear of recovers.
//
// Player.NextAction aborts by panicking when the match is disposed while it is
// waiting, and the match event loop recovers that abort. A recover in between
// would swallow it, so the effect would carry on resolving for a game that no
// longer exists. If effect code ever genuinely needs one, it has to re-raise
// anything match.IsMatchDisposed reports.
func TestNoRecoverInEffectCode(t *testing.T) {
	recoverCall := regexp.MustCompile(`\brecover\(\)`)

	var offenders []string

	for _, directory := range []string{filepath.Join("..", "game", "cards"), filepath.Join("..", "game", "fx")} {
		err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
				return nil
			}

			contents, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			for number, line := range strings.Split(string(contents), "\n") {
				if recoverCall.MatchString(line) {
					offenders = append(offenders, filepath.ToSlash(path)+":"+strconv.Itoa(number+1)+": "+strings.TrimSpace(line))
				}
			}

			return nil
		})

		require.NoError(t, err)
	}

	require.Empty(
		t,
		offenders,
		"card and effect code must not recover, or it will swallow the abort Player.NextAction raises "+
			"when the match is disposed while a prompt is open:\n%s",
		strings.Join(offenders, "\n"),
	)
}
