package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	v1 "github.com/ocurity/dracon/api/proto/v1"
	"github.com/ocurity/dracon/components/producers"
	"github.com/ocurity/dracon/components/producers/terraform-tfsec/types"
	"github.com/ocurity/dracon/pkg/sarif"
)

// Sarif flag to indicate the producer is being fed sarif input.
var Sarif bool

func main() {
	flag.BoolVar(&Sarif, "sarifOut", false, "Output is in sarif format}")

	if err := producers.ParseFlags(); err != nil {
		log.Fatal(err)
	}

	inFile, err := producers.ReadInFile()
	if err != nil {
		log.Fatal(err)
	}

	if Sarif {
		var sarifResults []*sarif.DraconIssueCollection
		var draconResults []*v1.Issue
		sarifResults, err := sarif.ToDracon(string(inFile))
		if err != nil {
			log.Fatal(err)
		}
		for _, result := range sarifResults {
			if result.ToolName != "defsec" {
				log.Printf("Toolname from Sarif results is not 'defsec' it is %s instead\n", result.ToolName)
			}
			draconResults = append(draconResults, result.Issues...)
		}
		if err := producers.WriteDraconOut("tfsec", draconResults); err != nil {
			log.Fatal(err)
		}
	} else {
		var results types.TfSecOut
		if err := producers.ParseJSON(inFile, &results); err != nil {
			log.Fatal(err)
		}
		if err := producers.WriteDraconOut("tfsec", parseOut(results)); err != nil {
			log.Fatal(err)
		}

	}
}

func parseOut(results types.TfSecOut) []*v1.Issue {
	issues := []*v1.Issue{}
	for _, res := range results.Results {
		description, _ := json.Marshal(res)
		issues = append(issues, &v1.Issue{
			Target: fmt.Sprintf("%s:%d-%d",
				res.Location.Filename,
				res.Location.StartLine,
				res.Location.EndLine),
			Type:        res.LongID,
			Title:       res.RuleDescription,
			Severity:    TfSecSeverityToDracon(res.Severity),
			Confidence:  v1.Confidence_CONFIDENCE_MEDIUM,
			Description: string(description),
		})
	}
	return issues
}

// TfSecSeverityToDracon maps tfsec Severity Strings to dracon struct.
func TfSecSeverityToDracon(severity string) v1.Severity {
	switch severity {
	case "LOW":
		return v1.Severity_SEVERITY_LOW
	case "MEDIUM":
		return v1.Severity_SEVERITY_MEDIUM
	case "HIGH":
		return v1.Severity_SEVERITY_HIGH
	case "CRITICAL":
		return v1.Severity_SEVERITY_CRITICAL
	default:
		return v1.Severity_SEVERITY_INFO
	}
}
