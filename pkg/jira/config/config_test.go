package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var sampleConfig = Config{
	DefaultValues: DefaultValues{
		Project:         "TEST",
		IssueType:       "Task",
		Components:      []string{"c1", "c2", "c3"},
		AffectsVersions: []string{"V1", "V2"},
		Labels:          []string(nil),
		CustomFields: []CustomField{{
			ID:        "customfield_10000",
			FieldType: "multi-value",
			Values:    []string{"foo", "bar"},
		}},
	},
	Mappings: []Mappings{{
		DraconField: "cvss",
		JiraField:   "customfield_10001",
		FieldType:   "float",
	}},
	DescriptionExtras: []string{
		"scan_start_time",
		"tool_name",
		"target",
		"type",
		"confidence_text",
		"annotations",
		"hash",
	},
	SyncMappings: []JiraToDraconVulnMappings{
		{
			JiraStatus:     "Test",
			JiraResolution: "Test",
			DraconStatus:   "Resolved",
		},
	},
}

func TestGetConfig(t *testing.T) {
	testConfig := `
	{
		"defaultValues": {
			"project": "TEST",
			"issueType": "Task",
			"customFields": [{
				"ID": "customfield_10000",
				"fieldType": "multi-value",
				"values": [
					"foo",
					"bar"
				]
			}],
			"components": [
				"c1",
				"c2",
				"c3"
			],
			"affectsVersions": [
				"V1",
				"V2"
			]
		},
		"addToDescription": [
			"scan_start_time",
			"tool_name",
			"target",
			"type",
			"confidence_text",
			"annotations",
			"hash"
		],
		"mappings": [{
			"draconField": "cvss",
			"jiraField": "customfield_10001",
			"fieldType": "float"
		}],
		"syncMappings": [{
			"jiraStatus": "Test",
			"jiraResolution": "Test",
			"draconStatus": "Resolved"
		}]
	} 
`
	reader := strings.NewReader(testConfig)
	res, err := New(reader)
	assert.Nil(t, err)
	assert.EqualValues(t, res, sampleConfig)
}
