package jira

import (
	"testing"
	"time"

	jira "github.com/andygrunwald/go-jira"
	"github.com/ocurity/dracon/pkg/jira/document"
	"github.com/stretchr/testify/require"
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

	require.EqualValues(t, res, exp)
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

	require.EqualValues(t, res1, exp1)
	require.EqualValues(t, res2, exp2)
	require.Equal(t, res3, exp3)
	require.Equal(t, res4, exp4)
}

func TestMakeDescription(t *testing.T) {
	res := makeDescription(sampleResult, "")
	exp := "spotbugs detected 'Unit Test Title' at //foo1/bar1:baz2 during scan with id babbb83-4627-41c6-8ba0-70ee866290e9.\nConfidence: Info\nThis issue has been detected 2 times before\nOriginal Description is: 'this is a test description'\nspotbugs reported severity as Info\nSmithy enrichers added the following annotations:\nfoo:bar\nfoobar:baz\n\n"
	require.Equal(t, res, exp)
}


func TestMakeSummary(t *testing.T) {
	res, extra := makeSummary(sampleResult)
	exp := "bar1:baz2 Unit Test Title"
	require.Equal(t, res, exp)
	require.Equal(t, extra, "")

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
	require.Equal(t, string(truncatedSummary), res)

	require.Equal(t, string(expectedExtra), extra)
}
