package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerCommand_Metadata(t *testing.T) {
	cmd := NewServerCommand()
	assert.Equal(t, "server", cmd.Use)
	assert.Contains(t, cmd.Aliases, "srv")
	assert.Contains(t, cmd.Aliases, "servers")
	assert.NotEmpty(t, cmd.Short)
}

func TestServerCommand_Subcommands(t *testing.T) {
	cmd := NewServerCommand()
	names := make(map[string]bool)
	for _, sub := range cmd.Commands() {
		names[sub.Name()] = true
	}
	assert.True(t, names["list"])
	assert.True(t, names["get"])
	assert.True(t, names["add"])
	assert.True(t, names["remove"])
	assert.True(t, names["validate"])
	assert.True(t, names["domains"])
}

func TestServerAddCommand_Args(t *testing.T) {
	cmd := NewAddCommand()
	require.NotNil(t, cmd.Args)

	err := cmd.Args(cmd, []string{})
	assert.Error(t, err, "should require 3 arguments")

	err = cmd.Args(cmd, []string{"name", "1.2.3.4"})
	assert.Error(t, err, "should require 3 arguments")

	err = cmd.Args(cmd, []string{"name", "1.2.3.4", "key-uuid"})
	assert.NoError(t, err, "should accept exactly 3 arguments")
}

func TestServerAddCommand_Flags(t *testing.T) {
	cmd := NewAddCommand()
	assert.NotNil(t, cmd.Flags().Lookup("port"))
	assert.NotNil(t, cmd.Flags().Lookup("user"))
	assert.NotNil(t, cmd.Flags().Lookup("validate"))
}

func TestServerRemoveCommand_Args(t *testing.T) {
	cmd := NewRemoveCommand()
	require.NotNil(t, cmd.Args)

	err := cmd.Args(cmd, []string{})
	assert.Error(t, err, "should require UUID argument")

	err = cmd.Args(cmd, []string{"some-uuid"})
	assert.NoError(t, err, "should accept single UUID")
}

func TestServerValidateCommand_Args(t *testing.T) {
	cmd := NewValidateCommand()
	require.NotNil(t, cmd.Args)

	err := cmd.Args(cmd, []string{})
	assert.Error(t, err, "should require UUID argument")

	err = cmd.Args(cmd, []string{"some-uuid"})
	assert.NoError(t, err, "should accept single UUID")
}

func TestColorizeBoolean(t *testing.T) {
	trueResult := colorizeBoolean(true)
	falseResult := colorizeBoolean(false)

	// pterm adds ANSI escape codes, so result should differ from plain text
	assert.Contains(t, trueResult, "yes")
	assert.Contains(t, falseResult, "no")
	// Colors should differ
	assert.NotEqual(t, trueResult, falseResult)
}
