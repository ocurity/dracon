package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/bigquery"
	v1 "github.com/ocurity/dracon/api/proto/v1"
	"github.com/ocurity/dracon/components/consumers"
)

var (
	dbURL     string
	projectID string
	dataset   string
)

func init() {
	// GOOGLE_APPLICATION_CREDENTIALS environment variable
	flag.StringVar(&projectID, "project-id", "", "The bigquery project id to use.")
	flag.StringVar(&dataset, "dataset", "", "The bigquery dataset to use.")
}

func main() {
	if err := consumers.ParseFlags(); err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	if err := run(ctx); err != nil {
		log.Fatalf("could not run: %s", err)
	}
}

func run(ctx context.Context) error {
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return err
	}
	inserter := client.Dataset(dataset).Table("dracon").Inserter()

	// Enumerate Dracon Issues to consume and create documents for each of them.
	if consumers.Raw {
		log.Println("Parsing Raw results")
		responses, err := consumers.LoadToolResponse()
		if err != nil {
			return fmt.Errorf("could not load raw results, file malformed: %w", err)
		}
		for _, res := range responses {
			for _, iss := range res.GetIssues() {
				return insert(ctx, inserter, iss)
			}
		}
	} else {
		log.Print("Parsing Enriched results")
		responses, err := consumers.LoadEnrichedToolResponse()
		if err != nil {
			return fmt.Errorf("could not load enriched results, file malformed: %w", err)
		}
		for _, res := range responses {
			for _, iss := range res.GetIssues() {
				return insert(ctx, inserter, iss)
			}
		}
	}

	return nil
}

func insert(ctx context.Context, inserter *bigquery.Inserter, issue interface{}, scanID, toolName string, scanStartTime time.Time) error {
	schema, err := bigquery.InferSchema(bqDraconIssue{})
	var data interface{}
	switch e := issue.(type) {
	case *v1.Issue:
		iss, _ := issue.(*v1.Issue)
		data = &bigquery.StructSaver{
			Schema:   schema,
			InsertID: e.Spec.ID,
			Struct: bqDraconIssue{
				confidence:    iss.GetConfidence(),
				cve:           iss.GetCve(),
				cvss:          iss.GetCvss(),
				description:   iss.GetDescription(),
				issueType:     iss.GetType(),
				scanID:        scanID,
				scanStartTime: scanStartTime,
				severity:      iss.GetSeverity(),
				source:        iss.GetSource(),
				target:        iss.GetTarget(),
				title:         iss.GetTitle(),
				toolName:      toolName,
			},
		}
	case *v1.EnrichedIssue:
		iss, _ := issue.(*v1.EnrichedIssue)
		data = &bigquery.StructSaver{
			Schema:   schema,
			InsertID: e.Spec.ID,
			Struct: bqDraconIssue{
				annotations:    iss.GetAnnotations(),
				confidence:     iss.GetRawIssue().GetConfidence(),
				previousCounts: iss.GetCount(),
				cve:            iss.GetRawIssue().GetCve(),
				cvss:           iss.GetRawIssue().GetCvss(),
				description:    iss.GetRawIssue().GetDescription(),
				falsePositive:  iss.GetFalsePositive(),
				firstFound:     iss.GetFirstFound(),
				issueType:      iss.GetRawIssue().GetType(),
				lastFound:      iss.GetUpdatedAt(),
				scanID:         scanID,
				scanStartTime:  scanStartTime,
				severity:       iss.GetRawIssue().GetSeverity(),
				source:         iss.GetRawIssue().GetSource(),
				target:         iss.GetRawIssue().GetTarget(),
				title:          iss.GetRawIssue().GetTitle(),
				toolName:       toolName,
			},
		}
	}
	return inserter.Put(ctx, data)
}

type bqDraconIssue struct {
	annotations    map[string]string `bigquery:"annotations"`
	confidence     string            `bigquery:"confidence"`
	cve            string            `bigquery:"cve"`
	cvss           float32           `bigquery:"cvss"`
	description    string            `bigquery:"description"`
	falsePositive  bool              `bigquery:"falsePositive"`
	firstFound     time.Time         `bigquery:"firstFound"`
	issueType      string            `bigquery:"issueType"`
	lastFound      time.Time         `bigquery:"lastFound"`
	previousCounts int               `bigquery:"previousCounts"`
	scanID         string            `bigquery:"scanID"`
	scanStartTime  time.Time         `bigquery:"scanStartTime"`
	severity       string            `bigquery:"severity"`
	source         string            `bigquery:"source"`
	target         string            `bigquery:"target"`
	title          string            `bigquery:"title"`
	toolName       string            `bigquery:"toolName"`
}
