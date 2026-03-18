package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewApplicationCommand_Metadata(t *testing.T) {
	cmd := NewAppCommand()

	assert.NotEmpty(t, cmd.Use)
	assert.NotEmpty(t, cmd.Short)

	subCmds := make(map[string]bool)
	for _, c := range cmd.Commands() {
		subCmds[c.Name()] = true
		// Also register aliases so tests can find commands by any name
		for _, alias := range c.Aliases {
			subCmds[alias] = true
		}
	}

	assert.True(t, subCmds["list"], "should have list subcommand")
	assert.True(t, subCmds["get"], "should have get subcommand")
	assert.True(t, subCmds["start"], "should have start subcommand")
	assert.True(t, subCmds["stop"], "should have stop subcommand")
	assert.True(t, subCmds["restart"], "should have restart subcommand")
	assert.True(t, subCmds["logs"], "should have logs subcommand")
	assert.True(t, subCmds["env"], "should have env subcommand")
}

func TestNewLogsCommand_Flags(t *testing.T) {
	cmd := NewLogsCommand()

	t.Run("has lines flag with default 100", func(t *testing.T) {
		f := cmd.Flags().Lookup("lines")
		require.NotNil(t, f, "lines flag must exist")
		assert.Equal(t, "100", f.DefValue)
	})

	t.Run("has follow flag with default false", func(t *testing.T) {
		f := cmd.Flags().Lookup("follow")
		require.NotNil(t, f, "follow flag must exist")
		assert.Equal(t, "false", f.DefValue)
	})
}

func TestNewStartCommand_Args(t *testing.T) {
	cmd := NewStartCommand()

	t.Run("rejects zero arguments", func(t *testing.T) {
		require.Error(t, cmd.Args(cmd, []string{}))
	})

	t.Run("rejects two arguments", func(t *testing.T) {
		require.Error(t, cmd.Args(cmd, []string{"uuid-1", "uuid-2"}))
	})

	t.Run("accepts exactly one argument", func(t *testing.T) {
		assert.NoError(t, cmd.Args(cmd, []string{"uuid-123"}))
	})
}

func TestNewStopCommand_Args(t *testing.T) {
	cmd := NewStopCommand()

	t.Run("rejects zero arguments", func(t *testing.T) {
		require.Error(t, cmd.Args(cmd, []string{}))
	})

	t.Run("rejects two arguments", func(t *testing.T) {
		require.Error(t, cmd.Args(cmd, []string{"uuid-1", "uuid-2"}))
	})

	t.Run("accepts exactly one argument", func(t *testing.T) {
		assert.NoError(t, cmd.Args(cmd, []string{"uuid-123"}))
	})
}

func TestNewRestartCommand_Args(t *testing.T) {
	cmd := NewRestartCommand()

	t.Run("rejects zero arguments", func(t *testing.T) {
		require.Error(t, cmd.Args(cmd, []string{}))
	})

	t.Run("rejects two arguments", func(t *testing.T) {
		require.Error(t, cmd.Args(cmd, []string{"uuid-1", "uuid-2"}))
	})

	t.Run("accepts exactly one argument", func(t *testing.T) {
		assert.NoError(t, cmd.Args(cmd, []string{"uuid-123"}))
	})
}
