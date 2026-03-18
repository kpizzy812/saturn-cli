package teams

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTeamsCommand_Metadata(t *testing.T) {
	cmd := NewTeamsCommand()
	assert.Equal(t, "teams", cmd.Use)
	assert.Contains(t, cmd.Aliases, "team")
	assert.NotEmpty(t, cmd.Short)
}

func TestTeamsCommand_Subcommands(t *testing.T) {
	cmd := NewTeamsCommand()
	names := make(map[string]bool)
	for _, sub := range cmd.Commands() {
		names[sub.Name()] = true
	}
	assert.True(t, names["list"])
	assert.True(t, names["get"])
	assert.True(t, names["current"])
	assert.True(t, names["members"])
}

func TestTeamGetCommand_Args(t *testing.T) {
	cmd := NewGetCommand()
	require.NotNil(t, cmd.Args)

	err := cmd.Args(cmd, []string{})
	assert.Error(t, err, "should require team ID argument")

	err = cmd.Args(cmd, []string{"42"})
	assert.NoError(t, err, "should accept single ID")
}
