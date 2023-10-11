// Package main of the trivy component transforms the json,cyclonedx or sarif output of Trivy to dracon
package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	v1 "github.com/ocurity/dracon/api/proto/v1"
	"github.com/ocurity/dracon/components/producers"
	"github.com/ocurity/dracon/components/producers/docker-trivy/types"
	"github.com/ocurity/dracon/pkg/cyclonedx"
	"github.com/ocurity/dracon/pkg/sarif"
)

// Combined flag to indicate the producer is being fed  aggregated input from multiple images.
var Combined bool

var (
	// Format is what was passed while running trivy -f.
	Format           string
	supportedFormats = []string{"json", "sarif", "cyclonedx"}
)

func main() {
	flag.BoolVar(&Combined, "combinedout", false, "Output is the combined output of Trivy against multiple images, expects {<img-name>:[<regular trivy output>],<other-img>:[<trivy out for 'other-img'>]}")
	flag.StringVar(&Format, "format", "json", fmt.Sprintf("trivy output format, by default json, supported formats are %v", supportedFormats))

	if err := producers.ParseFlags(); err != nil {
		log.Fatal(err)
	}
	var results []*v1.Issue
	inFile, err := producers.ReadInFile()
	if err != nil {
		log.Fatal(err)
	}
	switch Format {
	case "json":
		results, err = handleJSON(inFile)
	case "sarif":
		results, err = handleSarif(inFile)
	case "cyclonedx":
		results, err = handleCycloneDX(inFile)
	case "spdx":
		log.Fatal("SPDX is not supported, please use cyclonedx instead")
	case "spdx-json":
		log.Fatal("SPDX is not supported, please use cyclonedx instead")
	default:
		log.Fatal(fmt.Errorf("Format %s is not supported, supported formats are %v", Format, supportedFormats))
	}
	if err != nil {
		log.Fatal(err)
	}
	if err := producers.WriteDraconOut(
		"trivy", results,
	); err != nil {
		log.Fatal(err)
	}
}

func handleJSON(inFile []byte) ([]*v1.Issue, error) {
	if Combined {
		var results types.CombinedOut
		if err := producers.ParseJSON(inFile, &results); err != nil {
			return []*v1.Issue{}, err
		}
		return parseCombinedOut(results), nil
	}
	var results types.TrivyOut
	if err := producers.ParseJSON(inFile, &results); err != nil {
		return []*v1.Issue{}, err
	}
	return parseSingleOut(results), nil
}

func handleSarif(inFile []byte) ([]*v1.Issue, error) {
	var sarifResults []*sarif.DraconIssueCollection
	var draconResults []*v1.Issue
	sarifResults, err := sarif.ToDracon(string(inFile))
	if err != nil {
		return draconResults, err
	}
	for _, result := range sarifResults {
		if strings.ToLower(result.ToolName) != "trivy" {
			log.Printf("Toolname from Sarif results is not 'trivy' it is %s instead\n", result.ToolName)
		}
		draconResults = append(draconResults, result.Issues...)
	}
	return draconResults, nil
}

func handleCycloneDX(inFile []byte) ([]*v1.Issue, error) {
	return cyclonedx.ToDracon(inFile, "json")
}

func parseCombinedOut(results types.CombinedOut) []*v1.Issue {
	issues := []*v1.Issue{}
	for img, output := range results {
		log.Printf("Parsing Combined Output for %s\n", img)
		for _, res := range output.Results {
			target := res.Target
			for _, vuln := range res.Vulnerable {
				issues = append(issues, parseResult(vuln, target))
			}
		}
	}
	return issues
}

func parseSingleOut(results types.TrivyOut) []*v1.Issue {
	issues := []*v1.Issue{}
	for _, res := range results.Results {
		target := res.Target

		for _, vuln := range res.Vulnerable {
			issues = append(issues, parseResult(vuln, target))
		}
	}
	return issues
}

// TrivySeverityToDracon maps Trivy Severity Strings to dracon struct.
func TrivySeverityToDracon(severity string) v1.Severity {
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

func parseResult(r *types.TrivyVulnerability, target string) *v1.Issue {
	return &v1.Issue{
		Target:     target,
		Type:       "Container image vulnerability",
		Title:      fmt.Sprintf("[%s][%s] %s", target, r.CVE, r.Title),
		Severity:   TrivySeverityToDracon(r.Severity),
		Confidence: v1.Confidence_CONFIDENCE_UNSPECIFIED,
		Cvss:       r.CVSS.Nvd.V3Score,
		Description: fmt.Sprintf("CVSS Score: %v\nCvssVector: %s\nCve: %s\nCwe: %s\nReference: %s\nOriginal Description:%s\n",
			r.CVSS.Nvd.V3Score, r.CVSS.Nvd.V3Vector, r.CVE, strings.Join(r.CweIDs, ","), r.PrimaryURL, r.Description),
	}
}
