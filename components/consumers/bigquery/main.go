// Package main of the bigquery consumer puts dracon issues into the target bigquery dataset, it will create teh dataset and the schema if one does not exist
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/ocurity/dracon/components/consumers"
	"github.com/ocurity/dracon/pkg/enumtransformers"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
)

var (
	projectID   string
	datasetName string
	gcpToken    string
)

func init() {
	// GOOGLE_APPLICATION_CREDENTIALS environment variable
	flag.StringVar(&gcpToken, "gcp-token", "", "The token used to access bigquery")
	flag.StringVar(&projectID, "project-id", "", "The bigquery project id to use.")
	flag.StringVar(&datasetName, "dataset", "", "The bigquery dataset to use.")
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
	token := &oauth2.Token{AccessToken: gcpToken}
	client, err := bigquery.NewClient(context.Background(), projectID, option.WithTokenSource(oauth2.StaticTokenSource(token)))
	if err != nil {
		return err
	}
	dataset := client.Dataset(datasetName)
	if _, err := dataset.Metadata(context.Background()); err != nil {
		log.Println("Dataset", dataset, "does not exist", "creating")
		if err = dataset.Create(context.Background(), &bigquery.DatasetMetadata{
			Name:        datasetName,
			Description: "a dataset to store findings from the Dracon ASOC framework",
			Location:    "EU",
		}); err != nil {
			return err
		}
	}
	table := dataset.Table("dracon")
	tmeta, err := table.Metadata(context.Background())
	if err != nil {
		log.Println("Table dracon does not exist creating")
		schema, err := bigquery.InferSchema(bqDraconIssue{})
		if err != nil {
			return err
		}
		if err = table.Create(context.Background(), &bigquery.TableMetadata{
			Name:        "dracon",
			Description: "a table to store dracon findings",
			Schema:      schema,
		}); err != nil {
			return err
		}
	} else if tmeta.Schema == nil {
		log.Println("Schema for table dracon does not exist creating")
		schema, err := bigquery.InferSchema(bqDraconIssue{})
		if err != nil {
			return err
		}
		_, err = table.Update(context.Background(), bigquery.TableMetadataToUpdate{
			Name:        "dracon",
			Description: "a table to store dracon findings",
			Schema:      schema,
		}, tmeta.ETag)
		if err != nil {
			return err
		}
	}
	inserter := table.Inserter()
	// Enumerate Dracon Issues to consume and create documents for each of them.
	if consumers.Raw {
		log.Println("Parsing Raw results")
		responses, err := consumers.LoadToolResponse()
		if err != nil {
			return fmt.Errorf("could not load raw results, file malformed: %w", err)
		}
		for _, res := range responses {
			for _, iss := range res.GetIssues() {
				if err = insert(ctx, inserter, bqDraconIssue{
					Confidence:    enumtransformers.ConfidenceToText(iss.GetConfidence()),
					Cve:           iss.GetCve(),
					Cvss:          iss.GetCvss(),
					Description:   iss.GetDescription(),
					IssueType:     iss.GetType(),
					ScanID:        res.GetScanInfo().GetScanUuid(),
					ScanStartTime: res.GetScanInfo().GetScanStartTime().AsTime(),
					Severity:      enumtransformers.SeverityToText(iss.GetSeverity()),
					Source:        iss.GetSource(),
					Target:        iss.GetTarget(),
					Title:         iss.GetTitle(),
					ToolName:      res.GetToolName(),
				}, iss.GetUuid()); err != nil {
					return err
				}
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
				var annotations []*bqDraconAnnotations
				for k, v := range iss.GetAnnotations() {
					annotations = append(annotations, &bqDraconAnnotations{Key: k, Value: v})
				}
				if err = insert(ctx, inserter, bqDraconIssue{
					Annotations:    annotations,
					Confidence:     enumtransformers.ConfidenceToText(iss.GetRawIssue().GetConfidence()),
					PreviousCounts: int(iss.GetCount()),
					Cve:            iss.GetRawIssue().GetCve(),
					Cvss:           iss.GetRawIssue().GetCvss(),
					Description:    iss.GetRawIssue().GetDescription(),
					FalsePositive:  iss.GetFalsePositive(),
					FirstFound:     iss.GetFirstSeen().AsTime(),
					IssueType:      iss.GetRawIssue().GetType(),
					LastFound:      iss.GetUpdatedAt().AsTime(),
					ScanID:         res.GetOriginalResults().GetScanInfo().GetScanUuid(),
					Severity:       enumtransformers.SeverityToText(iss.GetRawIssue().GetSeverity()),
					Source:         iss.GetRawIssue().GetSource(),
					Target:         iss.GetRawIssue().GetTarget(),
					Title:          iss.GetRawIssue().GetTitle(),
					ToolName:       res.GetOriginalResults().GetToolName(),
				}, iss.GetRawIssue().GetUuid()); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func insert(ctx context.Context, inserter *bigquery.Inserter, issue bqDraconIssue, issueUUID string) error {
	schema, err := bigquery.InferSchema(bqDraconIssue{})
	if err != nil {
		return err
	}
	data := &bigquery.StructSaver{
		Schema:   schema,
		InsertID: issueUUID,
		Struct:   issue,
	}
	return inserter.Put(ctx, data)
}

type bqDraconAnnotations struct {
	Key   string `bigquery:"key"`
	Value string `bigquery:"value"`
}

type bqDraconIssue struct {
	Annotations    []*bqDraconAnnotations `bigquery:"annotations"`
	Confidence     string                 `bigquery:"confidence"`
	Cve            string                 `bigquery:"cve"`
	Cvss           float64                `bigquery:"cvss"`
	Description    string                 `bigquery:"description"`
	FalsePositive  bool                   `bigquery:"falsePositive"`
	FirstFound     time.Time              `bigquery:"firstFound"`
	IssueType      string                 `bigquery:"issueType"`
	LastFound      time.Time              `bigquery:"lastFound"`
	PreviousCounts int                    `bigquery:"previousCounts"`
	ScanID         string                 `bigquery:"scanID"`
	ScanStartTime  time.Time              `bigquery:"scanStartTime"`
	Severity       string                 `bigquery:"severity"`
	Source         string                 `bigquery:"source"`
	Target         string                 `bigquery:"target"`
	Title          string                 `bigquery:"title"`
	ToolName       string                 `bigquery:"toolName"`
}
