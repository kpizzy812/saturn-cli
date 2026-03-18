package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatabaseCommand_Metadata(t *testing.T) {
	cmd := NewDatabaseCommand()
	assert.Equal(t, "database", cmd.Use)
	assert.Contains(t, cmd.Aliases, "db")
	assert.Contains(t, cmd.Aliases, "dbs")
	assert.Contains(t, cmd.Aliases, "databases")
	assert.NotEmpty(t, cmd.Short)
}

func TestDatabaseCommand_Subcommands(t *testing.T) {
	cmd := NewDatabaseCommand()
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
	assert.True(t, names["backup"])
}

func TestDatabaseCreateCommand_Flags(t *testing.T) {
	cmd := NewCreateCommand()
	assert.NotNil(t, cmd.Flags().Lookup("server-uuid"))
	assert.NotNil(t, cmd.Flags().Lookup("project-uuid"))
	assert.NotNil(t, cmd.Flags().Lookup("environment-name"))
	assert.NotNil(t, cmd.Flags().Lookup("environment-uuid"))
	assert.NotNil(t, cmd.Flags().Lookup("name"))
}

func TestDatabaseStartCommand_Args(t *testing.T) {
	cmd := NewStartCommand()
	require.NotNil(t, cmd.Args)

	err := cmd.Args(cmd, []string{})
	assert.Error(t, err, "should require UUID argument")

	err = cmd.Args(cmd, []string{"some-uuid"})
	assert.NoError(t, err, "should accept single UUID")
}

func TestDatabaseStopCommand_Args(t *testing.T) {
	cmd := NewStopCommand()
	require.NotNil(t, cmd.Args)

	err := cmd.Args(cmd, []string{})
	assert.Error(t, err, "should require UUID argument")

	err = cmd.Args(cmd, []string{"some-uuid"})
	assert.NoError(t, err, "should accept single UUID")
}

func TestDatabaseRestartCommand_Args(t *testing.T) {
	cmd := NewRestartCommand()
	require.NotNil(t, cmd.Args)

	err := cmd.Args(cmd, []string{})
	assert.Error(t, err, "should require UUID argument")

	err = cmd.Args(cmd, []string{"some-uuid"})
	assert.NoError(t, err, "should accept single UUID")
}

func TestDatabaseCreateCommand_ValidatesType(t *testing.T) {
	cmd := NewCreateCommand()
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	// Should reject unknown DB type
	cmd.SetArgs([]string{"oracle"})
	err := cmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid database type")
}

func TestDatabaseCreateCommand_ValidTypes(t *testing.T) {
	validTypes := []string{"postgresql", "mysql", "mariadb", "mongodb", "redis", "keydb", "clickhouse", "dragonfly"}
	cmd := NewCreateCommand()
	// Just verify --help works for each valid type to confirm they're accepted
	for _, dbType := range validTypes {
		// We just check the Args validator directly
		err := cmd.Args(cmd, []string{dbType})
		assert.NoError(t, err, "type %q should be valid positional arg format", dbType)
	}
}

func TestColorizeDatabaseStatus(t *testing.T) {
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
		{"unknown-status", false},
	}

	for _, tt := range tests {
		result := colorizeDatabaseStatus(tt.status)
		if tt.hasColor {
			assert.NotEqual(t, tt.status, result, "status %q should be colorized", tt.status)
		} else {
			assert.Equal(t, tt.status, result, "unknown status should be returned as-is")
		}
	}
}
