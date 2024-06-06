package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	v1 "github.com/ocurity/dracon/api/proto/v1"
	"github.com/ocurity/dracon/components/producers/semgrep/types"
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

	var results types.SemgrepResults
	if err := json.Unmarshal(inFile, &results); err != nil {
		log.Fatal(err)
	}

	issues, err := parseIssues(results)
	if err != nil {
		log.Fatal(err)
	}
	if err := producers.WriteDraconOut(
		"semgrep",
		issues,
	); err != nil {
		log.Fatal(err)
	}
}

func parseIssues(out types.SemgrepResults) ([]*v1.Issue, error) {
	issues := []*v1.Issue{}

	results := out.Results

	for _, r := range results {

		// Map the semgrep severity levels to dracon severity levels
		severityMap := map[string]v1.Severity{
			"INFO":    v1.Severity_SEVERITY_INFO,
			"WARNING": v1.Severity_SEVERITY_MEDIUM,
			"ERROR":   v1.Severity_SEVERITY_HIGH,
		}

		sev := severityMap[r.Extra.Severity]
		iss := &v1.Issue{
			Target:      fmt.Sprintf("%s:%v-%v", r.Path, r.Start.Line, r.End.Line),
			Type:        r.Extra.Message,
			Title:       r.CheckID,
			Severity:    sev,
			Cvss:        0.0,
			Confidence:  v1.Confidence_CONFIDENCE_MEDIUM,
			Description: fmt.Sprintf("%s\n extra lines: %s", r.Extra.Message, r.Extra.Lines),
			Cwe:         handleSemgrepCWE(r.Extra.Metadata.CWE),
		}
		cs, err := context.ExtractCode(iss)
		if err != nil {
			return nil, err
		}
		iss.ContextSegment = &cs
		issues = append(issues, iss)
	}
	return issues, nil
}

// Semgrep CWEs can be a string or an array of strings
func handleSemgrepCWE(cwe interface{}) []int32 {
	cweInts := []int32{}

	switch v := cwe.(type) {
	case []interface{}:
		for _, s := range v {
			cweAsInt := convertCWEStringToInt(s.(string))
			cweInts = append(cweInts, int32(cweAsInt))
		}
	case string:
		cweAsInt := convertCWEStringToInt(v)
		cweInts = append(cweInts, int32(cweAsInt))
	default:
		log.Fatalf("invalid cwe type: %T", cwe)
	}

	return cweInts
}

// Convert Semgrep CWE string to int.
// They always follow the schema `CWE-<number>:<description>`
func convertCWEStringToInt(cwe string) int32 {
	parts := strings.Split(cwe, "-")
	if len(parts) < 2 {
		log.Fatal("invalid cwe format yy")
	}
	subParts := strings.Split(parts[1], ":")
	if len(subParts) < 1 {
		log.Fatal("invalid cwe format zz")
	}
	cweAsInt, err := strconv.Atoi(subParts[0])
	if err != nil {
		log.Fatal(err)
	}
	return int32(cweAsInt)
}
