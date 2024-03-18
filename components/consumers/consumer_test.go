package consumers

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	v1 "github.com/ocurity/dracon/api/proto/v1"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ocurity/dracon/pkg/putil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadToolResponse(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dracon-test")
	require.NoError(t, err)

	tmpFile, err := os.CreateTemp(tmpDir, "dracon-test-*.pb")
	require.NoError(t, err)

	defer require.NoError(t, os.Remove(tmpFile.Name()))

	issues := []*v1.Issue{
		{
			Target:      "/dracon/source/foobar",
			Title:       "/dracon/source/barfoo",
			Description: "/dracon/source/example.yaml",
		},
	}
	timestamp := time.Now().UTC().Format(time.RFC3339)
	scanID := "ab3d3290-cd9f-482c-97dc-ec48bdfcc4de"
	tags := map[string]string{
		"assetID":       "someID",
		"assetPriority": "priotity",
	}
	scanTags, err := json.Marshal(tags)
	assert.NoError(t, err)

	require.NoError(t, os.Setenv(EnvDraconStartTime, timestamp))
	require.NoError(t, os.Setenv(EnvDraconScanID, scanID))
	require.NoError(t, os.Setenv(EnvDraconScanTags, string(scanTags)))

	resultTempDir := tmpFile.Name()
	resultFile := "test-tool"
	assert.NoError(t, putil.WriteResults(resultFile, issues, resultTempDir, scanID, timestamp, tags))

	toolRes, err := putil.LoadToolResponse(resultTempDir)
	assert.NoError(t, err)

	assert.Equal(t, "test-tool", toolRes[0].GetToolName())
	assert.Equal(t, scanID, toolRes[0].GetScanInfo().GetScanUuid())
	assert.Equal(t, tags, toolRes[0].GetScanInfo().GetScanTags())
}

func TestFlatenLaunchToolResponse(t *testing.T) {

	issues := []*v1.Issue{}
	for i := 0; i < 4; i++ {
		issues = append(issues, &v1.Issue{
			Target:      "/dracon/source/foobar",
			Title:       "/dracon/source/barfoo",
			Description: "/dracon/source/example.yaml",
			Type:        fmt.Sprintf("%d", i),
			Cvss:        float64(i),
		})
	}
	scanID := "ab3d3290-cd9f-482c-97dc-ec48bdfcc4de"
	tags := map[string]string{
		"assetID":       "someID",
		"assetPriority": "priotity",
	}
	ts := timestamppb.New(time.Now().UTC())
	response := v1.LaunchToolResponse{
		ScanInfo: &v1.ScanInfo{
			ScanUuid:      scanID,
			ScanTags:      tags,
			ScanStartTime: ts,
		},
		ToolName: "unitTests",
		Issues:   issues,
	}
	expected := []map[string]string{
		{
			"CVE":                   "",
			"CVSS":                  "0.000000",
			"Confidence":            "CONFIDENCE_UNSPECIFIED",
			"ConfidenceText":        "CONFIDENCE_UNSPECIFIED",
			"CycloneDXSBOM":         "",
			"Description":           "/dracon/source/example.yaml",
			"ScanID":                "ab3d3290-cd9f-482c-97dc-ec48bdfcc4de",
			"ScanStartTime":         ts.AsTime().Format(time.RFC3339),
			"ScanTag:assetID":       "someID",
			"ScanTag:assetPriority": "priotity",
			"Severity":              "SEVERITY_UNSPECIFIED",
			"SeverityText":          "SEVERITY_UNSPECIFIED",
			"Source":                "",
			"Target":                "/dracon/source/foobar",
			"Title":                 "/dracon/source/barfoo",
			"ToolName":              "unitTests",
			"Type":                  "0"},
		{"CVE": "",
			"CVSS":                  "1.000000",
			"Confidence":            "CONFIDENCE_UNSPECIFIED",
			"ConfidenceText":        "CONFIDENCE_UNSPECIFIED",
			"CycloneDXSBOM":         "",
			"Description":           "/dracon/source/example.yaml",
			"ScanID":                "ab3d3290-cd9f-482c-97dc-ec48bdfcc4de",
			"ScanStartTime":         ts.AsTime().Format(time.RFC3339),
			"ScanTag:assetID":       "someID",
			"ScanTag:assetPriority": "priotity",
			"Severity":              "SEVERITY_UNSPECIFIED",
			"SeverityText":          "SEVERITY_UNSPECIFIED",
			"Source":                "",
			"Target":                "/dracon/source/foobar",
			"Title":                 "/dracon/source/barfoo",
			"ToolName":              "unitTests",
			"Type":                  "1"},
		{"CVE": "", "CVSS": "2.000000",
			"Confidence":            "CONFIDENCE_UNSPECIFIED",
			"ConfidenceText":        "CONFIDENCE_UNSPECIFIED",
			"CycloneDXSBOM":         "",
			"Description":           "/dracon/source/example.yaml",
			"ScanID":                "ab3d3290-cd9f-482c-97dc-ec48bdfcc4de",
			"ScanStartTime":         ts.AsTime().Format(time.RFC3339),
			"ScanTag:assetID":       "someID",
			"ScanTag:assetPriority": "priotity",
			"Severity":              "SEVERITY_UNSPECIFIED",
			"SeverityText":          "SEVERITY_UNSPECIFIED",
			"Source":                "",
			"Target":                "/dracon/source/foobar",
			"Title":                 "/dracon/source/barfoo",
			"ToolName":              "unitTests",
			"Type":                  "2"},
		{"CVE": "",
			"CVSS":                  "3.000000",
			"Confidence":            "CONFIDENCE_UNSPECIFIED",
			"ConfidenceText":        "CONFIDENCE_UNSPECIFIED",
			"CycloneDXSBOM":         "",
			"Description":           "/dracon/source/example.yaml",
			"ScanID":                "ab3d3290-cd9f-482c-97dc-ec48bdfcc4de",
			"ScanStartTime":         ts.AsTime().Format(time.RFC3339),
			"ScanTag:assetID":       "someID",
			"ScanTag:assetPriority": "priotity",
			"Severity":              "SEVERITY_UNSPECIFIED",
			"SeverityText":          "SEVERITY_UNSPECIFIED",
			"Source":                "",
			"Target":                "/dracon/source/foobar",
			"Title":                 "/dracon/source/barfoo",
			"ToolName":              "unitTests",
			"Type":                  "3"},
	}

	result := FlatenLaunchToolResponse(&response)
	assert.Equal(t, result, expected)
}
