package main

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
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
		cwe, err := handleSemgrepCWE(r.Extra.Metadata.CWE)
		if err != nil {
			slog.Warn("Couldn't parse CWE, skipping", err)
			cwe = []int32{}
		}

		iss := &v1.Issue{
			Target:      fmt.Sprintf("%s:%v-%v", r.Path, r.Start.Line, r.End.Line),
			Type:        r.Extra.Message,
			Title:       r.CheckID,
			Severity:    sev,
			Cvss:        0.0,
			Confidence:  v1.Confidence_CONFIDENCE_MEDIUM,
			Description: fmt.Sprintf("%s\n extra lines: %s", r.Extra.Message, r.Extra.Lines),
			Cwe:         cwe,
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
func handleSemgrepCWE(cwe interface{}) ([]int32, error) {
	cweInts := []int32{}

	switch v := cwe.(type) {
	case []interface{}:
		for _, s := range v {
			cweAsInt, err := convertCWEStringToInt(s.(string))
			if err != nil {
				return nil, err
			}
			cweInts = append(cweInts, int32(cweAsInt))
		}
	case string:
		cweAsInt, err := convertCWEStringToInt(v)
		if err != nil {
			return nil, err
		}
		cweInts = append(cweInts, int32(cweAsInt))
	default:
		return nil, fmt.Errorf("unexpected type for cwe: %T", v)
	}

	return cweInts, nil
}

// Convert Semgrep CWE string to int.
// They always follow the schema `CWE-<number>:<description>`
func convertCWEStringToInt(cwe string) (int32, error) {
	parts := strings.Split(cwe, "-")
	if len(parts) < 2 {
		return 0, fmt.Errorf("invalid cwe format, no '-' found")
	}
	subParts := strings.Split(parts[1], ":")
	if len(subParts) < 1 {
		return 0, fmt.Errorf("invalid cwe format, no ':' found")
	}
	cweAsInt, err := strconv.Atoi(subParts[0])
	if err != nil {
		return 0, err
	}
	return int32(cweAsInt), nil
}
