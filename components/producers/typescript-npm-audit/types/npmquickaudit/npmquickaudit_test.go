package npmquickaudit

import (
	"testing"

	v1 "github.com/ocurity/dracon/api/proto/v1"
	atypes "github.com/ocurity/dracon/components/producers/typescript-npm-audit/types"

	"github.com/stretchr/testify/assert"
	gock "gopkg.in/h2non/gock.v1"
)

var invalidJSON = `Not a valid JSON object`

func TestNewReportInvalidJSON(t *testing.T) {
	report, err := NewReport([]byte(invalidJSON))
	assert.Nil(t, report)

	_, ok := err.(*atypes.ParsingError)
	assert.True(t, ok)
}

var invalidAuditReportJSON = `{
	"description": "A valid JSON object, but not a Quick Audit report"
}
`

func TestNewReportInvalidReport(t *testing.T) {
	report, err := NewReport([]byte(invalidAuditReportJSON))
	assert.Nil(t, report)

	_, ok := err.(*atypes.FormatError)
	assert.True(t, ok)
}

var quickAuditReportJSON = `{
  "auditReportVersion": 2,
  "vulnerabilities": {
	"fbjs": {
	  "name": "fbjs",
	  "severity": "low",
	  "via": ["isomorphic-fetch"],
	  "effects": [],
	  "range": "0.7.0 - 1.0.0",
	  "nodes": ["node_modules/fbjs"],
	  "fixAvailable": {
		"name": "react",
		"version": "17.0.1",
		"isSemVerMajor": true
	  }
	},
	"isomorphic-fetch": {
	  "name": "isomorphic-fetch",
	  "severity": "low",
	  "via": ["node-fetch"],
	  "effects": ["fbjs"],
	  "range": "2.0.0 - 2.2.1",
	  "nodes": ["node_modules/isomorphic-fetch"],
	  "fixAvailable": true
	},
	"node-fetch": {
	  "name": "node-fetch",
	  "severity": "low",
	  "via": [
		{
		  "source": 1556,
		  "name": "node-fetch",
		  "dependency": "node-fetch",
		  "title": "Denial of Service",
		  "url": "https://npmjs.com/advisories/1556",
		  "severity": "low",
		  "range": "< 2.6.1 || >= 3.0.0-beta.1 < 3.0.0-beta.9"
		}
	  ],
	  "effects": ["isomorphic-fetch"],
	  "range": "<=2.6.0 || 3.0.0-beta.1 - 3.0.0-beta.8",
	  "nodes": ["node_modules/node-fetch"],
	  "fixAvailable": {
		"name": "react",
		"version": "17.0.1",
		"isSemVerMajor": true
	  }
	}
  },
  "metadata": {
	"vulnerabilities": {
	  "info": 0,
	  "low": 3,
	  "moderate": 0,
	  "high": 0,
	  "critical": 0,
	  "total": 3
	},
	"dependencies": {
	  "prod": 39,
	  "dev": 766,
	  "optional": 766,
	  "peer": 766,
	  "peerOptional": 0,
	  "total": 805
	}
  }
}
`

var quickAuditReport = &Report{
	PackagePath: "test",
	Version:     2,
	Vulnerabilities: map[string]Vulnerability{
		"fbjs": {
			Package:  "fbjs",
			Severity: "low",
			Via: []Advisory{
				{
					Transitive: true,
					Package:    "isomorphic-fetch",
					Dependency: "isomorphic-fetch",
				},
			},
			Effects: []string{},
			Range:   "0.7.0 - 1.0.0",
			Fix: Fix{
				Available: true,
				Package:   "react",
				Version:   "17.0.1",
				IsMajor:   true,
			},
		},
		"isomorphic-fetch": {
			Package:  "isomorphic-fetch",
			Severity: "low",
			Via: []Advisory{
				{
					Transitive: true,
					Package:    "node-fetch",
					Dependency: "node-fetch",
				},
			},
			Effects: []string{"fbjs"},
			Range:   "2.0.0 - 2.2.1",
			Fix: Fix{
				Available: true,
			},
		},
		"node-fetch": {
			Package:  "node-fetch",
			Severity: "low",
			Via: []Advisory{
				{
					Transitive: false,
					ID:         1556,
					Package:    "node-fetch",
					Dependency: "node-fetch",
					Title:      "Denial of Service",
					URL:        "https://npmjs.com/advisories/1556",
					Severity:   "low",
					Range:      "< 2.6.1 || >= 3.0.0-beta.1 < 3.0.0-beta.9",
				},
			},
			Effects: []string{"isomorphic-fetch"},
			Range:   "<=2.6.0 || 3.0.0-beta.1 - 3.0.0-beta.8",
			Fix: Fix{
				Available: true,
				Package:   "react",
				Version:   "17.0.1",
				IsMajor:   true,
			},
		},
	},
}

func TestNewReportValid(t *testing.T) {
	report, err := NewReport([]byte(quickAuditReportJSON))
	assert.NoError(t, err)
	report.SetPackagePath("test")
	assert.True(t, assert.ObjectsAreEqual(quickAuditReport, report))
}

var quickAuditIssues = []*v1.Issue{
	{
		Target:      "test:node-fetch",
		Type:        "Vulnerable Dependency",
		Title:       "Denial of Service",
		Severity:    v1.Severity_SEVERITY_LOW,
		Confidence:  v1.Confidence_CONFIDENCE_HIGH,
		Description: "Vulnerable versions: < 2.6.1 || >= 3.0.0-beta.1 < 3.0.0-beta.9\nRecommendation: Upgrade to version 2.6.1 or 3.0.0-beta.9\nOverview: Node Fetch did not honor the size option after following a redirect, which means that when a content size was over the limit, a FetchError would never get thrown and the process would end without failure.\n\nFor most people, this fix will have a little or no impact. However, if you are relying on node-fetch to gate files above a size, the impact could be significant, for example: If you don't double-check the size of the data after fetch() has completed, your JS thread could get tied up doing work on a large file (DoS) and/or cost you money in computing.\nReferences: \nNPM advisory URL: https://npmjs.com/advisories/1556\n",
	},
}

func TestAsIssuesValid(t *testing.T) {
	defer gock.Off()
	gock.New("https://npmjs.com").
		Get("/advisories/1556").
		MatchHeader("X-Spiferack", "1").
		Reply(200).
		AddHeader("Content-Type", "application/json").
		File("components/producers/typescript-npm-audit/types/npmquickaudit/npm_advisory_1556")

	issues := quickAuditReport.AsIssues()
	assert.True(t, assert.ObjectsAreEqual(quickAuditIssues, issues))

	assert.True(t, gock.IsDone())
}
