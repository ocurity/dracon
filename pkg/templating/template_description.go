// Package templating includes helper methods that apply
// go templates to Dracon Raw and Enriched Issues and return the resulting str
package templating

import (
	"bytes"
	"text/template"

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

// TemplateStringEnriched applies the provided go template to the Enriched Issue provided and returns the resulting str
func TemplateStringEnriched(inputTemplate string, issue *v1.EnrichedIssue) (*string, error) {
	if inputTemplate == "" {
		inputTemplate = defaultEnrichedFindingTemplate
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
