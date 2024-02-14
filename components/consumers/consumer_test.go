package consumers

import (
	"encoding/json"
	v1 "github.com/ocurity/dracon/api/proto/v1"
	"io/ioutil"
	"log"
	"os"
	"path"
	"testing"
	"time"

	"github.com/ocurity/dracon/pkg/putil"
	"github.com/stretchr/testify/assert"
)

func TestLoadToolResponse(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "dracon-test")
	assert.Nil(t, err)
	tmpFile, err := ioutil.TempFile(tmpDir, "dracon-test-*.pb")
	assert.Nil(t, err)
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
	tags := map[string]string{
		"assetID":       "someID",
		"assetPriority": "priotity",
	}
	scanTags, err := json.Marshal(tags)
	assert.Nil(t, err)
	os.Setenv(EnvDraconStartTime, timestamp)
	os.Setenv(EnvDraconScanID, scanID)
	os.Setenv(EnvDraconScanTags, string(scanTags))
	err = putil.WriteResults("test-tool", issues, tmpFile.Name(), scanID, timestamp, tags)
	assert.Nil(t, err)

	log.Println(tmpDir)
	inResults = path.Dir(tmpDir)

	toolRes, err := LoadToolResponse()
	assert.Nil(t, err)
	log.Println(toolRes)

	assert.Equal(t, "test-tool", toolRes[0].GetToolName())
	assert.Equal(t, scanID, toolRes[0].GetScanInfo().GetScanUuid())
	assert.Equal(t, tags, toolRes[0].GetScanInfo().GetScanTags())
}
