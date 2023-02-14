package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	v1 "github.com/ocurity/dracon/api/proto/v1"

	"github.com/ocurity/dracon/components/producers"
)

// DependencyVulnerability represents the Vulnerability block of Dependency check scan json output format.
type DependencyVulnerability struct {
	target      string
	cvss3       float64
	cwes        []interface{}
	notes       string
	name        string
	severity    string
	cvss2       float64
	description string
	cve         string
}

// UnmarshalJSON returns a list of Dependency Vulns from dependency check json.
func UnmarshalJSON(jsonBytes []byte) []DependencyVulnerability {
	var result []DependencyVulnerability
	var v map[string]interface{}
	if !json.Valid(jsonBytes) {
		log.Fatal("Inputfile not valid JSON")
	}
	if err := json.Unmarshal(jsonBytes, &v); err != nil {
		log.Fatal(err)
	}
	dependencies := v["dependencies"].([]interface{})
	for _, dependency := range dependencies {
		depmap := dependency.(map[string]interface{})
		if vulns, ok := depmap["vulnerabilities"]; ok {
			target := depmap["filePath"].(string)
			for _, vuln := range vulns.([]interface{}) {
				vv := vuln.(map[string]interface{})
				cvss3 := 0.0
				cvss2 := 0.0
				cve := ""
				if vv["source"].(string) == "NVD" {
					cve = vv["name"].(string)
				}
				if vv["cvssv3"] != nil {
					v3 := vv["cvssv3"].(map[string]interface{})
					cvss3 = v3["baseScore"].(float64)
				}
				if vv["cvssv2"] != nil {
					v2 := vv["cvssv2"].(map[string]interface{})
					cvss2 = v2["score"].(float64)
				}
				result = append(result, DependencyVulnerability{
					target:      target,
					cvss3:       cvss3,
					cwes:        vv["cwes"].([]interface{}),
					notes:       vv["notes"].(string),
					name:        vv["name"].(string),
					severity:    vv["severity"].(string),
					cvss2:       cvss2,
					description: vv["description"].(string),
					cve:         cve,
				})
			}
		}
	}
	return result
}

func parseIssues(out []DependencyVulnerability) []*v1.Issue {
	issues := []*v1.Issue{}
	for _, r := range out {
		cvss := r.cvss2
		if r.cvss3 != 0.0 {
			cvss = r.cvss3
		}
		issues = append(issues, &v1.Issue{
			Target:      r.target,
			Type:        "Vulnerable Dependency",
			Title:       r.target,
			Severity:    v1.Severity(v1.Severity_value[fmt.Sprintf("SEVERITY_%s", r.severity)]),
			Cvss:        cvss,
			Confidence:  v1.Confidence_CONFIDENCE_MEDIUM,
			Description: r.description,
			Cve:         r.cve,
		})
	}
	return issues
}

func main() {
	if err := producers.ParseFlags(); err != nil {
		log.Fatal(err)
	}
	jsonBytes, err := ioutil.ReadFile(producers.InResults)
	if err != nil {
		log.Fatal(err)
	}

	issues := UnmarshalJSON(jsonBytes)
	if err := producers.WriteDraconOut(
		"dependencyCheck",
		parseIssues(issues),
	); err != nil {
		log.Fatal(err)
	}
}
