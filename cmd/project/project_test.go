package project

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProjectCommand_Metadata(t *testing.T) {
	cmd := NewProjectCommand()
	assert.Equal(t, "project", cmd.Use)
	assert.Contains(t, cmd.Aliases, "proj")
	assert.Contains(t, cmd.Aliases, "projects")
	assert.NotEmpty(t, cmd.Short)
}

func TestProjectCommand_Subcommands(t *testing.T) {
	cmd := NewProjectCommand()
	names := make(map[string]bool)
	for _, sub := range cmd.Commands() {
		names[sub.Name()] = true
	}
	assert.True(t, names["list"])
	assert.True(t, names["get"])
	assert.True(t, names["create"])
}

func TestProjectCreateCommand_Flags(t *testing.T) {
	cmd := NewCreateCommand()
	assert.NotNil(t, cmd.Flags().Lookup("name"))
	assert.NotNil(t, cmd.Flags().Lookup("description"))
}

func TestProjectCreateCommand_RequiresName(t *testing.T) {
	cmd := NewCreateCommand()
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	cmd.SetArgs([]string{})
	err := cmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "--name is required")
}

func TestProjectGetCommand_Args(t *testing.T) {
	cmd := NewGetCommand()
	assert.NotNil(t, cmd.Args)

	err := cmd.Args(cmd, []string{})
	assert.Error(t, err, "should require UUID argument")

	err = cmd.Args(cmd, []string{"some-uuid"})
	assert.NoError(t, err, "should accept single UUID")
}
