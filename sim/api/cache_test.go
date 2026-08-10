package api

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestResolveCardCatalog covers the ways the server is started, which each leave
// the card catalog in a different place relative to the working directory. The
// container image flattens the Go module into the working directory, the README
// starts the server from inside the module, and running it from the repository
// root leaves the module in a subdirectory. Missing the file is not fatal, but it
// costs every card its rules text, so all three have to resolve.
func TestResolveCardCatalog(t *testing.T) {
	cases := []struct {
		name     string
		place    string
		startIn  string
		expected string
	}{
		{
			name:     "container image, module flattened into the working directory",
			place:    "DuelMastersCards.json",
			startIn:  ".",
			expected: "DuelMastersCards.json",
		},
		{
			name:     "repository root, module in a subdirectory",
			place:    filepath.Join("sim", "DuelMastersCards.json"),
			startIn:  ".",
			expected: "sim/DuelMastersCards.json",
		},
		{
			name:     "inside the module, as the README describes",
			place:    "DuelMastersCards.json",
			startIn:  "sim",
			expected: "../DuelMastersCards.json",
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()

			path := filepath.Join(root, test.place)
			require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
			require.NoError(t, os.WriteFile(path, []byte(`{"cards":[]}`), 0o644))
			require.NoError(t, os.MkdirAll(filepath.Join(root, "sim"), 0o755))

			restore := chdir(t, filepath.Join(root, test.startIn))
			defer restore()

			resolved, found := resolveCardCatalog()
			require.True(t, found, "the card catalog should have been found")
			require.Equal(t, test.expected, filepath.ToSlash(resolved))
		})
	}

	t.Run("reports not found when the catalog is absent", func(t *testing.T) {
		restore := chdir(t, t.TempDir())
		defer restore()

		_, found := resolveCardCatalog()
		require.False(t, found)
	})
}

func chdir(t *testing.T, dir string) func() {
	t.Helper()

	previous, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(dir))

	return func() {
		require.NoError(t, os.Chdir(previous))
	}
}
