package main

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"

	v1 "github.com/ocurity/dracon/api/proto/v1"
	"github.com/ocurity/dracon/pkg/context"

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

	var results GoSecOut
	if err := json.Unmarshal(inFile, &results); err != nil {
		log.Fatal(err)
	}

	issues, err := parseIssues(&results)
	if err != nil {
		log.Fatal(err)
	}
	if err := producers.WriteDraconOut(
		"gosec",
		issues,
	); err != nil {
		log.Fatal(err)
	}
}

func parseIssues(out *GoSecOut) ([]*v1.Issue, error) {
	issues := []*v1.Issue{}
	for _, r := range out.Issues {
		iss := &v1.Issue{
			Target:      fmt.Sprintf("%s:%v", r.File, r.Line),
			Type:        r.RuleID,
			Title:       r.Details,
			Severity:    v1.Severity(v1.Severity_value[fmt.Sprintf("SEVERITY_%s", r.Severity)]),
			Cvss:        0.0,
			Confidence:  v1.Confidence(v1.Confidence_value[fmt.Sprintf("CONFIDENCE_%s", r.Confidence)]),
			Description: r.Code,
		}

		// Extract the code snippet, if possible
		code, err := context.ExtractCode(iss)
		if err != nil {
			slog.Warn("Failed to extract code snippet", "error", err)
			code = ""
		}
		iss.ContextSegment = &code

		issues = append(issues, iss)
	}
	return issues, nil
}

// GoSecOut represents the output of a GoSec run.
type GoSecOut struct {
	Issues []GoSecIssue `json:"Issues"`
	// Stats  GoSecStats   `json:"Stats"`
}

// GoSecIssue represents a GoSec Result.
type GoSecIssue struct {
	Severity   string `json:"severity"`
	Confidence string `json:"confidence"`
	RuleID     string `json:"rule_id"`
	Details    string `json:"details"`
	File       string `json:"file"`
	Code       string `json:"code"`
	Line       string `json:"line"`
	Column     string `json:"column"`
}
