package main

import (
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"strings"

	dtrack "github.com/DependencyTrack/client-go"
	"github.com/google/uuid"
	v1 "github.com/ocurity/dracon/api/proto/v1"
	"github.com/ocurity/dracon/components/consumers"
	cyclonedx "github.com/ocurity/dracon/pkg/cyclonedx"
)

var (
	authURL         string
	apiKey          string
	projectName     string
	projectVersion  string
	projectUUID     string
	client          *dtrack.Client
	ownerAnnotation string
)

func main() {
	flag.StringVar(&apiKey, "apiKey", "", "dependency track api key")
	flag.StringVar(&authURL, "url", "", "dependency track base url")
	flag.StringVar(&projectName, "projectName", "", "dependency track project name")
	flag.StringVar(&projectUUID, "projectUUID", "", "dependency track project name")
	flag.StringVar(&projectVersion, "projectVersion", "", "dependency track project version")
	flag.StringVar(
		&ownerAnnotation,
		"ownerAnnotation",
		"",
		"if this consumer is in running after any enricher that adds ownership annotations, then provide the annotation-key for this enricher so it can tag owners as tags",
	)

	flag.Parse()

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
		var bomIssue *v1.EnrichedIssue
		for _, issue := range res.GetIssues() {
			if issue.GetRawIssue().GetCycloneDXSBOM() != "" && bomIssue == nil {
				bomIssue = issue
			} else if bomIssue != nil && bomIssue.GetRawIssue().GetCycloneDXSBOM() != "" {
				log.Printf("Tool response for tool %s is malformed, we expected a single issue with an SBOM as part of the tool, got something else instead",
					res.GetOriginalResults().GetToolName())
				continue
			}
		}
		cdxbom, err := cyclonedx.FromDracon(bomIssue.GetRawIssue())
		if err != nil {
			return tokens, err
		}
		token, err := uploadBOM(bomIssue.GetRawIssue().GetCycloneDXSBOM(), cdxbom.Metadata.Component.Version)
		if err != nil {
			log.Fatal("could not upload bom to dependency track, err:", err)
		}
		log.Println("upload token is", token)
		tokens = append(tokens, token)
		if ownerAnnotation != "" {
			log.Println("tagging owners")
			owners := []string{}
			for key, value := range bomIssue.Annotations {
				if strings.Contains(key, ownerAnnotation) {
					owners = append(owners, value)
				}
			}
			if err := addOwnersTags(owners); err != nil {
				log.Println("could not tag owners, err:", err)
			}
		}
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
				log.Printf("Tool response for tool %s is malformed, we expected a single issue with an SBOM as part of the tool, got multiple issues with sboms instead",
					res.GetToolName())
				continue
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

func addOwnersTags(owners []string) error {
	// addOwnersTags expects a map of <ownerAnnotation>-<number>:<username> tagging owners
	// it then adds to the projectUUID the owners in the following tag format: Owner:<username>
	uuid := uuid.MustParse(projectUUID)
	project, err := client.Project.Get(context.Background(), uuid)
	if err != nil {
		log.Println("could not add project tags error getting project by uuid, err:", err)
		return err
	}
	for _, owner := range owners {
		found := false
		for _, t := range project.Tags {
			if t.Name == fmt.Sprintf("%s:%s", ownerAnnotation, owner) {
				found = true
				break
			}
		}
		if !found {
			project.Tags = append(project.Tags, dtrack.Tag{Name: fmt.Sprintf("%s:%s", ownerAnnotation, owner)})
		}
	}
	_, err = client.Project.Update(context.Background(), project)
	return err
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
