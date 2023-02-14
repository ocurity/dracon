package npmfullaudit

import (
	"testing"

	v1 "github.com/ocurity/dracon/api/proto/v1"
	atypes "github.com/ocurity/dracon/components/producers/typescript-npm-audit/types"

	"github.com/stretchr/testify/assert"
)

var invalidJSON = `Not a valid JSON object`

func TestParseInvalidJSON(t *testing.T) {
	report, err := NewReport([]byte(invalidJSON))
	assert.Nil(t, report)

	_, ok := err.(*atypes.ParsingError)
	assert.True(t, ok)
}

var invalidAuditReport = `{
	"description": "A valid JSON object, but not a Full Audit report"
}
`

func TestParseInvalidReport(t *testing.T) {
	report, err := NewReport([]byte(invalidAuditReport))
	assert.Nil(t, report)

	_, ok := err.(*atypes.FormatError)
	assert.True(t, ok)
}

var fullAuditReport = `{
	"actions": [{
		"isMajor": true,
		"action": "install",
		"resolves": [{
			"id": 1179,
			"path": "auditjs>npm>bin-links>cmd-shim>mkdirp>bar",
			"dev": false,
			"optional": false,
			"bundled": true
		}]

	}],
	"advisories": {
		"1084": {
			"findings": [{
				"version": "1.1.0",
				"paths": [
					"auditjs>npm>libnpx>yargs>os-locale>foo"
				]
			}],
			"id": 1084,
			"created": "2019-07-18T21:30:31.935Z",
			"updated": "2019-11-19T23:31:37.349Z",
			"deleted": null,
			"title": "Denial of Service",
			"found_by": {
				"link": "",
				"name": "asdfasdfa",
				"email": ""
			},
			"reported_by": {
				"link": "",
				"name": "asdfasfa",
				"email": ""
			},
			"module_name": "foo",
			"cves": [],
			"vulnerable_versions": "<1.2.3",
			"patched_versions": ">=1.2.3",
			"overview": "Versions of 'foo' prior to 1.2.3 are vulnerable to Denial of Service (DoS).  The package fails to remove old values from the cache even after a value passes its 'maxAge' property. This may allow attackers to exhaust the system's memory if they are able to abuse the application logging.",
			"recommendation": "Upgrade to version 1.2.3 or later.",
			"references": "- [Snyk Report](https://snyk.io/vuln/npm:foo:20180117)",
			"access": "public",
			"severity": "low",
			"cwe": "CWE-400",
			"metadata": {
				"module_type": "",
				"exploitability": 2,
				"affected_components": ""
			},
			"url": "https://npmjs.com/advisories/123"
		}
	},
	"muted": [],
	"metadata": {
		"vulnerabilities": {
			"info": 0,
			"low": 136,
			"moderate": 0,
			"high": 23,
			"critical": 0
		},
		"dependencies": 748,
		"devDependencies": 0,
		"optionalDependencies": 5,
		"totalDependencies": 753
	},
	"runId": "75a2f929-191f-4b32-a0a1-28bfd230a36d"
}
`

func TestParseValidReport(t *testing.T) {
	report, err := NewReport([]byte(fullAuditReport))
	report.SetPackagePath("test")
	assert.NoError(t, err)

	issues := report.AsIssues()
	assert.Len(t, issues, 1)

	expectedIssue := &v1.Issue{
		Target:      "test:foo",
		Type:        "Vulnerable Dependency",
		Title:       "Denial of Service",
		Severity:    v1.Severity_SEVERITY_LOW,
		Cvss:        0.0,
		Confidence:  v1.Confidence_CONFIDENCE_HIGH,
		Description: "Vulnerable Versions: <1.2.3\nRecommendation: Upgrade to version 1.2.3 or later.\nOverview: Versions of 'foo' prior to 1.2.3 are vulnerable to Denial of Service (DoS).  The package fails to remove old values from the cache even after a value passes its 'maxAge' property. This may allow attackers to exhaust the system's memory if they are able to abuse the application logging.\nReferences: - [Snyk Report](https://snyk.io/vuln/npm:foo:20180117)\nNPM Advisory URL: https://npmjs.com/advisories/123\n",
	}

	assert.Equal(t, issues[0].Target, expectedIssue.Target)
	assert.Equal(t, issues[0].Type, expectedIssue.Type)
	assert.Equal(t, issues[0].Title, expectedIssue.Title)
	assert.Equal(t, issues[0].Severity, expectedIssue.Severity)
	assert.Equal(t, issues[0].Cvss, expectedIssue.Cvss)
	assert.Equal(t, issues[0].Confidence, expectedIssue.Confidence)
	assert.Equal(t, issues[0].Description, expectedIssue.Description)
}
