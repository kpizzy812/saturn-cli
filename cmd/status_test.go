package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStatusCommand(t *testing.T) {
	cmd := NewStatusCommand()
	assert.Equal(t, "status", cmd.Use)
	assert.NotEmpty(t, cmd.Short)
	assert.NotNil(t, cmd.RunE)
}

func TestStatusCommandLong(t *testing.T) {
	cmd := NewStatusCommand()
	assert.NotEmpty(t, cmd.Long)
}
