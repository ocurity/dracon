package main

import (
	"fmt"
	v1 "github.com/ocurity/dracon/api/proto/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"os"
	"strings"
	"testing"
)

func countResults(res []*v1.EnrichedIssue) (int, int) {
	r := 0
	f := 0
	for _, finding := range res {
		if strings.Contains(fmt.Sprintf("%#v\n", finding.Annotations), "\"reachable\":\"false\"") {
			f++
		}
		if strings.Contains(fmt.Sprintf("%#v\n", finding.Annotations), "\"reachable\":\"true\"") {
			r++
		}
	}
	return r, f
}

func readPb(t *testing.T, err error, outFile string) []*v1.EnrichedIssue {
	pbBytes, err := os.ReadFile(outFile)
	require.NoError(t, err)

	res := v1.EnrichedLaunchToolResponse{}
	require.NoError(t, proto.Unmarshal(pbBytes, &res))
	return res.GetIssues()
}

func TestReachableEnricher(t *testing.T) {
	//prepare
	outDir, err := os.MkdirTemp("/tmp", "")
	require.NoError(t, err)

	dir := "testdata"
	readPath = dir
	writePath = outDir
	sliceFile = dir + "/sampleReachables.json"

	run()
	assert.FileExists(t, outDir+"/reachability.reachability.enriched.pb", "file was not created")
	assert.FileExists(t, outDir+"/bandit.reachability.enriched.pb", "file was not created")
	assert.FileExists(t, outDir+"/pip-safety.reachability.enriched.pb", "file was not created")

	res := readPb(t, err, outDir+"/reachability.reachability.enriched.pb")

	r, f := countResults(res)
	assert.Equal(t, 0, r)
	assert.Equal(t, 1, f)

	res = readPb(t, err, outDir+"/bandit.reachability.enriched.pb")

	r, f = countResults(res)
	assert.Equal(t, 2, r)
	assert.Equal(t, 0, f)

	res = readPb(t, err, outDir+"/pip-safety.reachability.enriched.pb")

	r, f = countResults(res)
	assert.Equal(t, 14, r)
	assert.Equal(t, 9, f)

}
