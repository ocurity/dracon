package templating

import (
	"testing"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	v1 "github.com/ocurity/dracon/api/proto/v1"
	"github.com/stretchr/testify/assert"
)

func Test_TemplateStringRaw(t *testing.T) {
	type args struct {
		inputTemplate string
		issue         v1.Issue
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "template references some of of the issue",
			args: args{
				inputTemplate: "Dracon found '{{.Title}}' at '{{.Target}}', severity '{{.Severity}}'",
				issue: v1.Issue{
					Target:      "/foo/bar/baz:32",
					Title:       "whoops, XSS!",
					Severity:    v1.Severity_SEVERITY_HIGH,
					Cvss:        3.2,
					Confidence:  v1.Confidence_CONFIDENCE_UNSPECIFIED,
					Description: "this is a description",
					Source:      "github.com/foo/bar",
					Cve:         "CVE-2020-1111",
					Uuid:        "2d7d1bd6-f1a0-11ed-a05b-0242ac120003",
				},
			},
			want:    "Dracon found 'whoops, XSS!' at '/foo/bar/baz:32', severity 'SEVERITY_HIGH'",
			wantErr: false,
		},
		{
			name: "template references all of the issue",
			args: args{
				inputTemplate: "Dracon found '{{.Title}}' at '{{.Target}}', severity '{{.Severity}}', rule id: '{{.Type}}', CVSS '{{.Cvss}}' Confidence '{{.Confidence}}' Original Description: {{.Description}}, Cve {{.Cve}}",
				issue: v1.Issue{
					Target:      "/foo/bar/baz:32",
					Title:       "whoops, XSS!",
					Severity:    v1.Severity_SEVERITY_HIGH,
					Cvss:        3.2,
					Confidence:  v1.Confidence_CONFIDENCE_UNSPECIFIED,
					Description: "this is a description",
					Source:      "github.com/foo/bar",
					Cve:         "CVE-2020-1111",
					Uuid:        "2d7d1bd6-f1a0-11ed-a05b-0242ac120003",
				},
			},
			want:    "Dracon found 'whoops, XSS!' at '/foo/bar/baz:32', severity 'SEVERITY_HIGH', rule id: '', CVSS '3.2' Confidence 'CONFIDENCE_UNSPECIFIED' Original Description: this is a description, Cve CVE-2020-1111",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TemplateStringRaw(tt.args.inputTemplate, &tt.args.issue)
			if (err != nil) != tt.wantErr {
				t.Errorf("templateStringRaw() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if *got != tt.want {
				t.Errorf("templateStringRaw() = '%s', want '%s'", *got, tt.want)
			}
		})
	}
}

func Test_TemplateStringEnriched(t *testing.T) {
	tstampFS, err := time.Parse(time.RFC3339, "2020-04-13T11:51:53+01:00")
	assert.Nil(t, err)
	firstSeen := timestamppb.New(tstampFS)
	tstampUAT, err := time.Parse(time.RFC3339, "2020-04-13T11:51:53+01:00")
	assert.Nil(t, err)
	updatedAt := timestamppb.New(tstampUAT)

	type args struct {
		inputTemplate string
		issue         v1.EnrichedIssue
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "template references all of the issue",
			args: args{
				inputTemplate: "Dracon found '{{.RawIssue.Title}}' at '{{.RawIssue.Target}}', severity '{{.RawIssue.Severity}}', rule id: '{{.RawIssue.Type}}', CVSS '{{.RawIssue.Cvss}}' Confidence '{{.RawIssue.Confidence}}' Original Description: {{.RawIssue.Description}}, Cve {{.RawIssue.Cve}},\n{{ range $key,$element := .Annotations }}{{$key}}:{{$element}}\n{{end}}",
				issue: v1.EnrichedIssue{
					RawIssue: &v1.Issue{
						Target:      "/foo/bar/baz:32",
						Title:       "whoops, XSS!",
						Severity:    v1.Severity_SEVERITY_HIGH,
						Cvss:        3.2,
						Confidence:  v1.Confidence_CONFIDENCE_HIGH,
						Description: "this is a description",
						Source:      "github.com/foo/bar",
						Cve:         "CVE-2020-1111",
						Type:        "G101",
						Uuid:        "2d7d1bd6-f1a0-11ed-a05b-0242ac120003",
					},
					FirstSeen:     firstSeen,
					Count:         15,
					FalsePositive: false,
					UpdatedAt:     updatedAt,
					Hash:          "",
					Annotations:   map[string]string{"Policy X": "false", "Some Other Annotation": "value"},
				},
			},
			want:    "Dracon found 'whoops, XSS!' at '/foo/bar/baz:32', severity 'SEVERITY_HIGH', rule id: 'G101', CVSS '3.2' Confidence 'CONFIDENCE_HIGH' Original Description: this is a description, Cve CVE-2020-1111,\nPolicy X:false\nSome Other Annotation:value\n",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TemplateStringEnriched(tt.args.inputTemplate, &tt.args.issue)
			if (err != nil) != tt.wantErr {
				t.Errorf("templateStringEnriched() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if *got != tt.want {
				t.Errorf("templateStringEnriched() = `%s`, want `%s`", *got, tt.want)
			}
		})
	}
}
