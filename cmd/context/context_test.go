package context

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewContextCommand_Metadata(t *testing.T) {
	cmd := NewContextCommand()

	assert.NotEmpty(t, cmd.Use)
	assert.NotEmpty(t, cmd.Short)

	subCmds := make(map[string]bool)
	for _, c := range cmd.Commands() {
		subCmds[c.Name()] = true
		for _, alias := range c.Aliases {
			subCmds[alias] = true
		}
	}

	assert.True(t, subCmds["add"], "should have add subcommand")
	assert.True(t, subCmds["list"], "should have list subcommand")
	assert.True(t, subCmds["use"], "should have use subcommand")
	assert.True(t, subCmds["delete"] || subCmds["remove"], "should have delete or remove subcommand")
}

func TestNewAddCommand_Args(t *testing.T) {
	cmd := NewAddCommand()

	t.Run("rejects zero arguments", func(t *testing.T) {
		require.Error(t, cmd.Args(cmd, []string{}))
	})

	t.Run("rejects one argument", func(t *testing.T) {
		require.Error(t, cmd.Args(cmd, []string{"myserver"}))
	})

	t.Run("rejects two arguments", func(t *testing.T) {
		require.Error(t, cmd.Args(cmd, []string{"myserver", "https://example.com"}))
	})

	t.Run("accepts exactly three arguments", func(t *testing.T) {
		assert.NoError(t, cmd.Args(cmd, []string{"myserver", "https://example.com", "token-123"}))
	})

	t.Run("rejects four arguments", func(t *testing.T) {
		require.Error(t, cmd.Args(cmd, []string{"myserver", "https://example.com", "token-123", "extra"}))
	})
}

func TestNewAddCommand_Flags(t *testing.T) {
	cmd := NewAddCommand()

	t.Run("has default flag with default false", func(t *testing.T) {
		f := cmd.Flags().Lookup("default")
		require.NotNil(t, f, "default flag must exist")
		assert.Equal(t, "false", f.DefValue)
	})

	t.Run("has force flag with default false", func(t *testing.T) {
		f := cmd.Flags().Lookup("force")
		require.NotNil(t, f, "force flag must exist")
		assert.Equal(t, "false", f.DefValue)
	})
}

func TestNewUseCommand_Args(t *testing.T) {
	cmd := NewUseCommand()

	t.Run("rejects zero arguments", func(t *testing.T) {
		require.Error(t, cmd.Args(cmd, []string{}))
	})

	t.Run("rejects two arguments", func(t *testing.T) {
		require.Error(t, cmd.Args(cmd, []string{"ctx-1", "ctx-2"}))
	})

	t.Run("accepts exactly one argument", func(t *testing.T) {
		assert.NoError(t, cmd.Args(cmd, []string{"myserver"}))
	})
}

func TestNewRemoveCommand_Args(t *testing.T) {
	cmd := NewDeleteCommand()

	t.Run("rejects zero arguments", func(t *testing.T) {
		require.Error(t, cmd.Args(cmd, []string{}))
	})

	t.Run("rejects two arguments", func(t *testing.T) {
		require.Error(t, cmd.Args(cmd, []string{"ctx-1", "ctx-2"}))
	})

	t.Run("accepts exactly one argument", func(t *testing.T) {
		assert.NoError(t, cmd.Args(cmd, []string{"myserver"}))
	})
}
