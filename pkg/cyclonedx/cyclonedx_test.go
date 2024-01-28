package cyclonedx

import (
	"encoding/json"
	"os"
	"testing"

	v1 "github.com/ocurity/dracon/api/proto/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToDraconLibrary(t *testing.T) {
	rawLibraryBOM, err := os.ReadFile("./testdata/libraryBOM.json")
	require.NoError(t, err)

	issues, err := ToDracon(rawLibraryBOM, "json")
	require.NoError(t, err)

	libraryBOM := string(rawLibraryBOM)

	expectedIssues := []*v1.Issue{
		{
			Target:        "pkg:npm/juice-shop@11.1.2",
			Type:          "SBOM",
			Title:         "SBOM for pkg:npm/juice-shop@11.1.2",
			Severity:      v1.Severity_SEVERITY_INFO,
			CycloneDXSBOM: &libraryBOM,
		},
	}
	assert.Equal(t, expectedIssues[0].Target, issues[0].Target)
	assert.Equal(t, expectedIssues[0].Type, issues[0].Type)
	assert.Equal(t, expectedIssues[0].Title, issues[0].Title)
	assert.Equal(t, expectedIssues[0].Severity, issues[0].Severity)
	var sbom1, sbom2 map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(*expectedIssues[0].CycloneDXSBOM), &sbom1))
	require.NoError(t, json.Unmarshal([]byte(*issues[0].CycloneDXSBOM), &sbom2))
	assert.Equal(t, sbom1, sbom2)
}

func TestToDraconSaaSInfra(t *testing.T) {
	rawSaaSBOM, err := os.ReadFile("./testdata/saasBOM.json")
	require.NoError(t, err)

	issues, err := ToDracon(rawSaaSBOM, "json")
	require.NoError(t, err)

	saasBOM := string(rawSaaSBOM)
	expectedIssues := []*v1.Issue{
		{
			Target:        "acme-application",
			Type:          "SBOM",
			Title:         "SBOM for acme-application",
			Severity:      v1.Severity_SEVERITY_INFO,
			CycloneDXSBOM: &saasBOM,
		},
	}
	assert.Equal(t, expectedIssues[0].Target, issues[0].Target)
	assert.Equal(t, expectedIssues[0].Type, issues[0].Type)
	assert.Equal(t, expectedIssues[0].Title, issues[0].Title)
	assert.Equal(t, expectedIssues[0].Severity, issues[0].Severity)

	var sbom1, sbom2 map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(*expectedIssues[0].CycloneDXSBOM), &sbom1))
	require.NoError(t, json.Unmarshal([]byte(*issues[0].CycloneDXSBOM), &sbom2))
	assert.Equal(t, sbom1, sbom2)
}
