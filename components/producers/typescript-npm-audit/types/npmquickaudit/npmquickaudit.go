// Package npmquickaudit provides types and functions for working with audit
// reports from npm's "Quick Audit" endpoint (/-/npm/v1/security/audits/quick)
// and transforming them into data structures understood by the Dracon enricher.
// These reports are JSON objects organised by vulnerable package name; they do
// not contain as much information about the vulnerabilities affecting each
// package as npm Full Audit reports (hence the name).
package npmquickaudit

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"

	v1 "github.com/ocurity/dracon/api/proto/v1"
	"github.com/ocurity/dracon/components/producers"
	atypes "github.com/ocurity/dracon/components/producers/typescript-npm-audit/types"
)

// PrintableType package const, printed as part of the report or errors.
const PrintableType = "npm Quick Audit report"

// Report represents an npm Quick Audit report. The key for Vulnerabilities
// represents a package name.
type Report struct {
	PackagePath     string                   `json:"-"`
	Version         int                      `json:"auditReportVersion"`
	Vulnerabilities map[string]Vulnerability `json:"vulnerabilities"`
}

// Vulnerability represents the set of vulnerabilities present in a particular
// package.
type Vulnerability struct {
	Package  string     `json:"name"`
	Severity string     `json:"severity"`
	Via      []Advisory `json:"via"`
	Effects  []string   `json:"effects"`
	Range    string     `json:"range"`
	Fix      Fix        `json:"fixAvailable"`
}

// Advisory represents a single vulnerability in a particular package. This
// vulnerability may arise either in this package itself (non-transitive), or
// because this package depends on a vulnerable package described elsewhere in
// the report (transitive). For transitive advisories, only the Transitive,
// Package and Dependency fields have values assigned.
type Advisory struct {
	Transitive bool   `json:"-"`
	ID         int    `json:"source"`
	Package    string `json:"name"`
	Dependency string `json:"dependency"`
	Title      string `json:"title"`
	URL        string `json:"url"`
	Severity   string `json:"severity"`
	Range      string `json:"range"`
}

// UnmarshalJSON converts NPM Audit JSON results to Advisory structs.
func (a *Advisory) UnmarshalJSON(data []byte) error {
	// An advisory in the audit report is either a string containing a package
	// name (which exists elsewhere in the report), or an object representing an
	// advisory from the npm Registry
	var packageName string
	if err := json.Unmarshal(data, &packageName); err == nil {
		*a = Advisory{
			Transitive: true,
			Package:    packageName,
			Dependency: packageName,
		}

		return nil
	}

	type tempAdvisory Advisory
	advisory := tempAdvisory{}
	if err := json.Unmarshal(data, &advisory); err == nil {
		newAdvisory := Advisory(advisory)
		newAdvisory.Transitive = false
		*a = newAdvisory

		return nil
	}

	return &json.UnmarshalTypeError{
		Value: "Advisory",
		Type:  reflect.TypeOf(a),
	}
}

// Fix represents a proposed fix for a particular advisory.
type Fix struct {
	Available bool
	Package   string `json:"name"`
	Version   string `json:"version"`
	IsMajor   bool   `json:"isSemVerMajor"`
}

// UnmarshalJSON transforms between NPM Audit fix json and the Fix struct above.
func (f *Fix) UnmarshalJSON(data []byte) error {
	var isAvailable bool
	if err := json.Unmarshal(data, &isAvailable); err == nil {
		*f = Fix{
			Available: isAvailable,
		}

		return nil
	}

	type tempFix Fix
	fix := tempFix{}
	if err := json.Unmarshal(data, &fix); err == nil {
		newFix := Fix(fix)
		newFix.Available = true
		*f = newFix

		return nil
	}

	return &json.UnmarshalTypeError{
		Value: "Fix",
		Type:  reflect.TypeOf(f),
	}
}

// NewReport constructs a Report from an npm Full Audit report.
func NewReport(report []byte) (atypes.Report, error) {
	var r *Report
	if err := producers.ParseJSON(report, &r); err != nil {
		//nolint:errorlint
		switch err.(type) {
		case *json.SyntaxError,
			*json.UnmarshalTypeError, *json.UnsupportedTypeError, *json.UnsupportedValueError:
			return nil, &atypes.ParsingError{
				Type:          "npm_quick_audit",
				PrintableType: PrintableType,
				Err:           err,
			}
		default:
			return nil, err
		}
	}

	if r.Version != 2 {
		return nil, &atypes.FormatError{
			Type:          "npm_quick_audit",
			PrintableType: PrintableType,
		}
	}

	return r, nil
}

// SetPackagePath helper method to set the npm package path.
func (r *Report) SetPackagePath(packagePath string) {
	r.PackagePath = packagePath
}

// Type helper method to set the type.
func (r *Report) Type() string {
	return PrintableType
}

// AsIssues transforms between NPM issues and dracon issues.
func (r *Report) AsIssues() []*v1.Issue {
	issues := make([]*v1.Issue, 0)

	for _, vuln := range r.Vulnerabilities {
		for _, a := range vuln.Via {
			if a.Transitive {
				continue
			}

			var description string
			aData, err := NewAdvisoryData(a.URL)
			if err == nil {
				description = fmt.Sprintf("Vulnerable versions: %s\nRecommendation: %s\nOverview: %s\nReferences: %s\n",
					aData.VulnerableVersions, aData.Recommendation, aData.Overview, aData.References)
			} else {
				log.Printf("Failed to fetch NPM advisory data from %s (%v); issue will not contain advisory data\n",
					a.URL, err)
			}
			description += fmt.Sprintf("NPM advisory URL: %s\n", a.URL)

			var targetName string
			if r.PackagePath != "" {
				targetName = r.PackagePath + ":"
			}
			targetName += a.Package

			issues = append(issues, &v1.Issue{
				Target:      targetName,
				Type:        "Vulnerable Dependency",
				Title:       a.Title,
				Severity:    v1.Severity(v1.Severity_value[fmt.Sprintf("SEVERITY_%s", strings.ToUpper(a.Severity))]),
				Confidence:  v1.Confidence_CONFIDENCE_HIGH,
				Description: description,
			})
		}
	}

	return issues
}
