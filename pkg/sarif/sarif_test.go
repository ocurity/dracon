package sarif

import (
	"testing"
	"time"

	v1 "github.com/ocurity/dracon/api/proto/v1"
	"github.com/owenrumney/go-sarif/v2/sarif"
	"github.com/stretchr/testify/assert"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_ParseOut(t *testing.T) {
	results, err := sarif.FromString(exampleOutput)
	if err != nil {
		t.Logf(err.Error())
		t.Fail()
	}

	expectedIssues := []*v1.Issue{
		{
			Target:      "main.go:83-83",
			Type:        "G404",
			Title:       "Use of weak random number generator (math/rand instead of crypto/rand)",
			Severity:    v1.Severity_SEVERITY_HIGH,
			Confidence:  v1.Confidence_CONFIDENCE_UNSPECIFIED,
			Description: "Message: Use of weak random number generator (math/rand instead of crypto/rand)",
		},
		{
			Target:      "main.go:347-347",
			Type:        "G104",
			Title:       "Errors unhandled.",
			Severity:    v1.Severity_SEVERITY_MEDIUM,
			Confidence:  v1.Confidence_CONFIDENCE_UNSPECIFIED,
			Description: "Message: Errors unhandled.",
		},
	}
	for _, run := range results.Runs {
		rules := make(map[string]*sarif.ReportingDescriptor)
		for _, rule := range run.Tool.Driver.Rules {
			rules[rule.ID] = rule
		}
		issues := parseOut(*run, rules, "trivy")

		assert.EqualValues(t, expectedIssues, issues)
	}
}

func Test_ToDracon(t *testing.T) {
	var issues []*DraconIssueCollection
	issues, err := ToDracon(trivyOutput)
	if err != nil {
		t.Logf(err.Error())
		t.Fail()
	}

	expectedIssues := []*v1.Issue{
		{
			Target:      "library/ubuntu:1-1",
			Type:        "CVE-2016-20013",
			Title:       "Package: libc6\nInstalled Version: 2.35-0ubuntu3\nVulnerability CVE-2016-20013\nSeverity: LOW\nFixed Version: \nLink: [CVE-2016-20013](https://avd.aquasec.com/nvd/cve-2016-20013)",
			Severity:    v1.Severity_SEVERITY_INFO,
			Confidence:  v1.Confidence_CONFIDENCE_UNSPECIFIED,
			Description: "Message: Package: libc6\nInstalled Version: 2.35-0ubuntu3\nVulnerability CVE-2016-20013\nSeverity: LOW\nFixed Version: \nLink: [CVE-2016-20013](https://avd.aquasec.com/nvd/cve-2016-20013)",
		},
	}
	for _, issue := range issues {
		assert.EqualValues(t, issue.ToolName, "Trivy")
		assert.EqualValues(t, expectedIssues, issue.Issues)
	}
}

func Test_FromDraconEnrichedIssuesRun(t *testing.T) {
	scanStartTime := "2020-04-13T11:51:53Z"
	tstampStart, err := time.Parse(time.RFC3339, scanStartTime)
	assert.Nil(t, err)
	startTime := timestamppb.New(tstampStart)

	tstampFS, err := time.Parse(time.RFC3339, scanStartTime)
	assert.Nil(t, err)

	firstSeen := timestamppb.New(tstampFS)
	tstampUAT, _ := time.Parse(time.RFC3339, scanStartTime)
	updatedAt := timestamppb.New(tstampUAT)

	scanUUID := "babbb83-4627-41c6-8ba0-70ee866290e9"
	responses := []*v1.EnrichedLaunchToolResponse{
		{
			OriginalResults: &v1.LaunchToolResponse{
				ToolName: "test",
				Issues: []*v1.Issue{
					{
						Description: "this is a test description",
						Confidence:  v1.Confidence_CONFIDENCE_INFO,
						Severity:    v1.Severity_SEVERITY_INFO,
						Cvss:        0.0,
						Source:      "//foo/bar:baz",
						Target:      "//foo1/bar1:baz2",
						Title:       "Unit Test Title",
						Type:        "test type",
						Cve:         "CVE-0000-99999",
					},
				},
				ScanInfo: &v1.ScanInfo{
					ScanUuid:      scanUUID,
					ScanStartTime: startTime,
				},
			},
			Issues: []*v1.EnrichedIssue{
				{
					FirstSeen:     firstSeen,
					UpdatedAt:     updatedAt,
					Hash:          "cf23df2207d99a74fbe169e3eba035e633b65d94",
					FalsePositive: true,
					Count:         2,
					Annotations: map[string]string{
						"Policy":    "fail",
						"footag":    "bartag",
						"Signature": "blah",
					},
					RawIssue: &v1.Issue{
						Description: "this is a test description",
						Confidence:  v1.Confidence_CONFIDENCE_INFO,
						Severity:    v1.Severity_SEVERITY_INFO,
						Cvss:        0.0,
						Source:      "//foo/bar:baz",
						Target:      "//foo1/bar1:baz2",
						Title:       "Unit Test Title",
						Type:        "test type",
						Cve:         "CVE-0000-99999",
					},
				},
				{
					FirstSeen:     firstSeen,
					UpdatedAt:     updatedAt,
					Hash:          "cf23df2207d99a74fbe169e3eba035e633b65d94",
					FalsePositive: false,
					Count:         1,
					Annotations: map[string]string{
						"Policy":    "fail",
						"footag":    "bartag",
						"Signature": "blah",
					},
					RawIssue: &v1.Issue{
						Description: "this is a test description",
						Confidence:  v1.Confidence_CONFIDENCE_INFO,
						Severity:    v1.Severity_SEVERITY_INFO,
						Cvss:        0.0,
						Source:      "//foo/bar:baz",
						Target:      "//foo1/bar1:baz2",
						Title:       "Unit Test Title",
						Type:        "test type",
						Cve:         "CVE-0000-99999",
					},
				},
			},
		},
	}
	typ := "test type"
	level := "note"
	msg := "this is a test description"
	uri := "//foo1/bar1:baz2"
	ruleIndex := uint(0)
	fp := "False Positive:false"
	hash := "Hash:cf23df2207d99a74fbe169e3eba035e633b65d94"
	policy := "Policy:fail"
	foo := "footag:bartag"
	sign := "Signature:blah"
	confidence := "Confidence:CONFIDENCE_INFO"
	source := "Source://foo/bar:baz"
	cve := "CVE-0000-99999"

	expected := &sarif.Report{
		Version: "2.1.0",
		Schema:  "https://json.schemastore.org/sarif-2.1.0-rtm.5.json",
		Runs: []*sarif.Run{
			{
				AutomationDetails: &sarif.RunAutomationDetails{
					Description: &sarif.Message{
						Text: &scanStartTime,
					},
					GUID: &scanUUID,
					ID:   &scanUUID,
				},
				Tool: sarif.Tool{
					Driver: &sarif.ToolComponent{
						Name: "test",
						Rules: []*sarif.ReportingDescriptor{{
							ID: typ,
						}},
					},
				},
				Results: []*sarif.Result{
					{
						RuleID:    &typ,
						RuleIndex: &ruleIndex,
						Level:     &level,
						Message: sarif.Message{
							Text: &msg,
						},
						Locations: []*sarif.Location{
							{
								PhysicalLocation: &sarif.PhysicalLocation{
									ArtifactLocation: &sarif.ArtifactLocation{
										URI: &uri,
									},
								},
							},
						},
						Attachments: []*sarif.Attachment{
							{Description: &sarif.Message{Text: &confidence}},
							{Description: &sarif.Message{Text: &source}},
							{Description: &sarif.Message{Text: &cve}},
							{Description: &sarif.Message{Text: &fp}},
							{Description: &sarif.Message{Text: &hash}},
							{Description: &sarif.Message{Text: &policy}},
							{Description: &sarif.Message{Text: &foo}},
							{Description: &sarif.Message{Text: &sign}},
						},
					},
				},
			},
		},
	}

	report, err := FromDraconEnrichedIssuesRun(responses, false)
	assert.Nil(t, err)
	assert.NotNil(t, report)
	assert.EqualValues(t, report, expected)
}

func Test_FromDraconRawIssuesRun(t *testing.T) {
	scanStartTime := "2020-04-13T11:51:53Z"
	tstampStart, err := time.Parse(time.RFC3339, scanStartTime)
	assert.Nil(t, err)

	startTime := timestamppb.New(tstampStart)
	scanUUID := "babbb83-4627-41c6-8ba0-70ee866290e9"
	responses := []*v1.LaunchToolResponse{
		{
			ToolName: "test",
			Issues: []*v1.Issue{
				{
					Description: "this is a test description",
					Confidence:  v1.Confidence_CONFIDENCE_INFO,
					Severity:    v1.Severity_SEVERITY_INFO,
					Cvss:        0.0,
					Source:      "//foo/bar:baz",
					Target:      "//foo1/bar1:baz2",
					Title:       "Unit Test Title",
					Type:        "test type",
					Cve:         "CVE-0000-99999",
				},
			},
			ScanInfo: &v1.ScanInfo{
				ScanUuid:      scanUUID,
				ScanStartTime: startTime,
			},
		},
	}

	typ := "test type"
	level := "note"
	msg := "this is a test description"
	uri := "//foo1/bar1:baz2"
	ruleIndex := uint(0)
	confidence := "Confidence:CONFIDENCE_INFO"
	source := "Source://foo/bar:baz"
	cve := "CVE-0000-99999"
	expected := &sarif.Report{
		Version: "2.1.0",
		Schema:  "https://json.schemastore.org/sarif-2.1.0-rtm.5.json",
		Runs: []*sarif.Run{
			{
				AutomationDetails: &sarif.RunAutomationDetails{
					Description: &sarif.Message{
						Text: &scanStartTime,
					},
					GUID: &scanUUID,
					ID:   &scanUUID,
				},
				Tool: sarif.Tool{
					Driver: &sarif.ToolComponent{
						Name: "test",
						Rules: []*sarif.ReportingDescriptor{{
							ID: typ,
						}},
					},
				},
				Results: []*sarif.Result{
					{
						RuleID:    &typ,
						RuleIndex: &ruleIndex,
						Level:     &level,
						Message: sarif.Message{
							Text: &msg,
						},
						Attachments: []*sarif.Attachment{
							{Description: &sarif.Message{Text: &confidence}},
							{Description: &sarif.Message{Text: &source}},
							{Description: &sarif.Message{Text: &cve}},
						},
						Locations: []*sarif.Location{
							{
								PhysicalLocation: &sarif.PhysicalLocation{
									ArtifactLocation: &sarif.ArtifactLocation{
										URI: &uri,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	report, err := FromDraconRawIssuesRun(responses)
	assert.Nil(t, err)
	assert.NotNil(t, report)
	assert.EqualValues(t, report, expected)
}

func Test_draconIssueToSarif(t *testing.T) {
	issue := &v1.Issue{
		Description: "this is a test description",
		Confidence:  v1.Confidence_CONFIDENCE_INFO,
		Severity:    v1.Severity_SEVERITY_INFO,
		Cvss:        0.0,
		Source:      "//foo/bar:baz",
		Target:      "/workspace/source-code-ws/foo1/bar1:baz2",
		Title:       "Unit Test Title",
		Type:        "test type",
		Cve:         "CVE-0000-99999",
	}
	typ := "test type"
	level := "note"
	msg := "this is a test description"
	uri := "/foo1/bar1:baz2"
	confidence := "Confidence:CONFIDENCE_INFO"
	source := "Source://foo/bar:baz"
	cve := "CVE-0000-99999"
	expected := &sarif.Result{
		RuleID:  &typ,
		Level:   &level,
		Message: sarif.Message{Text: &msg},
		Attachments: []*sarif.Attachment{
			{Description: &sarif.Message{Text: &confidence}},
			{Description: &sarif.Message{Text: &source}},
			{Description: &sarif.Message{Text: &cve}},
		},
		Locations: []*sarif.Location{
			{
				PhysicalLocation: &sarif.PhysicalLocation{
					ArtifactLocation: &sarif.ArtifactLocation{
						URI: &uri,
					},
				},
			},
		},
	}

	sarifResults, err := draconIssueToSarif(issue, &sarif.ReportingDescriptor{ID: typ})
	assert.Nil(t, err)
	assert.NotNil(t, sarifResults)
	assert.EqualValues(t, sarifResults, expected)
}

var trivyOutput = `{
	"version": "2.1.0",
	"$schema": "https://json.schemastore.org/sarif-2.1.0-rtm.5.json",
	"runs": [
	  {
		"tool": {
		  "driver": {
			"fullName": "Trivy Vulnerability Scanner",
			"informationUri": "https://github.com/aquasecurity/trivy",
			"name": "Trivy",
			"version": "0.29.2"
		  }
		},
		"results": [
		  {
			"ruleId": "CVE-2016-20013",
			"ruleIndex": 3,
			"level": "note",
			"message": {
			  "text": "Package: libc6\nInstalled Version: 2.35-0ubuntu3\nVulnerability CVE-2016-20013\nSeverity: LOW\nFixed Version: \nLink: [CVE-2016-20013](https://avd.aquasec.com/nvd/cve-2016-20013)"
			},
			"locations": [
			  {
				"physicalLocation": {
				  "artifactLocation": {
					"uri": "library/ubuntu",
					"uriBaseId": "ROOTPATH"
				  },
				  "region": {
					"startLine": 1,
					"startColumn": 1,
					"endLine": 1,
					"endColumn": 1
				  }
				}
			  }
			]
		  }
		],
		"columnKind": "utf16CodeUnits",
		"originalUriBaseIds": {
		  "ROOTPATH": {
			"uri": "file:///"
		  }
		}
	  }
	]
  }`

var exampleOutput = `{
	"runs": [{
		"results": [{
				"level": "error",
				"locations": [{
					"physicalLocation": {
						"artifactLocation": {
							"uri": "main.go"
						},
						"region": {
							"endColumn": 7,
							"endLine": 83,
							"snippet": {
								"text": "r := rand.New(rand.NewSource(time.Now().UnixNano()))"
							},
							"sourceLanguage": "go",
							"startColumn": 7,
							"startLine": 83
						}
					}
				}],
				"message": {
					"text": "Use of weak random number generator (math/rand instead of crypto/rand)"
				},
				"ruleId": "G404"
			},
			{
				"level": "warning",
				"locations": [{
					"physicalLocation": {
						"artifactLocation": {
							"uri": "main.go"
						},
						"region": {
							"endColumn": 2,
							"endLine": 347,
							"snippet": {
								"text": "zipWriter.Close()"
							},
							"sourceLanguage": "go",
							"startColumn": 2,
							"startLine": 347
						}
					}
				}],
				"message": {
					"text": "Errors unhandled."
				},
				"ruleId": "G104",
				"ruleIndex": 3
			}
		],
		"tool": {
			"driver": {
				"guid": "8b518d5f-906d-39f9-894b-d327b1a421c5",
				"informationUri": "https://github.com/securego/gosec/",
				"name": "gosec"
			}
		}
	}],
	"$schema": "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
	"version": "2.1.0"
}`
