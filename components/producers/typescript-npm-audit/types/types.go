// Package types provides common types for audit report formats.
package types

import (
	v1 "github.com/ocurity/dracon/api/proto/v1"
)

// Report is an interface for audit report formats.
type Report interface {
	// SetPackagePath sets the path to the package to which this audit report
	// belongs.
	SetPackagePath(string)

	// Type returns a short textual description of this report type.
	Type() string

	// AsIssues transforms this Report into a slice of Dracon Issues that can be
	// processed by the Dracon enricher.
	AsIssues() []*v1.Issue
}
