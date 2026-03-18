package output

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSuccess(t *testing.T) {
	var buf bytes.Buffer
	s := NewStyler(&buf)
	s.Success("deploy complete")
	assert.Contains(t, buf.String(), "deploy complete")
}

func TestError(t *testing.T) {
	var buf bytes.Buffer
	s := NewStyler(&buf)
	s.Error("failed to connect")
	assert.Contains(t, buf.String(), "failed to connect")
}

func TestInfo(t *testing.T) {
	var buf bytes.Buffer
	s := NewStyler(&buf)
	s.Info("checking status")
	assert.Contains(t, buf.String(), "checking status")
}

func TestWarning(t *testing.T) {
	var buf bytes.Buffer
	s := NewStyler(&buf)
	s.Warning("deprecated flag")
	assert.Contains(t, buf.String(), "deprecated flag")
}

func TestURL(t *testing.T) {
	var buf bytes.Buffer
	s := NewStyler(&buf)
	s.URL("https://app.saturn.ac")
	assert.Contains(t, buf.String(), "https://app.saturn.ac")
}
