package output

import "github.com/pterm/pterm"

// Spinner wraps pterm spinner for CLI progress indication
type Spinner struct {
	spinner *pterm.SpinnerPrinter
}

// NewSpinner creates and starts a spinner with the given text
func NewSpinner(text string) *Spinner {
	sp, _ := pterm.DefaultSpinner.
		WithRemoveWhenDone(false).
		Start(text)
	return &Spinner{spinner: sp}
}

// UpdateText changes the spinner text
func (s *Spinner) UpdateText(text string) {
	s.spinner.UpdateText(text)
}

// Success stops the spinner with a success message
func (s *Spinner) Success(msg string) {
	s.spinner.Success(msg)
}

// Fail stops the spinner with a failure message
func (s *Spinner) Fail(msg string) {
	s.spinner.Fail(msg)
}

// Warning stops the spinner with a warning message
func (s *Spinner) Warning(msg string) {
	s.spinner.Warning(msg)
}

// Stop stops the spinner without message
func (s *Spinner) Stop() {
	_ = s.spinner.Stop()
}
