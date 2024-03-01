package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	v1 "github.com/ocurity/dracon/api/proto/v1"

	"github.com/ocurity/dracon/components/producers"
	packageurl "github.com/package-url/packageurl-go"
)

func main() {
	if err := producers.ParseFlags(); err != nil {
		log.Fatal(err)
	}

	inFile, err := producers.ReadInFile()
	if err != nil {
		log.Fatal(err)
	}

	var results BlackduckOut
	if err := json.Unmarshal(inFile, &results); err != nil {
		log.Fatal(err)
	}

	issues, err := parseIssues(&results)
	if err != nil {
		log.Fatal(err)
	}
	if err := producers.WriteDraconOut(
		"blackduck",
		issues,
	); err != nil {
		log.Fatal(err)
	}
}

func cweIdsToInt(ids []string) []int32 {
	res := []int32{}
	for _, id := range ids {
		newID := strings.ReplaceAll(id, "CWE-", "")
		newIntID, err := strconv.Atoi(newID)
		if err != nil {
			log.Println(err)
			continue
		}
		res = append(res, int32(newIntID))
	}
	return res
}

func BDTargetToPurl(originName, originID string) string {
	switch originName {
	case "anaconda":
		originName = "conda"
	case "centos":
		originName = "rpm"
	case "debian":
		originName = "deb"
	case "fedora":
		originName = "rpm"
	case "gradle":
		originName = "maven"
	case "npmjs":
		originName = "npm"
	case "opensuse":
		originName = "rpm"
	case "packagist":
		originName = "composer"
	case "redhat":
		originName = "rpm"
	}
	found := false
	for knownType := range packageurl.KnownTypes {
		if originName == knownType {
			found = true
		}
	}
	for knownType := range packageurl.CandidateTypes {
		if originName == knownType {
			found = true
		}
	}
	if !found {
		log.Println("Cannot create package URL from package of type", originName, originID)
	}
	splitOrigin := strings.Split(originID, ":")
	if len(splitOrigin) == 2 { // no namespace e.g. npm/reactjs
		return packageurl.NewPackageURL(originName, splitOrigin[0], "", splitOrigin[1], packageurl.Qualifiers{}, "").ToString()
	} else if len(splitOrigin) == 3 {
		return packageurl.NewPackageURL(originName, splitOrigin[0], splitOrigin[1], splitOrigin[2], packageurl.Qualifiers{}, "").ToString()
	} else {
		log.Println("originID was split into", len(splitOrigin), "parts, this looks like a bug, for reference, original originID was", originID)
	}
	return originName + "/" + originID
}

func parseIssues(out *BlackduckOut) ([]*v1.Issue, error) {
	issues := []*v1.Issue{}
	for _, r := range out.Data {
		for _, vuln := range r.Vulnerabilities.Items {
			cwe := cweIdsToInt(vuln.CweIds)
			severity := v1.Severity_SEVERITY_UNSPECIFIED
			if vuln.UseCvss3 {
				switch vuln.Cvss3.Severity {
				case "CRITICAL":
					severity = v1.Severity_SEVERITY_CRITICAL
				case "HIGH":
					severity = v1.Severity_SEVERITY_HIGH
				case "MEDIUM":
					severity = v1.Severity_SEVERITY_MEDIUM
				case "LOW":
					severity = v1.Severity_SEVERITY_LOW
				case "INFO":
					severity = v1.Severity_SEVERITY_INFO
				case "UNASSIGNED":
					severity = v1.Severity_SEVERITY_UNSPECIFIED
				}
			}
			cve := ""
			if !strings.HasPrefix(vuln.Name, "CVE") {
				for _, metaLink := range vuln.Meta.Links {
					if strings.HasPrefix(metaLink.Label, "CVE") {
						cve = metaLink.Label
					}
				}
			} else {
				cve = vuln.Name
			}
			description := fmt.Sprintf("%s\nSolution Available: %t\nWorkaround Available: %t\nExploit Available: %t\nOriginal Description: %s", vuln.Summary, vuln.SolutionAvailable, vuln.WorkaroundAvailable, vuln.ExploitAvailable, vuln.Description)
			iss := &v1.Issue{
				Cvss:        vuln.Cvss3.OverallScore,
				Cwe:         cwe,
				Cve:         cve,
				Target:      BDTargetToPurl(r.Origins.OriginName, r.Origins.OriginID),
				Type:        vuln.ID,
				Title:       vuln.Summary,
				Severity:    severity,
				Confidence:  v1.Confidence_CONFIDENCE_UNSPECIFIED,
				Description: description,
			}
			issues = append(issues, iss)

		}
	}

	return issues, nil
}

// BlackduckOut models the output of a blackduck scan
type BlackduckOut struct {
	TotalComponentsFound int      `json:"totalComponentsFound,omitempty"`
	MatchedFlag          bool     `json:"matchedFlag,omitempty"`
	Headers              []string `json:"headers,omitempty"`
	Appid                string   `json:"appid,omitempty"`
	Appname              string   `json:"appname,omitempty"`
	Releaseid            string   `json:"releaseid,omitempty"`
	ProjectID            string   `json:"projectId,omitempty"`
	Data                 []struct {
		ComponentName    string `json:"componentName,omitempty"`
		ComponentVersion string `json:"componentVersion,omitempty"`
		Componentid      string `json:"componentid,omitempty"`
		MatchedFiles     struct {
			Items []struct {
				FilePath struct {
					Path                 string `json:"path,omitempty"`
					ArchiveContext       string `json:"archiveContext,omitempty"`
					CompositePathContext string `json:"compositePathContext,omitempty"`
					FileName             string `json:"fileName,omitempty"`
				} `json:"filePath,omitempty"`
				Usages []string `json:"usages,omitempty"`
			} `json:"items,omitempty"`
		} `json:"matchedFiles,omitempty"`
		Vulnerabilities struct {
			Items []struct {
				ID                string    `json:"id,omitempty"`
				Summary           string    `json:"summary,omitempty"`
				PublishedDate     time.Time `json:"publishedDate,omitempty"`
				LastModified      time.Time `json:"lastModified,omitempty"`
				Source            string    `json:"source,omitempty"`
				RemediationStatus string    `json:"remediationStatus,omitempty"`
				CreatedAt         time.Time `json:"createdAt,omitempty"`
				UpdatedAt         time.Time `json:"updatedAt,omitempty"`
				CreatedBy         struct {
					UserName  string `json:"userName,omitempty"`
					FirstName string `json:"firstName,omitempty"`
					LastName  string `json:"lastName,omitempty"`
					User      string `json:"user,omitempty"`
				} `json:"createdBy,omitempty"`
				UpdatedBy struct {
					UserName  string `json:"userName,omitempty"`
					FirstName string `json:"firstName,omitempty"`
					LastName  string `json:"lastName,omitempty"`
					User      string `json:"user,omitempty"`
				} `json:"updatedBy,omitempty"`
				CweIds []string `json:"cweIds,omitempty"`
				Cvss3  struct {
					Severity     string  `json:"severity,omitempty"`
					Vector       string  `json:"vector,omitempty"`
					OverallScore float64 `json:"overallScore,omitempty"`
				} `json:"cvss3,omitempty"`
				UseCvss3            bool     `json:"useCvss3,omitempty"`
				OverallScore        float64  `json:"overallScore,omitempty"`
				SolutionAvailable   bool     `json:"solutionAvailable,omitempty"`
				WorkaroundAvailable bool     `json:"workaroundAvailable,omitempty"`
				ExploitAvailable    bool     `json:"exploitAvailable,omitempty"`
				BdsaTags            []string `json:"bdsaTags,omitempty"`
				Meta                struct {
					Allow []string `json:"allow,omitempty"`
					Href  string   `json:"href,omitempty"`
					Links []struct {
						Rel   string `json:"rel,omitempty"`
						Href  string `json:"href,omitempty"`
						Label string `json:"label,omitempty"`
					} `json:"links,omitempty"`
				} `json:"_meta,omitempty"`
				Name        string `json:"name,omitempty"`
				Description string `json:"description,omitempty"`
			} `json:"items,omitempty"`
			AppliedFilters []any `json:"appliedFilters,omitempty"`
		} `json:"vulnerabilities,omitempty"`
		Policyviolations struct{} `json:"policyviolations,omitempty"`
		Origins          struct {
			Name                          string `json:"name,omitempty"`
			Origin                        string `json:"origin,omitempty"`
			ExternalNamespace             string `json:"externalNamespace,omitempty"`
			ExternalID                    string `json:"externalId,omitempty"`
			ExternalNamespaceDistribution bool   `json:"externalNamespaceDistribution,omitempty"`
			Meta                          struct {
				Allow []any `json:"allow,omitempty"`
				Links []struct {
					Rel  string `json:"rel,omitempty"`
					Href string `json:"href,omitempty"`
				} `json:"links,omitempty"`
			} `json:"_meta,omitempty"`
			OriginName string `json:"originName,omitempty"`
			OriginID   string `json:"originId,omitempty"`
		} `json:"origins,omitempty"`
	} `json:"data,omitempty"`
}
