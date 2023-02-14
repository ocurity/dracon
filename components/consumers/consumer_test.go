package consumers

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"testing"
	"time"

	v1 "github.com/ocurity/dracon/api/proto/v1"

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
	os.Setenv(EnvDraconStartTime, timestamp)
	os.Setenv(EnvDraconScanID, scanID)
	err = putil.WriteResults("test-tool", issues, tmpFile.Name(), scanID, timestamp)
	assert.Nil(t, err)

	log.Println(tmpDir)
	inResults = path.Dir(tmpDir)

	toolRes, err := LoadToolResponse()
	assert.Nil(t, err)
	log.Println(toolRes)

	assert.Equal(t, "test-tool", toolRes[0].GetToolName())
	assert.Equal(t, scanID, toolRes[0].GetScanInfo().GetScanUuid())
}
