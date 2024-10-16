package main

import (
	"encoding/json"
	"log"
	"log/slog"

	v1 "github.com/ocurity/dracon/api/proto/v1"

	"github.com/ocurity/dracon/components/producers"
)

func main() {
	if err := producers.ParseFlags(); err != nil {
		log.Fatal(err)
	}

	inFile, err := producers.ReadInFile()
	if err != nil {
		log.Fatal(err)
	}

	var results ModelScanOut
	if err := json.Unmarshal(inFile, &results); err != nil {
		log.Fatal(err)
	}

	issues, err := parseIssues(&results)
	if err != nil {
		log.Fatal(err)
	}
	if err := producers.WriteDraconOut(
		"modelscan",
		issues,
	); err != nil {
		log.Fatal(err)
	}
}

func parseIssues(out *ModelScanOut) ([]*v1.Issue, error) {
	issues := []*v1.Issue{}
	slog.Info("found Critical issues", slog.Int("numCrit", len(out.IssuesBySeverity.Critical)))
	for _, crit := range out.IssuesBySeverity.Critical {
		issues = append(issues,
			&v1.Issue{
				Target:      crit.Source,
				Type:        crit.Scanner,
				Description: crit.Description,
				Title:       crit.Description,
				Severity:    v1.Severity_SEVERITY_UNSPECIFIED,
				Confidence:  v1.Confidence_CONFIDENCE_UNSPECIFIED,
			})
	}
	slog.Info("found High issues", slog.Int("numHigh", len(out.IssuesBySeverity.High)))
	for _, crit := range out.IssuesBySeverity.High {
		issues = append(issues,
			&v1.Issue{
				Target:      crit.Source,
				Type:        crit.Scanner,
				Description: crit.Description,
				Title:       crit.Description,
				Severity:    v1.Severity_SEVERITY_UNSPECIFIED,
				Confidence:  v1.Confidence_CONFIDENCE_UNSPECIFIED,
			})
	}
	slog.Info("found Medium issues", slog.Int("numMedium", len(out.IssuesBySeverity.Medium)))
	for _, crit := range out.IssuesBySeverity.Medium {
		issues = append(issues,
			&v1.Issue{
				Target:      crit.Source,
				Type:        crit.Scanner,
				Description: crit.Description,
				Title:       crit.Description,
				Severity:    v1.Severity_SEVERITY_UNSPECIFIED,
				Confidence:  v1.Confidence_CONFIDENCE_UNSPECIFIED,
			})
	}
	slog.Info("found Low issues", slog.Int("numLow", len(out.IssuesBySeverity.Low)))
	for _, crit := range out.IssuesBySeverity.Low {
		issues = append(issues,
			&v1.Issue{
				Target:      crit.Source,
				Type:        crit.Scanner,
				Description: crit.Description,
				Title:       crit.Description,
				Severity:    v1.Severity_SEVERITY_UNSPECIFIED,
				Confidence:  v1.Confidence_CONFIDENCE_UNSPECIFIED,
			})
	}
	return issues, nil
}

type ModelScanOut struct {
	ModelscanVersion string `json:"modelscan_version"`
	Timestamp        string `json:"timestamp"`
	InputPath        string `json:"input_path"`
	TotalIssues      int    `json:"total_issues"`
	Summary          struct {
		TotalIssuesBySeverity struct {
			Low      int `json:"LOW"`
			Medium   int `json:"MEDIUM"`
			High     int `json:"HIGH"`
			Critical int `json:"CRITICAL"`
		} `json:"total_issues_by_severity"`
	} `json:"summary"`
	IssuesBySeverity struct {
		Critical []struct {
			Description string `json:"description"`
			Operator    string `json:"operator"`
			Module      string `json:"module"`
			Source      string `json:"source"`
			Scanner     string `json:"scanner"`
		} `json:"CRITICAL"`
		High []struct {
			Description string `json:"description"`
			Operator    string `json:"operator"`
			Module      string `json:"module"`
			Source      string `json:"source"`
			Scanner     string `json:"scanner"`
		} `json:"HIGH"`
		Medium []struct {
			Description string `json:"description"`
			Operator    string `json:"operator"`
			Module      string `json:"module"`
			Source      string `json:"source"`
			Scanner     string `json:"scanner"`
		} `json:"MEDIUM"`
		Low []struct {
			Description string `json:"description"`
			Operator    string `json:"operator"`
			Module      string `json:"module"`
			Source      string `json:"source"`
			Scanner     string `json:"scanner"`
		} `json:"LOW"`
	} `json:"issues_by_severity"`
	Errors  []any `json:"errors"`
	Scanned struct {
		TotalScanned int      `json:"total_scanned"`
		ScannedFiles []string `json:"scanned_files"`
	} `json:"scanned"`
}
