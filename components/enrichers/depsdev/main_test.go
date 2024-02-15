package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	cdx "github.com/CycloneDX/cyclonedx-go"
	"github.com/google/uuid"
	v1 "github.com/ocurity/dracon/api/proto/v1"
	"github.com/ocurity/dracon/pkg/cyclonedx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	license = "Foo License v0"
)

func genSampleIssue(t *testing.T) []*v1.Issue {
	id := uuid.New()

	bom, err := os.ReadFile("./testdata/sampleSBOM.json")
	require.NoError(t, err)

	sampleSBOM := string(bom)
	rI := &v1.Issue{
		Target:        "some/target",
		Type:          "some type",
		Title:         "/some/target sbom",
		Severity:      v1.Severity_SEVERITY_INFO,
		Cvss:          0,
		Confidence:    v1.Confidence_CONFIDENCE_INFO,
		Description:   "foo bar",
		Cve:           "",
		Uuid:          id.String(),
		CycloneDXSBOM: &sampleSBOM,
	}
	return []*v1.Issue{rI}
}

func prepareIssue(t *testing.T) string {
	// prepare
	dir, err := os.MkdirTemp("/tmp", "")
	require.NoError(t, err)

	rawIssues := genSampleIssue(t)
	id := uuid.New()
	scanUUUID := id.String()
	startTime, _ := time.Parse(time.RFC3339, time.Now().UTC().Format(time.RFC3339))
	orig := v1.LaunchToolResponse{
		Issues:   rawIssues,
		ToolName: "depsdevSAT",
		ScanInfo: &v1.ScanInfo{
			ScanUuid:      scanUUUID,
			ScanStartTime: timestamppb.New(startTime),
		},
	}
	// write sample raw issues in mktemp
	out, _ := proto.Marshal(&orig)
	require.NoError(t, os.WriteFile(dir+"/depsdevSAT.tagged.pb", out, 0o600))

	readPath = dir
	writePath = dir
	return dir
}

// TODO make this be the common setup method
// todo add test for deps dev and scorecard stuff
func setup(t *testing.T) (string, *httptest.Server) {
	dir := prepareIssue(t)

	// setup server
	response := Response{
		Version: Version{
			Licenses: []string{license},
			Projects: []Project{
				{
					ScorecardV2: ScorecardV2{
						Date:  "irrelevant",
						Score: 5.5,
						Check: []Check{
							{
								Name:   "foo",
								Score:  2,
								Reason: "bar",
							},
						},
					},
				},
			},
		},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.String(), "/_/s/go/p/")
		require.NoError(t, json.NewEncoder(w).Encode(response))
	}))
	depsdevBaseURL = srv.URL
	return dir, srv
}

func TestParseIssuesDepsDevScoreCardInfoWritten(t *testing.T) {
	dir, srv := setup(t)
	defer srv.Close()
	scoreCardInfo = "true"
	// run enricher
	run()
	assert.FileExists(t, dir+"/depsdevSAT.depsdev.enriched.pb", "file was not created")

	// load *enriched.pb
	pbBytes, err := os.ReadFile(dir + "/depsdevSAT.depsdev.enriched.pb")
	require.NoError(t, err, "could not read enriched file")

	res := v1.EnrichedLaunchToolResponse{}
	require.NoError(t, proto.Unmarshal(pbBytes, &res))

	expectedProperties := []cdx.Property{
		{Name: "aquasecurity:trivy:PkgType", Value: "gomod"},
		{Name: "ScorecardScore", Value: "5.500000"},
		{Name: "ScorecardInfo", Value: "{\n\t\"date\": \"irrelevant\",\n\t\"repo\": {},\n\t\"scorecard\": {},\n\t\"check\": [\n\t\t{\n\t\t\t\"name\": \"foo\",\n\t\t\t\"documentation\": {},\n\t\t\t\"score\": 2,\n\t\t\t\"reason\": \"bar\"\n\t\t}\n\t],\n\t\"score\": 5.5\n}"},
		{Name: "aquasecurity:trivy:PkgType", Value: "gomod"},
		{Name: "ScorecardScore", Value: "5.500000"},
		{Name: "ScorecardInfo", Value: "{\n\t\"date\": \"irrelevant\",\n\t\"repo\": {},\n\t\"scorecard\": {},\n\t\"check\": [\n\t\t{\n\t\t\t\"name\": \"foo\",\n\t\t\t\"documentation\": {},\n\t\t\t\"score\": 2,\n\t\t\t\"reason\": \"bar\"\n\t\t}\n\t],\n\t\"score\": 5.5\n}"},
		{Name: "aquasecurity:trivy:PkgType", Value: "gomod"},
		{Name: "ScorecardScore", Value: "5.500000"},
		{Name: "ScorecardInfo", Value: "{\n\t\"date\": \"irrelevant\",\n\t\"repo\": {},\n\t\"scorecard\": {},\n\t\"check\": [\n\t\t{\n\t\t\t\"name\": \"foo\",\n\t\t\t\"documentation\": {},\n\t\t\t\"score\": 2,\n\t\t\t\"reason\": \"bar\"\n\t\t}\n\t],\n\t\"score\": 5.5\n}"},
	}
	//  ensure every component has a license attached to it
	for _, finding := range res.Issues {
		bom, err := cyclonedx.FromDracon(finding.RawIssue)
		require.NoError(t, err, "Could not read enriched cyclone dx info")

		properties := []cdx.Property{}

		for _, component := range *bom.Components {
			properties = append(properties, *component.Properties...)
		}
		assert.Equal(t, properties, expectedProperties)
	}
}

func TestParseIssuesDepsDevExternalReferenceLinksWritten(t *testing.T) {
	dir, srv := setup(t)
	defer srv.Close()

	// run enricher
	run()
	assert.FileExists(t, dir+"/depsdevSAT.depsdev.enriched.pb", "file was not created")

	// load *enriched.pb
	pbBytes, err := os.ReadFile(dir + "/depsdevSAT.depsdev.enriched.pb")
	require.NoError(t, err, "could not read enriched file")

	res := v1.EnrichedLaunchToolResponse{}
	require.NoError(t, proto.Unmarshal(pbBytes, &res))

	expectedExternalReferences := []cdx.ExternalReference{
		{
			URL:  fmt.Sprintf("%s/go/p/cloud.google.com%%2Fgo%%2Fcompute/v/v1.14.0", srv.URL),
			Type: "other",
		}, {
			URL:  fmt.Sprintf("%s/go/p/cloud.google.com%%2Fgo%%2Fcompute%%2Fmetadata/v/v0.2.3", srv.URL),
			Type: "other",
		}, {
			URL:  fmt.Sprintf("%s/go/p/github.com%%2FAzure%%2Fazure-pipeline-go/v/v0.2.3", srv.URL),
			Type: "other",
		},
	}
	//  ensure every component has a license attached to it
	for _, finding := range res.Issues {
		bom, err := cyclonedx.FromDracon(finding.RawIssue)
		require.NoError(t, err, "Could not read enriched cyclone dx info")

		externalReferences := []cdx.ExternalReference{}

		for _, component := range *bom.Components {
			externalReferences = append(externalReferences, *component.ExternalReferences...)
		}
		assert.Equal(t, externalReferences, expectedExternalReferences)
	}
}

func TestParseIssuesLicensesWritten(t *testing.T) {
	dir, srv := setup(t)
	defer srv.Close()

	licensesInEvidence = "false"

	// run enricher
	run()
	assert.FileExists(t, dir+"/depsdevSAT.depsdev.enriched.pb", "file was not created")

	// load *enriched.pb
	pbBytes, err := os.ReadFile(dir + "/depsdevSAT.depsdev.enriched.pb")
	require.NoError(t, err, "could not read enriched file")
	res := v1.EnrichedLaunchToolResponse{}
	require.NoError(t, proto.Unmarshal(pbBytes, &res))

	//  ensure every component has a license attached to it
	for _, finding := range res.Issues {
		bom, err := cyclonedx.FromDracon(finding.RawIssue)
		require.NoError(t, err, "Could not read enriched cyclone dx info")
		found := false
		for _, component := range *bom.Components {
			for _, lic := range *component.Licenses {
				found = true
				assert.Equal(t, lic.License.Name, license)
			}
		}
		assert.True(t, found)
	}
}

func TestParseIssuesLicensesWrittenACcurateLicenses(t *testing.T) {
	dir, srv := setup(t)
	defer srv.Close()
	licensesInEvidence = "true"

	run()

	assert.FileExists(t, dir+"/depsdevSAT.depsdev.enriched.pb", "file was not created")

	// load *enriched.pb
	pbBytes, err := os.ReadFile(dir + "/depsdevSAT.depsdev.enriched.pb")
	require.NoError(t, err, "could not read enriched file")

	res := v1.EnrichedLaunchToolResponse{}
	require.NoError(t, proto.Unmarshal(pbBytes, &res))

	//  ensure every component has a license attached to it
	for _, finding := range res.Issues {
		bom, err := cyclonedx.FromDracon(finding.RawIssue)
		require.NoError(t, err, "Could not read enriched cyclone dx info")
		found := false
		for _, component := range *bom.Components {
			for _, lic := range *component.Evidence.Licenses {
				found = true
				assert.Equal(t, lic.License.Name, license)
			}
		}
		assert.True(t, found)
	}
}
