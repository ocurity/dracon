package types

import (
	"fmt"
)

// ParsingError indicates that a given string could not be parsed as a
// particular audit report type.
type ParsingError struct {
	Type          string
	PrintableType string
	Err           error
}

func (e *ParsingError) Error() string {
	return fmt.Sprintf("failed to parse input as %s", e.PrintableType)
}

func (e *ParsingError) Unwrap() error {
	return e.Err
}

// FormatError indicates that a given string could be parsed as a particular
// audit report type, but does not satisfy the required format for that type of
// audit report.
type FormatError struct {
	Type          string
	PrintableType string
}

func (e *FormatError) Error() string {
	return fmt.Sprintf("input is not %s", e.PrintableType)
}
