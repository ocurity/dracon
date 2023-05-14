package jira

import (
	"testing"
	"time"

	jira "github.com/andygrunwald/go-jira"
	"github.com/stretchr/testify/assert"
	"github.com/trivago/tgo/tcontainer"

	"github.com/ocurity/dracon/pkg/jira/config"
	"github.com/ocurity/dracon/pkg/jira/document"
)

var (
	sampleConfig = config.Config{
		DefaultValues: config.DefaultValues{
			Project:         "TOY",
			IssueType:       "Vulnerability",
			Components:      []string{"c1", "c2", "c3"},
			AffectsVersions: []string{"V1", "V2"},
			Labels:          []string(nil),
			CustomFields: []config.CustomField{{
				ID:        "customfield_10000",
				FieldType: "multi-value",
				Values:    []string{"foo", "bar"},
			}},
		},
		Mappings: []config.Mappings{{
			DraconField: "cvss",
			JiraField:   "customfield_10001",
			FieldType:   "float",
		}},
		DescriptionExtras: []string{"target", "tool_name"},
	}
	t, _         = time.Parse("0001-01-01T00:00:00Z", "0001-01-01T00:00:00Z")
	sampleResult = document.Document{
		ScanStartTime:  t,
		ScanID:         "babbb83-4627-41c6-8ba0-70ee866290e9",
		ToolName:       "spotbugs",
		Source:         "//foo/bar:baz",
		Target:         "//foo1/bar1:baz2",
		Type:           "test type",
		Title:          "Unit Test Title",
		SeverityText:   "Info",
		CVSS:           "0.000",
		ConfidenceText: "Info",
		Description:    "this is a test description",
		FirstFound:     t,
		FalsePositive:  "true",
		CVE:            "CVE-0000-99999",
	}

	expIssue = jira.Issue{
		Fields: &jira.IssueFields{
			Project: jira.Project{
				Key: "TOY",
			},
			Type: jira.IssueType{
				Name: "Vulnerability",
			},
			Description: "Dracon found 'Unit Test Title'"+
						 " at '//foo1/bar1:baz2',"+
						 " severity 'SEVERITY_INFO',"+
						 " rule id: 'test type',"+
						 " CVSS '0'"+
						 " Confidence 'CONFIDENCE_INFO'"+
						 " Original Description: this is a test description,"+
						 " Cve CVE-0000-99999,\n"+
						 "{code:}\n"+
						 "target:                    //foo1/bar1:baz2\n"+
						 "tool_name:                 spotbugs\n{code}\n",
			Summary: "bar1:baz2 Unit Test Title",
			Components: []*jira.Component{
				{Name: "c1"},
				{Name: "c2"},
				{Name: "c3"},
			},
			AffectsVersions: []*jira.AffectsVersion{
				{Name: "V1"},
				{Name: "V2"},
			},
			Labels: []string(nil),
			Unknowns: tcontainer.MarshalMap{
				"customfield_10000": []map[string]string{{"value": "foo"}, {"value": "bar"}},
				"customfield_10001": 0.000,
			},
		},
	}

	sampleClient = &Client{
		JiraClient:    authJiraClient("test_user", "test_token", "test_url"),
		DryRunMode:    true,
		Config:        sampleConfig,
		DefaultFields: getDefaultFields(sampleConfig),
	}
)

func TestNewClient(t *testing.T) {
	client := NewClient("test_user", "test_token", "test_url", true, sampleConfig)
	assert.EqualValues(t, client, sampleClient)
}

func TestAuthJiraClient(t *testing.T) {
	client := authJiraClient("test_user", "test_token", "test_url")
	assert.NotEmpty(t, client)
}

func TestAssembleIssue(t *testing.T) {
	issue := sampleClient.assembleIssue(sampleResult)
	assert.EqualValues(t, issue, &expIssue)
}

func TestCreateIssue(t *testing.T) {
	err := sampleClient.CreateIssue(sampleResult)
	assert.NoError(t, err)
}
