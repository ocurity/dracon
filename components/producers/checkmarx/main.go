// Package main of the checkmarx producer parses the XML Output of a Checkmarx scan
// creates a Dracon scan from it
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	v1 "github.com/ocurity/dracon/api/proto/v1"
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

	var results Flaws
	if err := producers.ParseXML(inFile, &results); err != nil {
		log.Fatal(err)
	}

	issues, err := parseIssues(&results)
	if err != nil {
		log.Fatal(err)
	}
	if err := producers.WriteDraconOut(
		"checkmarx",
		issues,
	); err != nil {
		log.Fatal(err)
	}
}

func parseIssues(out *Flaws) ([]*v1.Issue, error) {
	issues := []*v1.Issue{}
	source := fmt.Sprintf("%s:%s:%s",
		out.MetaData.AppID,
		out.MetaData.AppName,
		out.MetaData.ComponentName)

	for _, r := range out.Flaw {
		cvss, err := strconv.ParseFloat(r.CVSSScore, 64)
		if err != nil {
			log.Println("could not parse CVSSScore '", r.CVSSScore, "' from Checkmarx scan on", out.MetaData.Date, "CVSS set to 0, err:", err)
		}
		desc := DraconDescription{
			OriginalIssueDescription:                          r.IssueDescription,
			OriginalRemediationAdvice:                         r.RemediationDesc,
			OriginalExploitDescription:                        r.ExploitDesc,
			OriginalDefectInfo:                                r.DefectInfo,
			OriginalNotes:                                     r.Notes,
			OriginalTrace:                                     r.Trace,
			OriginalKnowledgeBaseReference:                    r.KBReference,
			OriginalRelatedExploitRange:                       r.RelatedExploitRange,
			OriginalAttackComplexity:                          r.AttackComplexity,
			OriginalLevelofAuthenticationNeeded:               r.LevelofAuthenticationNeeded,
			OriginalConfidentialityImpact:                     r.ConfidentialityImpact,
			OriginalIntegrityImpact:                           r.IntegrityImpact,
			OriginalAvailabilityImpact:                        r.AvailabilityImpact,
			OriginalCollateralDamagePotential:                 r.CollateralDamagePotential,
			OriginalTargetDistribution:                        r.TargetDistribution,
			OriginalConfidentialityRequirement:                r.ConfidentialityRequirement,
			OriginalIntegrityRequirement:                      r.IntegrityRequirement,
			OriginalAvailabilityRequirement:                   r.AvailabilityRequirement,
			OriginalAvailabilityofExploit:                     r.AvailabilityofExploit,
			OriginalTypeofFixAvailable:                        r.TypeofFixAvailable,
			OriginalLevelofVerificationthatVulnerabilityExist: r.LevelofVerificationthatVulnerabilityExist,
		}
		draconDesc, err := json.Marshal(desc)
		if err != nil {
			log.Println("Could not populate Dracon Description from Checkmarx fields, err", err)
		}
		target := fmt.Sprintf("%s:%s", r.FileName, r.LineNumber)
		iss := &v1.Issue{
			Source:      source,
			Target:      target,
			Type:        r.VulnerabilityType,
			Title:       fmt.Sprintf("%s - %s - %s", r.Status, r.Severity, r.ComponentName),
			Severity:    v1.Severity(v1.Severity_value[fmt.Sprintf("SEVERITY_%s", strings.ToUpper(r.Severity))]),
			Cvss:        cvss,
			Confidence:  v1.Confidence(v1.Confidence_value["CONFIDENCE_UNSPECIFIED"]),
			Description: string(draconDesc),
		}
		issues = append(issues, iss)
	}
	return issues, nil
}

// DraconDescription allows the user to map Checkmarx optional fields to the Dracon "description" field
type DraconDescription struct {
	OriginalIssueDescription                          string `json:"issue description,omitempty"`
	OriginalRemediationAdvice                         string `json:"remediation advice,omitempty"`
	OriginalExploitDescription                        string `json:"exploit description,omitempty"`
	OriginalDefectInfo                                string `json:"defect info,omitempty"`
	OriginalNotes                                     string `json:"notes,omitempty"`
	OriginalTrace                                     string `json:"trace,omitempty"`
	OriginalKnowledgeBaseReference                    string `json:"knowledge base reference,omitempty"`
	OriginalRelatedExploitRange                       string `json:"related exploit range,omitempty"`
	OriginalAttackComplexity                          string `json:"attack complexitt,omitempty"`
	OriginalLevelofAuthenticationNeeded               string `json:"level of authentication needed,omitempty"`
	OriginalConfidentialityImpact                     string `json:"confidentiality impact,omitempty"`
	OriginalIntegrityImpact                           string `json:"integrity impact,omitempty"`
	OriginalAvailabilityImpact                        string `json:"availability impact,omitempty"`
	OriginalCollateralDamagePotential                 string `json:"collateral damage potential,omitempty"`
	OriginalTargetDistribution                        string `json:"taret distribution,omitempty"`
	OriginalConfidentialityRequirement                string `json:"confidentiality requirement,omitempty"`
	OriginalIntegrityRequirement                      string `json:"integrity requirement,omitempty"`
	OriginalAvailabilityRequirement                   string `json:"availability requirement,omitempty"`
	OriginalAvailabilityofExploit                     string `json:"availability of exploit,omitempty"`
	OriginalTypeofFixAvailable                        string `json:"type of fix available,omitempty"`
	OriginalLevelofVerificationthatVulnerabilityExist string `json:"level of verification that vulnerability exists,omitempty"`
}

// Flaws is the checkmarx output xml
type Flaws struct {
	MetaData struct {
		AppID         string `xml:"appID,attr" json:"appid,omitempty"`
		AppName       string `xml:"appName,attr" json:"appname,omitempty"`
		ComponentName string `xml:"componentName,attr" json:"componentname,omitempty"`
		Date          string `xml:"date,attr" json:"date,omitempty"`
		ReleaseName   string `xml:"releaseName,attr" json:"releasename,omitempty"`
		SourceName    string `xml:"sourceName,attr" json:"sourcename,omitempty"`
		SourceDesc    string `xml:"sourceDesc,attr" json:"sourcedesc,omitempty"`
	} `xml:"metaData" json:"metadata,omitempty"`
	Flaw []struct {
		Text                                      string `xml:",chardata" json:"text,omitempty"`
		ID                                        string `xml:"id"`
		Status                                    string `xml:"status"`
		IssueDescription                          string `xml:"issueDescription"`
		RemediationDesc                           string `xml:"remediationDesc"`
		ExploitDesc                               string `xml:"exploitDesc"`
		IssueRecommendation                       string `xml:"issueRecommendation"`
		ComponentName                             string `xml:"componentName"`
		Module                                    string `xml:"module"`
		APIName                                   string `xml:"apiName"`
		VulnerabilityType                         string `xml:"vulnerabilityType"` // Basically CWE
		Classification                            string `xml:"classification"`
		Severity                                  string `xml:"severity"`
		FileName                                  string `xml:"fileName"`
		LineNumber                                string `xml:"lineNumber"`
		SrcContext                                string `xml:"srcContext"`
		DefectInfo                                string `xml:"defectInfo"`
		Notes                                     string `xml:"notes"`
		Trace                                     string `xml:"trace"`
		CallerName                                string `xml:"callerName"`
		FindingCodeRegion                         string `xml:"findingCodeRegion"`
		DateFirstOccurrence                       string `xml:"dateFirstOccurrence"`
		IssueBornDate                             string `xml:"issueBornDate"`
		IssueName                                 string `xml:"issueName"`
		KBReference                               string `xml:"kBReference"`
		CVSSScore                                 string `xml:"cVSSScore"`
		RelatedExploitRange                       string `xml:"relatedExploitRange"`
		AttackComplexity                          string `xml:"attackComplexity"`
		LevelofAuthenticationNeeded               string `xml:"levelofAuthenticationNeeded"`
		ConfidentialityImpact                     string `xml:"confidentialityImpact"`
		IntegrityImpact                           string `xml:"integrityImpact"`
		AvailabilityImpact                        string `xml:"availabilityImpact"`
		CollateralDamagePotential                 string `xml:"collateralDamagePotential"`
		TargetDistribution                        string `xml:"targetDistribution"`
		ConfidentialityRequirement                string `xml:"confidentialityRequirement"`
		IntegrityRequirement                      string `xml:"integrityRequirement"`
		AvailabilityRequirement                   string `xml:"availabilityRequirement"`
		AvailabilityofExploit                     string `xml:"availabilityofExploit"`
		TypeofFixAvailable                        string `xml:"typeofFixAvailable"`
		LevelofVerificationthatVulnerabilityExist string `xml:"levelofVerificationthatVulnerabilityExist"`
		CVSSEquation                              string `xml:"cVSSEquation"`
	} `xml:"flaw" json:"flaw"`
}
