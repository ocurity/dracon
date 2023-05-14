package main

import (
	"fmt"
	"log"

	v1 "github.com/ocurity/dracon/api/proto/v1"
	"github.com/ocurity/dracon/components/producers"
	"github.com/ocurity/dracon/components/producers/zaproxy/types"
)

func main() {
	if err := producers.ParseFlags(); err != nil {
		log.Fatal(err)
	}

	inFile, err := producers.ReadInFile()
	if err != nil {
		log.Fatal(err)
	}

	var results types.ZapOut
	if err := producers.ParseJSON(inFile, &results); err != nil {
		log.Fatal(err)
	}

	if err := producers.WriteDraconOut("zap", parseOut(&results)); err != nil {
		log.Fatal(err)
	}
}

func parseOut(results *types.ZapOut) []*v1.Issue {
	issues := []*v1.Issue{}
	for _, res := range results.Site {
		target := res.Name
		for _, alert := range res.Alerts {
			parsed := parseIssue(alert, target)
			issues = append(issues, parsed...)
		}
	}
	return issues
}

// zap doesn't provide cvss so assigned as 0.0.
func parseIssue(r *types.ZapAlerts, target string) []*v1.Issue {
	cvss := 0.0
	var results []*v1.Issue
	if r.Instances != nil {
		for _, instance := range r.Instances {
			results = append(results, &v1.Issue{
				Target:      instance.URI,
				Type:        r.CweID,
				Title:       r.Name,
				Severity:    riskcodeToSeverity(r.RiskCode),
				Confidence:  zapconfidenceToConfidence(r.Confidence),
				Cvss:        cvss,
				Description: fmt.Sprintf("Description: %s\nSolution: %s\nReference: %s\nAttack: %s", r.Description, r.Solution, r.Reference, instance.Attack),
			})
		}
	}
	return results
}

// riskcode values are 0-INFO,1-LOW,2-MEDIUM,3-HIGH only available from ZAP. It is determined by the ZAP contributors.
func riskcodeToSeverity(riskcode string) v1.Severity {
	switch riskcode {
	case "0":

		return v1.Severity_SEVERITY_INFO

	case "1":

		return v1.Severity_SEVERITY_LOW

	case "2":

		return v1.Severity_SEVERITY_MEDIUM

	case "3":

		return v1.Severity_SEVERITY_HIGH

	default:

		return v1.Severity_SEVERITY_CRITICAL
	}
}

// Confidence values are 0-INFO,1-LOW,2-MEDIUM,3-HIGH only available from ZAP. It is determined by the ZAP contributors.
func zapconfidenceToConfidence(confidence string) v1.Confidence {
	switch confidence {
	case "0":

		return v1.Confidence_CONFIDENCE_INFO

	case "1":

		return v1.Confidence_CONFIDENCE_LOW

	case "2":

		return v1.Confidence_CONFIDENCE_MEDIUM

	case "3":

		return v1.Confidence_CONFIDENCE_HIGH

	default:

		return v1.Confidence_CONFIDENCE_CRITICAL
	}
}
