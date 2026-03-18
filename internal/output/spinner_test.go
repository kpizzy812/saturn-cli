package output

import (
	"testing"

	"github.com/pterm/pterm"
	"github.com/stretchr/testify/assert"
)

func init() {
	// Suppress pterm output in tests to avoid noisy CI output
	pterm.DisableOutput()
}

func TestSpinnerStartStop(t *testing.T) {
	sp := NewSpinner("deploying...")
	assert.NotNil(t, sp)

	sp.UpdateText("building image...")
	sp.Success("deployed!")
}

func TestSpinnerFail(t *testing.T) {
	sp := NewSpinner("deploying...")
	sp.Fail("deployment failed")
}
