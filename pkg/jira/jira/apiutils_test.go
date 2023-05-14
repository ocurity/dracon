package jira

import (
	"testing"

	jira "github.com/andygrunwald/go-jira"
	"github.com/stretchr/testify/assert"
	"github.com/trivago/tgo/tcontainer"
)

func TestSetDefaultFields(t *testing.T) {
	res := getDefaultFields(sampleConfig)

	exp := defaultJiraFields{
		Project: jira.Project{
			Key: "TOY",
		},
		IssueType: jira.IssueType{
			Name: "Vulnerability",
		},
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
		CustomFields: tcontainer.MarshalMap{
			"customfield_10000": []map[string]string{{"value": "foo"}, {"value": "bar"}},
		},
	}

	assert.EqualValues(t, res, exp)
}

func TestMakeCustomField(t *testing.T) {
	res1 := makeCustomField("single-value", []string{"test-value"})
	exp1 := map[string]string{"value": "test-value"}

	res2 := makeCustomField("multi-value", []string{"value1", "value2", "value3"})
	exp2 := []map[string]string{
		{"value": "value1"},
		{"value": "value2"},
		{"value": "value3"},
	}

	res3 := makeCustomField("float", []string{"4.22"})
	exp3 := 4.22

	res4 := makeCustomField("simple-value", []string{"test-value"})
	exp4 := "test-value"

	assert.EqualValues(t, res1, exp1)
	assert.EqualValues(t, res2, exp2)
	assert.Equal(t, res3, exp3)
	assert.Equal(t, res4, exp4)
}

func TestMakeDescription(t *testing.T) {
	extras := []string{"tool_name", "target", "confidence_text"}
	res := makeDescription(sampleResult, extras,      "")
	exp :="Dracon found 'Unit Test Title' at '//foo1/bar1:baz2',"+
		" severity 'SEVERITY_INFO',"+
		" rule id: 'test type',"+
		" CVSS '0' "+
		"Confidence 'CONFIDENCE_INFO'"+
		" Original Description: this is a test description,"+
		" Cve CVE-0000-99999,\n"+
		"{code:}\n"+
		"tool_name:                 spotbugs\n"+
		"target:                    //foo1/bar1:baz2\n"+
		"confidence_text:           Info\n"+
		"{code}\n"
	assert.Equal(t, res, exp)
}

func TestMakeSummary(t *testing.T) {
	res, extra := makeSummary(sampleResult)
	exp := "bar1:baz2 Unit Test Title"
	assert.Equal(t, res, exp)
	assert.Equal(t, extra, "")

	longTitle := make([]rune, 300)
	truncatedSummary := make([]rune, 254)
	expectedExtra := make([]rune, 49)
	char := []rune{'\u0061'} // 'a'
	for i := range longTitle {
		longTitle[i] = char[0]
	}

	// truncatedSummary = []rune{'b','a','r',' '}
	for i := 4; i < len(truncatedSummary); i++ {
		truncatedSummary[i] = char[0]
	}
	truncatedSummary[0] = 'b'
	truncatedSummary[1] = 'a'
	truncatedSummary[2] = 'r'
	truncatedSummary[3] = ' '

	for i := range expectedExtra {
		expectedExtra[i] = char[0]
	}
	sampleResult.Target = "/foo/bar"
	sampleResult.Title = string(longTitle)

	res, extra = makeSummary(sampleResult)
	assert.Equal(t, string(truncatedSummary), res)

	assert.Equal(t, string(expectedExtra), extra)
}
