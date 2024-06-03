package consumers

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	draconapiv1 "github.com/ocurity/dracon/api/proto/v1"
	"github.com/ocurity/dracon/components"
	"github.com/ocurity/dracon/pkg/putil"
)

func TestLoadToolResponse(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dracon-test")
	require.NoError(t, err)

	tmpFile, err := os.CreateTemp(tmpDir, "dracon-test-*.pb")
	require.NoError(t, err)

	defer require.NoError(t, os.Remove(tmpFile.Name()))

	issues := []*draconapiv1.Issue{
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

	require.NoError(t, os.Setenv(components.EnvDraconStartTime, timestamp))
	require.NoError(t, os.Setenv(components.EnvDraconScanID, scanID))
	require.NoError(t, os.Setenv(components.EnvDraconScanTags, string(scanTags)))

	resultTempDir := tmpFile.Name()
	resultFile := "test-tool"
	assert.NoError(t, putil.WriteResults(resultFile, issues, resultTempDir, scanID, timestamp, tags))

	toolRes, err := putil.LoadToolResponse(resultTempDir)
	assert.NoError(t, err)

	assert.Equal(t, "test-tool", toolRes[0].GetToolName())
	assert.Equal(t, scanID, toolRes[0].GetScanInfo().GetScanUuid())
	assert.Equal(t, tags, toolRes[0].GetScanInfo().GetScanTags())
}
