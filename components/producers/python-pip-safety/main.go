// Package main of the pip-safety producer transforms safety json output to dracon issues
package main

import (
	"fmt"
	"log"

	v1 "github.com/ocurity/dracon/api/proto/v1"
	"github.com/ocurity/dracon/components/producers"
	"github.com/ocurity/dracon/components/producers/python-pip-safety/types"
)

func parseIssues(out []types.Vulnerability) []*v1.Issue {
	issues := []*v1.Issue{}
	for _, r := range out {
		issues = append(issues, &v1.Issue{
			Target:     fmt.Sprintf("%s:%s", r.PackageName, r.AnalyzedVersion),
			Type:       "Vulnerable Dependency",
			Title:      fmt.Sprintf("%s%s", r.PackageName, r.VulnerableSpec),
			Severity:   v1.Severity(v1.Severity_value[fmt.Sprintf("SEVERITY_%s", r.Severity)]),
			Cvss:       0.0,
			Confidence: v1.Confidence_CONFIDENCE_MEDIUM,
			Description: fmt.Sprintf("Advisory: %s\nFixed Versions: %v,Resources: %v, More Info: %s",
				r.Advisory,
				r.FixedVersions,
				r.Resources,
				r.MoreInfoURL),
			Cve: r.CVE,
		})
	}
	return issues
}

func main() {
	if err := producers.ParseFlags(); err != nil {
		log.Fatal(err)
	}

	inFile, err := producers.ReadInFile()
	if err != nil {
		log.Fatal(err)
	}

	issues := types.Out{}
	if err := producers.ParseJSON(inFile, &issues); err != nil {
		log.Fatal(err)
	}

	if err := producers.WriteDraconOut(
		"pip-safety",
		parseIssues(issues.Vulnerabilities),
	); err != nil {
		log.Fatal(err)
	}
}
