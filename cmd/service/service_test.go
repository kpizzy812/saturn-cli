package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceCommand_Metadata(t *testing.T) {
	cmd := NewServiceCommand()
	assert.Equal(t, "service", cmd.Use)
	assert.Contains(t, cmd.Aliases, "svc")
	assert.Contains(t, cmd.Aliases, "services")
	assert.NotEmpty(t, cmd.Short)
}

func TestServiceCommand_Subcommands(t *testing.T) {
	cmd := NewServiceCommand()
	names := make(map[string]bool)
	for _, sub := range cmd.Commands() {
		names[sub.Name()] = true
	}
	assert.True(t, names["list"])
	assert.True(t, names["get"])
	assert.True(t, names["create"])
	assert.True(t, names["start"])
	assert.True(t, names["stop"])
	assert.True(t, names["restart"])
	assert.True(t, names["delete"])
	assert.True(t, names["env"])
}

func TestServiceCreateCommand_Flags(t *testing.T) {
	cmd := NewCreateCommand()
	assert.NotNil(t, cmd.Flags().Lookup("server-uuid"))
	assert.NotNil(t, cmd.Flags().Lookup("project-uuid"))
	assert.NotNil(t, cmd.Flags().Lookup("environment-name"))
	assert.NotNil(t, cmd.Flags().Lookup("environment-uuid"))
	assert.NotNil(t, cmd.Flags().Lookup("name"))
	assert.NotNil(t, cmd.Flags().Lookup("instant-deploy"))
	assert.NotNil(t, cmd.Flags().Lookup("list-types"))
}

func TestServiceStartCommand_Args(t *testing.T) {
	cmd := NewStartCommand()
	require.NotNil(t, cmd.Args)

	err := cmd.Args(cmd, []string{})
	require.Error(t, err, "should require UUID argument")

	err = cmd.Args(cmd, []string{"some-uuid"})
	assert.NoError(t, err, "should accept single UUID")
}

func TestServiceStopCommand_Args(t *testing.T) {
	cmd := NewStopCommand()
	require.NotNil(t, cmd.Args)

	err := cmd.Args(cmd, []string{})
	require.Error(t, err, "should require UUID argument")

	err = cmd.Args(cmd, []string{"some-uuid"})
	assert.NoError(t, err, "should accept single UUID")
}

func TestServiceRestartCommand_Args(t *testing.T) {
	cmd := NewRestartCommand()
	require.NotNil(t, cmd.Args)

	err := cmd.Args(cmd, []string{})
	require.Error(t, err, "should require UUID argument")

	err = cmd.Args(cmd, []string{"some-uuid"})
	assert.NoError(t, err, "should accept single UUID")
}

func TestServiceCreateCommand_ValidatesType(t *testing.T) {
	cmd := NewCreateCommand()
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	// Should reject unknown service type
	cmd.SetArgs([]string{"not-a-real-service"})
	err := cmd.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid service type")
}

func TestServiceCreateCommand_RequiresArgs(t *testing.T) {
	cmd := NewCreateCommand()
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	// No args and no --list-types: should error
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	require.Error(t, err)
}

func TestColorizeServiceStatus(t *testing.T) {
	tests := []struct {
		status   string
		hasColor bool
	}{
		{"running", true},
		{"stopped", true},
		{"exited", true},
		{"restarting", true},
		{"error", true},
		{"failed", true},
		{"unhealthy", true},
		{"starting", true},
		{"deploying", true},
		{"unknown-status", false},
	}

	for _, tt := range tests {
		result := colorizeServiceStatus(tt.status)
		if tt.hasColor {
			// pterm adds ANSI escape codes
			assert.NotEqual(t, tt.status, result, "status %q should be colorized", tt.status)
		} else {
			assert.Equal(t, tt.status, result, "unknown status should be returned as-is")
		}
	}
}
