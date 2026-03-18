package output

import (
	"fmt"
	"io"
	"os"

	"github.com/pterm/pterm"
)

// Styler provides colored and styled terminal output
type Styler struct {
	writer io.Writer
}

// NewStyler creates a new Styler writing to the given writer
func NewStyler(w io.Writer) *Styler {
	return &Styler{writer: w}
}

// Success prints a green checkmark message
func (s *Styler) Success(msg string) {
	fmt.Fprintln(s.writer, pterm.FgGreen.Sprint("✓ "+msg))
}

// Error prints a red cross message
func (s *Styler) Error(msg string) {
	fmt.Fprintln(s.writer, pterm.FgRed.Sprint("✗ "+msg))
}

// Info prints a blue info message
func (s *Styler) Info(msg string) {
	fmt.Fprintln(s.writer, pterm.FgCyan.Sprint("ℹ "+msg))
}

// Warning prints a yellow warning message
func (s *Styler) Warning(msg string) {
	fmt.Fprintln(s.writer, pterm.FgYellow.Sprint("⚠ "+msg))
}

// URL prints a clickable-style URL (underlined cyan)
func (s *Styler) URL(url string) {
	fmt.Fprintln(s.writer, pterm.FgLightCyan.Sprint(url))
}

// Bold prints bold text
func (s *Styler) Bold(msg string) {
	fmt.Fprintln(s.writer, pterm.Bold.Sprint(msg))
}

// Muted prints gray/dim text
func (s *Styler) Muted(msg string) {
	fmt.Fprintln(s.writer, pterm.FgGray.Sprint(msg))
}

// DefaultStyler returns a Styler writing to stdout
func DefaultStyler() *Styler {
	return NewStyler(os.Stdout)
}
