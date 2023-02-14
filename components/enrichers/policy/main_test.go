package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	v1 "github.com/ocurity/dracon/api/proto/v1"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	// numFindings         = 20.
	expectedNumFindings = 5 // we have 5 distinct uuids
)

func genSampleIssues() []*v1.Issue {
	sev := []v1.Severity{
		v1.Severity_SEVERITY_CRITICAL,
		v1.Severity_SEVERITY_HIGH,
		v1.Severity_SEVERITY_MEDIUM,
		v1.Severity_SEVERITY_LOW,
		v1.Severity_SEVERITY_INFO,
	}
	raw := make([]*v1.Issue, 0)
	for i := 0; i < expectedNumFindings; i++ {
		id := uuid.New()
		rI := &v1.Issue{
			Target:      fmt.Sprintf("%d some/target", i),
			Type:        fmt.Sprintf("%d some type", i),
			Title:       fmt.Sprintf("%d /some/target is vulnerable", i),
			Severity:    sev[i%len(sev)],
			Cvss:        float64(i),
			Confidence:  v1.Confidence_CONFIDENCE_MEDIUM,
			Description: fmt.Sprintf("%d foo bar", i),
			Cve:         "CVE-2017-11770",
			Uuid:        id.String(),
		}

		raw = append(raw, rI)

	}
	return raw
}

func MockServer(t *testing.T) {
}

func TestParseIssues(t *testing.T) {
	// prepare
	dir, err := ioutil.TempDir("/tmp", "")
	if err != nil {
		log.Fatal(err)
	}
	rawIssues := genSampleIssues()

	id := uuid.New()
	scanUUUID := id.String()
	startTime, _ := time.Parse(time.RFC3339, time.Now().UTC().Format(time.RFC3339))
	orig := v1.LaunchToolResponse{
		Issues:   rawIssues,
		ToolName: "policySat",
		ScanInfo: &v1.ScanInfo{
			ScanUuid:      scanUUUID,
			ScanStartTime: timestamppb.New(startTime),
		},
	}
	// write sample raw issues in mktemp
	out, _ := proto.Marshal(&orig)
	ioutil.WriteFile(dir+"/policySat.tagged.pb", out, 0o600)

	readPath = dir
	writePath = dir
	policy = "cGFja2FnZSBleGFtcGxlLnBvbGljeVNhdAoKZGVmYXVsdCBhbGxvdyA6PSBmYWxzZQoKYWxsb3cgPXRydWUgewogICAgcHJpbnQoaW5wdXQpCiAgICBjaGVja19zZXZlcml0eQp9CgpjaGVja19zZXZlcml0eSB7CiAgICBpbnB1dC5zZXZlcml0eSA9PSAiU0VWRVJJVFlfTE9XIgp9CgpjaGVja19zZXZlcml0eSB7CiAgICBpbnB1dC5zZXZlcml0eSA9PSAiU0VWRVJJVFlfSElHSCIKfQpjaGVja19zZXZlcml0eSB7CiAgICBpbnB1dC5zZXZlcml0eSA9PSAiU0VWRVJJVFlfTUVESVVNIgp9Cg=="

	// setup server
	expected := "{}"
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type bod map[string]interface{}
		var body bod

		if strings.Contains(r.URL.String(), "/v1/policies") {
			fmt.Fprintf(w, "{}")
		}

		if strings.Contains(r.URL.String(), "/v1/data") {
			json.NewDecoder(r.Body).Decode(&body)
			if strings.Contains(fmt.Sprintf("%#v\n ", body), "SEVERITY_CRITICAL") {
				fmt.Fprintf(w, "{\"result\":{\"allow\":false}}")
			} else {
				fmt.Fprintf(w, "{\"result\":{\"allow\":true}}")
			}
		}

		fmt.Fprint(w, expected)
	}))
	defer svr.Close()
	regoServer = svr.URL

	run()

	assert.FileExists(t, dir+"/policySat.policy.enriched.pb", "file was not created")

	// load *enriched.pb
	pbBytes, err := ioutil.ReadFile(dir + "/policySat.policy.enriched.pb")
	res := v1.EnrichedLaunchToolResponse{}
	proto.Unmarshal(pbBytes, &res)

	//  ensure every finding has a policy decision and that for every finding the handler was called once
	for _, finding := range res.Issues {
		if finding.RawIssue.Severity != v1.Severity_SEVERITY_CRITICAL {
			assert.True(t, strings.Contains(fmt.Sprintf("%#v\n", finding.Annotations), "\"Policy Pass: example/policySat\":\"true\""))
		}
	}
}
