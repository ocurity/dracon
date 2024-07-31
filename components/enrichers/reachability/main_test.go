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

func TestReachableEnricher(t *testing.T) {
	// prepare
	//dir, err := os.MkdirTemp("/tmp", "")
	//require.NoError(t, err)

	dir := "testdata"
	readPath = dir
	writePath = dir
	sliceFile = dir + "/sampleReachables.json"

	run()

	pbBytes, err := os.ReadFile(dir + "/reachability.reachability.enriched.pb")
	require.NoError(t, err)

	res := v1.EnrichedLaunchToolResponse{}
	require.NoError(t, proto.Unmarshal(pbBytes, &res))

	for _, finding := range res.Issues {
		assert.Contains(t, fmt.Sprintf("%#v\n", finding.Annotations), "\"reachable\":\"false\"")
	}

	pbBytes, err = os.ReadFile(dir + "/bandit.reachability.enriched.pb")
	require.NoError(t, err)

	res = v1.EnrichedLaunchToolResponse{}
	require.NoError(t, proto.Unmarshal(pbBytes, &res))

	for _, finding := range res.Issues {
		assert.Contains(t, fmt.Sprintf("%#v\n", finding.Annotations), "\"reachable\":\"true\"")
	}

	pbBytes, err = os.ReadFile(dir + "/pip-safety.reachability.enriched.pb")
	require.NoError(t, err)

	res = v1.EnrichedLaunchToolResponse{}
	require.NoError(t, proto.Unmarshal(pbBytes, &res))

	r := 0
	f := 0
	for _, finding := range res.Issues {
		if strings.Contains(fmt.Sprintf("%#v\n", finding.Annotations), "\"reachable\":\"false\"") {
			f++
		}
		if strings.Contains(fmt.Sprintf("%#v\n", finding.Annotations), "\"reachable\":\"true\"") {
			r++
		}
	}
	assert.Equal(t, 14, r)
	assert.Equal(t, 9, f)

}
