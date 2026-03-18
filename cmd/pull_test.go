package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPullCommand(t *testing.T) {
	cmd := NewPullCommand()

	t.Run("command metadata", func(t *testing.T) {
		assert.Equal(t, "pull <uuid>", cmd.Use)
		assert.NotEmpty(t, cmd.Short)
	})

	t.Run("flags", func(t *testing.T) {
		outputFlag := cmd.Flags().Lookup("output")
		require.NotNil(t, outputFlag)
		assert.Equal(t, ".env.local", outputFlag.DefValue)

		forceFlag := cmd.Flags().Lookup("force")
		require.NotNil(t, forceFlag)
		assert.Equal(t, "false", forceFlag.DefValue)
	})
}
