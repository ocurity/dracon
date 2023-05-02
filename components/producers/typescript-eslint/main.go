// Package main of the eslint producer transforms eslint json output to dracon issues
package main

import (
	"fmt"
	"log"

	v1 "github.com/ocurity/dracon/api/proto/v1"
	"github.com/ocurity/dracon/components/producers/typescript-eslint/types"

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
	if err := producers.ParseJSON(inFile, &results); err != nil {
		log.Fatal(err)
	}
	issues := parseIssues(results)
	if err := producers.WriteDraconOut(
		"eslint",
		issues,
	); err != nil {
		log.Fatal(err)
	}
}

func parseIssues(out []types.ESLintIssue) []*v1.Issue {
	issues := []*v1.Issue{}
	for _, r := range out {
		for _, msg := range r.Messages {
			sev := v1.Severity_SEVERITY_LOW
			if msg.Severity == 1 {
				sev = v1.Severity_SEVERITY_MEDIUM
			} else if msg.Severity == 2 {
				sev = v1.Severity_SEVERITY_HIGH
			}
			iss := &v1.Issue{
				Target:      fmt.Sprintf("%s:%v:%v", r.FilePath, msg.Line, msg.Column),
				Type:        msg.RuleID,
				Title:       msg.RuleID,
				Severity:    sev,
				Cvss:        0.0,
				Confidence:  v1.Confidence_CONFIDENCE_MEDIUM,
				Description: msg.Message,
			}
			issues = append(issues, iss)
		}
	}
	return issues
}
