package main

import (
	"encoding/json"
	"log"
	"log/slog"

	v1 "github.com/ocurity/dracon/api/proto/v1"
	"github.com/ocurity/dracon/components/producers/typescript-eslint/types"
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

	var results []types.ESLintIssue
	if err := json.Unmarshal(inFile, &results); err != nil {
		log.Fatal(err)
	}
	issues, err := parseIssues(results)
	if err != nil {
		log.Fatal(err)
	}
	if err := producers.WriteDraconOut(
		"eslint",
		issues,
	); err != nil {
		log.Fatal(err)
	}
}

func parseIssues(out []types.ESLintIssue) ([]*v1.Issue, error) {
	issues := []*v1.Issue{}
	for _, r := range out {
		for _, msg := range r.Messages {
			// Convert the severity
			sev := v1.Severity_SEVERITY_LOW
			if msg.Severity == 1 {
				sev = v1.Severity_SEVERITY_MEDIUM
			} else if msg.Severity == 2 {
				sev = v1.Severity_SEVERITY_HIGH
			}

			// Ensure we always have a valid endLine
			endLine := msg.EndLine
			if endLine == 0 {
				endLine = msg.Line
			}

			iss := &v1.Issue{
				Target:      producers.GetFileTarget(r.FilePath, msg.Line, endLine),
				Type:        msg.RuleID,
				Title:       msg.RuleID,
				Severity:    sev,
				Cvss:        0.0,
				Confidence:  v1.Confidence_CONFIDENCE_MEDIUM,
				Description: msg.Message,
			}

			// Extract the code snippet, if possible
			code, err := context.ExtractCodeFromFileTarget(iss.Target)
			if err != nil {
				slog.Warn("Failed to extract code snippet", "error", err)
				code = ""
			}
			iss.ContextSegment = &code

			issues = append(issues, iss)
		}
	}
	return issues, nil
}
