package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"
	"time"

	"github.com/google/uuid"
	v1 "github.com/ocurity/dracon/api/proto/v1"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func genSampleIssues() []*v1.Issue {
	issues := []*v1.Issue{{}}
	for i := 0; i < 5; i++ {
		newIssue := &v1.Issue{
			Target:      fmt.Sprintf("%d some/target", i),
			Type:        fmt.Sprintf("%d some type"),
			Title:       fmt.Sprintf("%d /some/target is vulnerable"),
			Severity:    v1.Severity_SEVERITY_HIGH,
			Cvss:        float64(i),
			Confidence:  v1.Confidence_CONFIDENCE_MEDIUM,
			Description: fmt.Sprintf("%d foo bar"),
			Cve:         "CVE-2017-11770",
		}
		issues = append(issues, newIssue)
	}
	return issues
}

func TestParseIssues(t *testing.T) {
	// prepare
	dir, err := ioutil.TempDir("/tmp", "")
	if err != nil {
		log.Fatal(err)
	}
	issues := genSampleIssues()
	id := uuid.New()
	scanUUUID := id.String()
	startTime, _ := time.Parse(time.RFC3339, time.Now().UTC().Format(time.RFC3339))
	resp := v1.LaunchToolResponse{
		Issues:   issues,
		ToolName: "taggerSat",
		ScanInfo: &v1.ScanInfo{
			ScanUuid:      scanUUUID,
			ScanStartTime: timestamppb.New(startTime),
		},
	}

	// write sample raw issues in mktemp
	out, _ := proto.Marshal(&resp)
	ioutil.WriteFile(dir+"/taggerSat.pb", out, 0o600)

	readPath = dir
	writePath = dir

	run()

	assert.FileExists(t, dir+"/taggerSat.tagged.pb", "Tagged file was not created")

	// load *tagged.pb
	pbBytes, err := ioutil.ReadFile(dir + "/taggerSat.tagged.pb")
	res := v1.LaunchToolResponse{}
	proto.Unmarshal(pbBytes, &res)

	// ensure every issue has a uuid populated
	for _, issue := range res.Issues {
		assert.NotEqual(t, issue.Uuid, "")
		assert.NotNil(t, issue.Uuid)
	}
}
