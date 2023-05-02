// Package main of the dependency track consumer puts dracon issues into the target dependency track
package main

import (
	"context"
	"encoding/base64"
	"flag"
	"log"

	dtrack "github.com/DependencyTrack/client-go"
	"github.com/google/uuid"
	v1 "github.com/ocurity/dracon/api/proto/v1"
	"github.com/ocurity/dracon/components/consumers"
	cyclonedx "github.com/ocurity/dracon/pkg/cyclonedx"
)

var (
	authURL        string
	apiKey         string
	projectName    string
	projectVersion string
	projectUUID    string
	client         *dtrack.Client
)

func init() {
	flag.StringVar(&apiKey, "apiKey", "", "dependency track api key")
	flag.StringVar(&authURL, "url", "", "dependency track base url")
	flag.StringVar(&projectName, "projectName", "", "dependency track project name")
	flag.StringVar(&projectUUID, "projectUUID", "", "dependency track project name")
	flag.StringVar(&projectVersion, "projectVersion", "", "dependency track project version")
	flag.Parse()
}

func main() {
	if err := consumers.ParseFlags(); err != nil {
		log.Fatal(err)
	}
	if projectUUID == "" {
		log.Fatal("project uuid is mandatory for dependency track")
	}
	c, err := dtrack.NewClient(authURL, dtrack.WithAPIKey(apiKey))
	if err != nil {
		log.Panicf("could not instantiate client err: %#v\n", err)
	}
	client = c
	if consumers.Raw {
		responses, err := consumers.LoadToolResponse()
		if err != nil {
			log.Fatal("could not load raw results, file malformed: ", err)
		}
		if _, err := uploadBOMsFromRaw(responses); err != nil {
			log.Fatalf("could not upload boms from raw: %s", err)
		}

	} else {
		responses, err := consumers.LoadEnrichedToolResponse()
		if err != nil {
			log.Fatal("could not load enriched results, file malformed: ", err)
		}
		if _, err := uploadBOMSFromEnriched(responses); err != nil {
			log.Fatalf("could not upload boms from enriched: %s", err)
		}
	}
}

func uploadBOMSFromEnriched(responses []*v1.EnrichedLaunchToolResponse) ([]string, error) {
	var tokens []string
	for _, res := range responses {
		var bomIssue *v1.Issue
		for _, issue := range res.GetIssues() {
			if issue.GetRawIssue().GetCycloneDXSBOM() != "" && bomIssue == nil {
				bomIssue = issue.GetRawIssue()
			} else if bomIssue != nil && *bomIssue.CycloneDXSBOM != "" {
				log.Fatalf("Tool response for tool %s is malformed, we expected a single issue with an SBOM as part of the tool, got something else instead",
					res.GetOriginalResults().GetToolName())
			}
		}
		cdxbom, err := cyclonedx.FromDracon(bomIssue)
		if err != nil {
			return tokens, err
		}
		token, err := uploadBOM(bomIssue.GetCycloneDXSBOM(), cdxbom.Metadata.Component.Version)
		if err != nil {
			log.Fatal("could not upload bom to dependency track, err:", err)
		}
		log.Println("upload token is", token)
		tokens = append(tokens, token)
	}
	return tokens, nil
}

func uploadBOMsFromRaw(responses []*v1.LaunchToolResponse) ([]string, error) {
	var tokens []string
	for _, res := range responses {
		var bomIssue *v1.Issue
		for _, issue := range res.GetIssues() {
			if *issue.CycloneDXSBOM != "" && bomIssue == nil {
				bomIssue = issue
			} else if bomIssue != nil && *bomIssue.CycloneDXSBOM != "" {
				log.Fatalf("Tool response for tool %s is malformed, we expected a single issue with an SBOM as part of the tool, got multiple issues with sboms instead",
					res.GetToolName())
			}
		}
		cdxbom, err := cyclonedx.FromDracon(bomIssue)
		if err != nil {
			return tokens, err
		}
		token, err := uploadBOM(*bomIssue.CycloneDXSBOM, cdxbom.Metadata.Component.Version)
		if err != nil {
			log.Fatal("could not upload bod to dependency track, err:", err)
		}
		log.Println("upload token is", token)
		tokens = append(tokens, token)
	}
	return tokens, nil
}

func uploadBOM(bom string, projectVersion string) (string, error) {
	if projectVersion == "" {
		projectVersion = "Unknown"
	}
	uuid := uuid.MustParse(projectUUID)
	token, err := client.BOM.Upload(context.TODO(), dtrack.BOMUploadRequest{
		ProjectName:    projectName,
		ProjectVersion: projectVersion,
		ProjectUUID:    &uuid,
		BOM:            base64.StdEncoding.EncodeToString([]byte(bom)),
	})
	return string(token), err
}
