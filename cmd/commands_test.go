package cmd

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRootCommand(t *testing.T) {
	root := NewRootCommand()

	t.Run("has expected subcommands", func(t *testing.T) {
		names := make(map[string]bool)
		for _, c := range root.Commands() {
			names[c.Name()] = true
		}
		assert.True(t, names["deploy"], "deploy command should be registered")
		assert.True(t, names["server"] || names["servers"], "server command should be registered")
		assert.True(t, names["app"] || names["application"], "application command should be registered")
		assert.True(t, names["database"] || names["db"], "database command should be registered")
		assert.True(t, names["project"] || names["projects"], "project command should be registered")
		assert.True(t, names["pull"], "pull command should be registered")
	})
}

func TestPullCommandFlags(t *testing.T) {
	cmd := NewPullCommand()

	t.Run("has output flag", func(t *testing.T) {
		f := cmd.Flags().Lookup("output")
		require.NotNil(t, f)
		assert.Equal(t, ".env.local", f.DefValue)
	})

	t.Run("has force flag", func(t *testing.T) {
		f := cmd.Flags().Lookup("force")
		require.NotNil(t, f)
		assert.Equal(t, "false", f.DefValue)
	})

	t.Run("requires exactly 1 argument", func(t *testing.T) {
		// ArgsFunc should reject 0 and 2 args, accept exactly 1
		require.Error(t, cmd.Args(cmd, []string{}))
		require.Error(t, cmd.Args(cmd, []string{"a", "b"}))
		assert.NoError(t, cmd.Args(cmd, []string{"uuid-123"}))
	})
}

func TestDeployCommandAliases(t *testing.T) {
	root := NewRootCommand()
	var deployCmd *cobra.Command
	for _, c := range root.Commands() {
		if c.Name() == "deploy" {
			deployCmd = c
			break
		}
	}
	require.NotNil(t, deployCmd, "deploy command must exist")

	aliases := deployCmd.Aliases
	assert.Contains(t, aliases, "dep")
	assert.Contains(t, aliases, "deployment")
}

func TestServerCommandAliases(t *testing.T) {
	root := NewRootCommand()
	var serverCmd *cobra.Command
	for _, c := range root.Commands() {
		if c.Name() == "server" || c.Name() == "servers" {
			serverCmd = c
			break
		}
	}
	require.NotNil(t, serverCmd)
	assert.Contains(t, serverCmd.Aliases, "srv")
}

func TestProjectCommandAliases(t *testing.T) {
	root := NewRootCommand()
	var projectCmd *cobra.Command
	for _, c := range root.Commands() {
		if c.Name() == "project" || c.Name() == "projects" {
			projectCmd = c
			break
		}
	}
	require.NotNil(t, projectCmd)
	assert.Contains(t, projectCmd.Aliases, "proj")
}
