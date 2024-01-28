package consumers

import (
	"os"
	"testing"
	"time"

	v1 "github.com/ocurity/dracon/api/proto/v1"

	"github.com/ocurity/dracon/pkg/putil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadToolResponse(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dracon-test")
	require.NoError(t, err)

	tmpFile, err := os.CreateTemp(tmpDir, "dracon-test-*.pb")
	require.NoError(t, err)

	defer os.Remove(tmpFile.Name())
	issues := []*v1.Issue{
		{
			Target:      "/dracon/source/foobar",
			Title:       "/dracon/source/barfoo",
			Description: "/dracon/source/example.yaml",
		},
	}
	timestamp := time.Now().UTC().Format(time.RFC3339)
	scanID := "ab3d3290-cd9f-482c-97dc-ec48bdfcc4de"
	os.Setenv(EnvDraconStartTime, timestamp)
	os.Setenv(EnvDraconScanID, scanID)

	require.NoError(t, putil.WriteResults("test-tool", issues, tmpFile.Name(), scanID, timestamp))
	inResults = tmpDir

	toolRes, err := LoadToolResponse()
	require.NoError(t, err)

	assert.Equal(t, "test-tool", toolRes[0].GetToolName(), toolRes)
	assert.Equal(t, scanID, toolRes[0].GetScanInfo().GetScanUuid())
	assert.Equal(t, tags, toolRes[0].GetScanInfo().GetScanTags())
}
