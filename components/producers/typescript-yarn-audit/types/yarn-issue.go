package types

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	v1 "github.com/ocurity/dracon/api/proto/v1"
)

func yarnToIssueSeverity(severity string) v1.Severity {
	switch severity {
	case "low":
		return v1.Severity_SEVERITY_LOW
	case "moderate":
		return v1.Severity_SEVERITY_MEDIUM
	case "high":
		return v1.Severity_SEVERITY_HIGH
	case "critical":
		return v1.Severity_SEVERITY_CRITICAL
	default:
		return v1.Severity_SEVERITY_INFO

	}
}

type yarnAuditLine struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func (yl *yarnAuditLine) UnmarshalJSON(data []byte) error {
	var typ struct {
		Type string `json:"type"`
	}

	if err := json.Unmarshal(data, &typ); err != nil {
		return err
	}

	switch typ.Type {
	case "auditSummary":
		yl.Data = new(auditSummaryData)
	case "auditAdvisory":
		yl.Data = new(auditAdvisoryData)
	case "auditAction":
		yl.Data = new(auditActionData)
	default:
		log.Printf("Parsed unsupported type: %s", typ.Type)
	}

	type tmp yarnAuditLine // avoids infinite recursion
	return json.Unmarshal(data, (*tmp)(yl))
}

type auditActionData struct {
	Cmd        string      `json:"cmd"`
	IsBreaking bool        `json:"isBreaking"`
	Action     auditAction `json:"action"`
}

type auditAdvisoryData struct {
	Resolution auditResolution `json:"resolution"`
	Advisory   yarnAdvisory    `json:"advisory"`
}

// AsIssue returns data as a Dracon v1.Issue.
func (audit *auditAdvisoryData) AsIssue() *v1.Issue {
	var targetName string
	if audit.Resolution.Path != "" {
		targetName = audit.Resolution.Path + ": "
	}
	targetName += audit.Advisory.ModuleName

	return &v1.Issue{
		Target:      targetName,
		Type:        audit.Advisory.Cwe,
		Title:       audit.Advisory.Title,
		Severity:    yarnToIssueSeverity(audit.Advisory.Severity),
		Confidence:  v1.Confidence_CONFIDENCE_HIGH,
		Description: audit.Advisory.GetDescription(),
		Cve:         strings.Join(audit.Advisory.Cves, ", "),
	}
}

type auditSummaryData struct {
	Vulnerabilities      vulnerabilities `json:"vulnerabilities"`
	Dependencies         int             `json:"dependencies"`
	DevDependencies      int             `json:"devDependencies"`
	OptionalDependencies int             `json:"optionalDependencies"`
	TotalDependencies    int             `json:"totalDependencies"`
}

type auditAction struct {
	Action   string            `json:"action"`
	Module   string            `json:"module"`
	Target   string            `json:"target"`
	IsMajor  bool              `json:"isMajor"`
	Resolves []auditResolution `json:"resolves"`
}

type vulnerabilities struct {
	Info     int `json:"info"`
	Low      int `json:"low"`
	Moderate int `json:"moderate"`
	High     int `json:"high"`
	Critical int `json:"critical"`
}

type yarnAdvisory struct {
	Findings           []finding         `json:"findings"`
	Metadata           *advisoryMetaData `json:"metadata"`
	VulnerableVersions string            `json:"vulnerable_versions"`
	ModuleName         string            `json:"module_name"`
	Severity           string            `json:"severity"`
	GithubAdvisoryID   string            `json:"github_advisory_id"`
	Cves               []string          `json:"cves"`
	Access             string            `json:"access"`
	PatchedVersions    string            `json:"patched_versions"`
	Updated            string            `json:"updated"`
	Recommendation     string            `json:"recommendation"`
	Cwe                string            `json:"cwe"`
	FoundBy            *contact          `json:"found_by"`
	Deleted            bool              `json:"deleted"`
	ID                 int               `json:"id"`
	References         string            `json:"references"`
	Created            string            `json:"created"`
	ReportedBy         *contact          `json:"reported_by"`
	Title              string            `json:"title"`
	NpmAdvisoryID      interface{}       `json:"npm_advisory_id"`
	Overview           string            `json:"overview"`
	URL                string            `json:"url"`
}

func (advisory *yarnAdvisory) GetDescription() string {
	return fmt.Sprintf(
		"Vulnerable Versions: %s\nRecommendation: %s\nOverview: %s\nReferences:\n%s\nAdvisory URL: %s\n",
		advisory.VulnerableVersions,
		advisory.Recommendation,
		advisory.Overview,
		advisory.References,
		advisory.URL,
	)
}

type finding struct {
	Version  string   `json:"version"`
	Paths    []string `json:"paths"`
	Dev      bool     `json:"dev"`
	Optional bool     `json:"optional"`
	Bundled  bool     `json:"bundled"`
}

type auditResolution struct {
	ID       int    `json:"id"`
	Path     string `json:"path"`
	Dev      bool   `json:"dev"`
	Optional bool   `json:"optional"`
	Bundled  bool   `json:"bundled"`
}

type advisoryMetaData struct {
	ModuleType         string `json:"module_type"`
	Exploitability     int    `json:"exploitability"`
	AffectedComponents string `json:"affected_components"`
}

type contact struct {
	Name string `json:"name"`
}

// YarnAuditReport includes yarn audit data grouped by advisories, actions and summary.
type YarnAuditReport struct {
	AuditAdvisories []*auditAdvisoryData
	AuditActions    []*auditActionData
	AuditSummary    *auditSummaryData
}

// NewReport returns a YarnAuditReport, assuming each line is jsonline and returns any errors.
func NewReport(reportLines [][]byte) (*YarnAuditReport, []error) {
	var report YarnAuditReport

	var errors []error

	for _, line := range reportLines {
		var auditLine yarnAuditLine
		if err := json.Unmarshal(line, &auditLine); err != nil {
			log.Printf("Error parsing JSON line '%s': %s\n", line, err)
			errors = append(errors, err)
		} else {
			switch x := auditLine.Data.(type) {
			case *auditSummaryData:
				report.AuditSummary = x
			case *auditAdvisoryData:
				report.AuditAdvisories = append(report.AuditAdvisories, x)
			case *auditActionData:
				report.AuditActions = append(report.AuditActions, x)
			}
		}
	}

	if report.AuditAdvisories != nil && len(report.AuditAdvisories) > 0 {
		return &report, errors
	}

	return nil, errors
}

// AsIssues returns the YarnAuditReport as Dracon v1.Issue list. Currently only converts the YarnAuditReport.AuditAdvisories.
func (r *YarnAuditReport) AsIssues() []*v1.Issue {
	issues := make([]*v1.Issue, 0)

	for _, audit := range r.AuditAdvisories {
		issues = append(issues, audit.AsIssue())
	}

	return issues
}
