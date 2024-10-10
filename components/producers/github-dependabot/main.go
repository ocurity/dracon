package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"strconv"
	"strings"

	"github.com/google/go-github/v65/github"
	"github.com/package-url/packageurl-go"

	v1proto "github.com/ocurity/dracon/api/proto/v1"
	"github.com/ocurity/dracon/components/producers"
	wrapper "github.com/ocurity/dracon/pkg/github"
)

var (
	// RepositoryOwner is the owner of the GitHub repository
	RepositoryOwner string

	// RepositoryName is the name of the GitHub repository
	RepositoryName string

	// GitHubToken is the GitHub token used to authenticate
	GitHubToken string

	// Ref is the Ref/branch to get alerts for
	Ref string

	// Severity if specified, only code scanning alerts with this severity will be returned. Possible values are: critical, high, medium, low, warning, note, error
	Severity string

	// Ecosystem is a comma separated list of at least one of composer, go, maven, npm, nuget, pip, pub, rubygems, rust
	Ecosystem string
)

func main() {
	flag.StringVar(&RepositoryOwner, "repository-owner", "", "The owner of the GitHub repository")
	flag.StringVar(&RepositoryName, "repository-name", "", "The name of the GitHub repository")
	flag.StringVar(&GitHubToken, "github-token", "", "The GitHub token used to authenticate with the API")
	flag.StringVar(&Ref, "reference", "", "The Ref/branch to get alerts for")
	flag.StringVar(&Severity, "severity", "", "If specified, only code scanning alerts with this severity will be returned. Possible values are: critical, high, medium, low, warning, note, error")
	flag.StringVar(&Ecosystem, "ecosystem", "", "If specified, a comma separated list of at least one of composer, go, maven, npm, nuget, pip, pub, rubygems, rust")
	if err := producers.ParseFlags(); err != nil {
		log.Fatal(err)
	}
	apiClient := wrapper.NewClient(GitHubToken)

	alerts, err := listAlertsForRepo(apiClient, RepositoryOwner, RepositoryName)
	if err != nil {
		log.Fatal(err)
	}

	issues := parseIssues(alerts)

	if err := producers.WriteDraconOut(
		"github-code-scanning",
		issues,
	); err != nil {
		log.Fatal(err)
	}
}

func listAlertsForRepo(apiClient wrapper.Wrapper, owner, repo string) ([]*github.DependabotAlert, error) {
	open := "open"
	opt := &github.ListAlertsOptions{
		State:     &open,
		Severity:  &Severity,
		Ecosystem: &Ecosystem,

		ListOptions: github.ListOptions{
			PerPage: 30,
		},
	}

	var allAlerts []*github.DependabotAlert
	for {
		alerts, resp, err := apiClient.ListRepoDependabotAlerts(context.Background(), owner, repo, opt)
		if err != nil {
			return nil, err
		}

		allAlerts = append(allAlerts, alerts...)

		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}

	slog.Info("Successfully fetched alerts", "count", len(allAlerts), "repository", owner+"/"+repo)

	return allAlerts, nil
}

func parseIssues(alerts []*github.DependabotAlert) []*v1proto.Issue {
	issues := []*v1proto.Issue{}
	for _, alert := range alerts {
		ecosystem := *(alert.GetSecurityVulnerability().Package.Ecosystem)
		if ecosystem == "pip" {
			ecosystem = "pypi"
		}
		cwe := []int32{}
		for _, c := range alert.SecurityAdvisory.CWEs {
			numberOnly := strings.ReplaceAll(*c.CWEID, "CWE-", "")
			cweNum, err := strconv.Atoi(numberOnly)

			if err != nil {
				slog.Error("could not extract cwe number from ", slog.String("cweID", *c.CWEID))
				continue
			}
			cwe = append(cwe, int32(cweNum))
		}
		purl := producers.GetPURLTarget(ecosystem, "", *alert.GetSecurityVulnerability().Package.Name, "", packageurl.Qualifiers{}, "")
		slog.Info("Handling alert for", slog.String("purl", purl))
		cve := ""
		summary := ""
		description := ""
		cvss := float64(0)
		advisory := alert.GetSecurityAdvisory()
		severity := v1proto.Severity_SEVERITY_UNSPECIFIED
		if advisory != nil {
			if advisory.CVEID != nil {
				cve = *advisory.CVEID
			} else if advisory.GHSAID != nil {
				cve = *advisory.GHSAID
			}
			description = *advisory.Description
			summary = *advisory.Summary
			cvss = *advisory.CVSS.Score
			severity = parseGitHubSeverity(*alert.GetSecurityAdvisory().Severity)
		}
		issue := &v1proto.Issue{
			Target:      purl,
			Cve:         cve,
			Title:       summary,
			Description: description,
			Severity:    severity,
			Cvss:        cvss,
			Cwe:         cwe,
		}
		issues = append(issues, issue)
	}

	return issues
}

func parseGitHubSeverity(severity string) v1proto.Severity {
	switch severity {
	case "low":
		return v1proto.Severity_SEVERITY_LOW
	case "medium":
		return v1proto.Severity_SEVERITY_MEDIUM
	case "high":
		return v1proto.Severity_SEVERITY_HIGH
	case "critical":
		return v1proto.Severity_SEVERITY_CRITICAL
	default:
		return v1proto.Severity_SEVERITY_UNSPECIFIED
	}
}
