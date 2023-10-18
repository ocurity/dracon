package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	cdx "github.com/CycloneDX/cyclonedx-go"
	"github.com/google/uuid"
	v1 "github.com/ocurity/dracon/api/proto/v1"
	"github.com/ocurity/dracon/pkg/cyclonedx"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	license = "Foo License v0"
)

func genSampleIssue() []*v1.Issue {
	id := uuid.New()
	bom := sampleSbom
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
		CycloneDXSBOM: &bom,
	}
	return []*v1.Issue{rI}
}

func prepareIssue() string {
	// prepare
	dir, err := ioutil.TempDir("/tmp", "")
	if err != nil {
		log.Fatal(err)
	}
	rawIssues := genSampleIssue()

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
	ioutil.WriteFile(dir+"/depsdevSAT.tagged.pb", out, 0o600)

	readPath = dir
	writePath = dir
	return dir
}

func MockServer(t *testing.T) {
}

// TODO make this be the common setup method
// todo add test for deps dev and scorecard stuff
func setup(t *testing.T) (string, *httptest.Server) {
	dir := prepareIssue()

	// setup server
	response := Response{
		Version: Version{
			Licenses: []string{license},
			Projects: []Project{
				{
					ScorecardV2: ScorecardV2{
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
		json.NewEncoder(w).Encode(response)
	}))
	depsdevBaseURL = srv.URL
	return dir, srv
}


func TestParseIssuesDepsDevScoreCardInfoWritten(t *testing.T) {
	dir, srv := setup(t)
	defer srv.Close()

	// run enricher
	run()
	assert.FileExists(t, dir+"/depsdevSAT.depsdev.enriched.pb", "file was not created")

	// load *enriched.pb
	pbBytes, err := ioutil.ReadFile(dir + "/depsdevSAT.depsdev.enriched.pb")
	assert.NoError(t, err, "could not read enriched file")
	res := v1.EnrichedLaunchToolResponse{}
	proto.Unmarshal(pbBytes, &res)
	expectedExternalReferences := []cdx.ExternalReference{
		{
			URL:  "http://127.0.0.1:46679//go/p/cloud.google.com%2Fgo%2Fcompute/v/v1.14.0",
			Type: "other",
		}, {
			URL:  "http://127.0.0.1:46679//go/p/cloud.google.com%2Fgo%2Fcompute%2Fmetadata/v/v0.2.3",
			Type: "other",
		}, {
			URL:  "http://127.0.0.1:46679//go/p/github.com%2FAzure%2Fazure-pipeline-go/v/v0.2.3",
			Type: "other",
		},
	}
	//  ensure every component has a license attached to it
	for _, finding := range res.Issues {
		bom, err := cyclonedx.FromDracon(finding.RawIssue)
		assert.NoError(t, err, "Could not read enriched cyclone dx info")

		externalReferences := []cdx.ExternalReference{}

		for _, component := range *bom.Components {
			externalReferences = append(externalReferences, *component.ExternalReferences...)
		}
		assert.Equal(t, externalReferences, expectedExternalReferences)
	}
}


func TestParseIssuesDepsDevExternalReferenceLinksWritten(t *testing.T) {
	dir, srv := setup(t)
	defer srv.Close()

	// run enricher
	run()
	assert.FileExists(t, dir+"/depsdevSAT.depsdev.enriched.pb", "file was not created")

	// load *enriched.pb
	pbBytes, err := ioutil.ReadFile(dir + "/depsdevSAT.depsdev.enriched.pb")
	assert.NoError(t, err, "could not read enriched file")
	res := v1.EnrichedLaunchToolResponse{}
	proto.Unmarshal(pbBytes, &res)
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
		assert.NoError(t, err, "Could not read enriched cyclone dx info")

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
	pbBytes, err := ioutil.ReadFile(dir + "/depsdevSAT.depsdev.enriched.pb")
	assert.NoError(t, err, "could not read enriched file")
	res := v1.EnrichedLaunchToolResponse{}
	proto.Unmarshal(pbBytes, &res)

	//  ensure every component has a license attached to it
	for _, finding := range res.Issues {
		bom, err := cyclonedx.FromDracon(finding.RawIssue)
		assert.NoError(t, err, "Could not read enriched cyclone dx info")
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
	pbBytes, err := ioutil.ReadFile(dir + "/depsdevSAT.depsdev.enriched.pb")
	assert.NoError(t, err, "could not read enriched file")
	res := v1.EnrichedLaunchToolResponse{}
	proto.Unmarshal(pbBytes, &res)

	//  ensure every component has a license attached to it
	for _, finding := range res.Issues {
		bom, err := cyclonedx.FromDracon(finding.RawIssue)
		assert.NoError(t, err, "Could not read enriched cyclone dx info")
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

const sampleSbom = `{
	"bomFormat": "CycloneDX",
	"specVersion": "1.4",
	"serialNumber": "urn:uuid:2a73a682-6094-45e0-9a87-544abd01fd8a",
	"version": 1,
	"metadata": {
	  "timestamp": "2023-03-01T08:53:41+00:00",
	  "tools": [
		{
		  "vendor": "aquasecurity",
		  "name": "trivy",
		  "version": "0.36.1"
		}
	  ],
	  "component": {
		"bom-ref": "7d68c283-756d-4c19-a255-cde007a0437a",
		"type": "application",
		"name": "/code",
		"properties": [
		  {
			"name": "aquasecurity:trivy:SchemaVersion",
			"value": "2"
		  }
		]
	  }
	},
	"components": [
	  {
		"bom-ref": "pkg:golang/cloud.google.com/go/compute@1.14.0",
		"type": "library",
		"name": "cloud.google.com/go/compute",
		"version": "1.14.0",
		"purl": "pkg:golang/cloud.google.com/go/compute@1.14.0",
		"properties": [
		  {
			"name": "aquasecurity:trivy:PkgType",
			"value": "gomod"
		  }
		]
	  },
	  {
		"bom-ref": "pkg:golang/cloud.google.com/go/compute/metadata@0.2.3",
		"type": "library",
		"name": "cloud.google.com/go/compute/metadata",
		"version": "0.2.3",
		"purl": "pkg:golang/cloud.google.com/go/compute/metadata@0.2.3",
		"properties": [
		  {
			"name": "aquasecurity:trivy:PkgType",
			"value": "gomod"
		  }
		]
	  },
	  {
		"bom-ref": "pkg:golang/github.com/azure/azure-pipeline-go@0.2.3",
		"type": "library",
		"name": "github.com/Azure/azure-pipeline-go",
		"version": "0.2.3",
		"purl": "pkg:golang/github.com/azure/azure-pipeline-go@0.2.3",
		"properties": [
		  {
			"name": "aquasecurity:trivy:PkgType",
			"value": "gomod"
		  }
		]
	  }
	],
	"dependencies": [],
	"vulnerabilities": []
  }
  `
