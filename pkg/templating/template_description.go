// Package templating includes helper methods that apply
// go templates to Dracon Raw and Enriched Issues and return the resulting str
package templating

import (
	"bytes"
	"text/template"
	"time"

	"github.com/go-errors/errors"

	v1 "github.com/ocurity/dracon/api/proto/v1"
)

const (
	defaultEnrichedFindingTemplate = "Dracon found '{{.RawIssue.Title}}' at '{{.RawIssue.Target}}', severity '{{.RawIssue.Severity}}', rule id: '{{.RawIssue.Type}}', CVSS '{{.RawIssue.Cvss}}' Confidence '{{.RawIssue.Confidence}}' Original Description: {{.RawIssue.Description}}, Cve {{.RawIssue.Cve}},\n{{ range $key,$element := .Annotations }}{{$key}}:{{$element}}\n{{end}}"
	defaultRawFindingTemplate      = "Dracon found '{{.Title}}' at '{{.Target}}', severity '{{.Severity}}', rule id: '{{.Type}}', CVSS '{{.Cvss}}' Confidence '{{.Confidence}}' Original Description: {{.Description}}, Cve {{.Cve}}"
)

// TemplateStringRaw applies the provided go template to the Raw Issue provided and returns the resulting str
func TemplateStringRaw(inputTemplate string, issue *v1.Issue) (*string, error) {
	if inputTemplate == "" {
		inputTemplate = defaultRawFindingTemplate
	}
	tmpl, err := template.New("description").Parse(inputTemplate)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)

	err = tmpl.Execute(buf, issue)
	if err != nil {
		return nil, err
	}
	res := buf.String()
	return &res, nil
}

type enrichedIssue struct {
	*v1.EnrichedIssue
	ToolName       string
	ScanStartTime  string
	ScanID         string
	ConfidenceText string
	SeverityText   string
	Count          uint
	FirstFound     string
}
type enrichedIssueOption func(*enrichedIssue) error

func EnrichedIssueWithToolName(toolname string) enrichedIssueOption {
	return func(ei *enrichedIssue) error {
		if toolname == "" {
			return errors.New("invalid toolname ''")
		}
		ei.ToolName = toolname
		return nil
	}
}

func EnrichedIssueWithScanStartTime(startTime time.Time) enrichedIssueOption {
	return func(ei *enrichedIssue) error {
		if time.Time.IsZero(startTime) {
			return errors.New("invalid startTime zero")
		}
		ei.ScanStartTime = startTime.Format(time.RFC3339)
		return nil
	}
}

func EnrichedIssueWithConfidenceText(confidence string) enrichedIssueOption {
	return func(ei *enrichedIssue) error {
		if confidence == "" {
			return errors.New("invalid confidence ''")
		}
		ei.ConfidenceText = confidence
		return nil
	}
}

func EnrichedIssueWithSeverityText(severity string) enrichedIssueOption {
	return func(ei *enrichedIssue) error {
		if severity == "" {
			return errors.New("invalid severity ''")
		}
		ei.SeverityText = severity
		return nil
	}
}

func EnrichedIssueWithCount(count uint) enrichedIssueOption {
	return func(ei *enrichedIssue) error {
		if count <= 0 {
			return errors.Errorf("invalid count %d", count)
		}
		ei.Count = count
		return nil
	}
}

func EnrichedIssueWithScanID(scanID string) enrichedIssueOption {
	return func(ei *enrichedIssue) error {
		if scanID == "" {
			return errors.New("invalid scanID ")
		}
		ei.ScanID = scanID
		return nil
	}
}

func EnrichedIssueWithFirstFound(firstFound time.Time) enrichedIssueOption {
	return func(ei *enrichedIssue) error {
		if time.Time.IsZero(firstFound) {
			return errors.New("invalid firstFound zero")
		}
		ei.FirstFound = firstFound.Format(time.RFC3339)
		return nil
	}
}

// TemplateStringEnriched applies the provided go template to the Enriched Issue provided and returns the resulting str
func TemplateStringEnriched(inputTemplate string, issue *v1.EnrichedIssue, opts ...enrichedIssueOption) (*string, error) {
	enrichedIssue := &enrichedIssue{
		EnrichedIssue: issue,
	}
	for _, opt := range opts {
		if err := opt(enrichedIssue); err != nil {
			return nil, errors.Errorf("could not apply enriched issue option: %w", err)
		}
	}
	if inputTemplate == "" {
		inputTemplate = defaultEnrichedFindingTemplate
	}
	tmpl, err := template.New("description").Parse(inputTemplate)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)

	err = tmpl.Execute(buf, enrichedIssue)
	if err != nil {
		return nil, err
	}
	res := buf.String()
	return &res, nil
}
