package templating

import (
	"bytes"
	"fmt"
	"text/template"

	v1 "github.com/ocurity/dracon/api/proto/v1"
)

func TemplateStringRaw(inputTemplate string, issue *v1.Issue) (*string, error) {
	if inputTemplate == "" {
		return nil, fmt.Errorf("you need to specify a template")
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
func TemplateStringEnriched(inputTemplate string, issue *v1.EnrichedIssue) (*string, error) {
	if inputTemplate == "" {
		return nil, fmt.Errorf("you need to specify a template")
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
