package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"strconv"
	"strings"

	"github.com/google/go-github/v63/github"
	"golang.org/x/oauth2"

	v1 "github.com/ocurity/dracon/api/proto/v1"
	"github.com/ocurity/dracon/components/producers"
)

var (
	// RepositoryOwner is the owner of the GitHub repository
	RepositoryOwner string

	// RepositoryName is the name of the GitHub repository
	RepositoryName string

	// GitHubToken is the GitHub token used to authenticate
	GitHubToken string
)

func main() {
	flag.StringVar(&RepositoryOwner, "repository-owner", "", "The owner of the GitHub repository")
	flag.StringVar(&RepositoryName, "repository-name", "", "The name of the GitHub repository")
	flag.StringVar(&GitHubToken, "github-token", "", "The GitHub token used to authenticate with the API")
	if err := producers.ParseFlags(); err != nil {
		log.Fatal(err)
	}

	alerts, err := listAlertsForRepo(RepositoryOwner, RepositoryName, GitHubToken)
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

func listAlertsForRepo(owner, repo, token string) ([]*github.Alert, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	oAuthClient := oauth2.NewClient(context.Background(), ts)
	apiClient := github.NewClient(oAuthClient)

	opt := &github.AlertListOptions{
		State: "open",
		ListOptions: github.ListOptions{
			PerPage: 30,
		},
	}

	var allAlerts []*github.Alert
	for {
		alerts, resp, err := apiClient.CodeScanning.ListAlertsForRepo(context.Background(), owner, repo, opt)
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

func parseIssues(alerts []*github.Alert) []*v1.Issue {
	issues := []*v1.Issue{}
	for _, alert := range alerts {

		issue := &v1.Issue{
			Target: producers.GetFileTarget(
				alert.GetMostRecentInstance().GetLocation().GetPath(),
				alert.GetMostRecentInstance().GetLocation().GetStartLine(),
				alert.GetMostRecentInstance().GetLocation().GetEndLine(),
			),
			Type:        strconv.Itoa(alert.GetNumber()),
			Title:       *alert.GetRule().Description,
			Severity:    parseGitHubSeverity(*alert.GetRule().Severity),
			Cvss:        0,
			Confidence:  v1.Confidence_CONFIDENCE_UNSPECIFIED,
			Description: alert.GetMostRecentInstance().GetMessage().GetText(),
			Source:      alert.GetHTMLURL(),
			Cwe:         parseGithubCWEsFromTags(alert.Rule.Tags),
		}
		issues = append(issues, issue)
	}

	return issues
}

func parseGitHubSeverity(severity string) v1.Severity {
	switch severity {
	case "low":
		return v1.Severity_SEVERITY_LOW
	case "medium":
		return v1.Severity_SEVERITY_MEDIUM
	case "high":
		return v1.Severity_SEVERITY_HIGH
	case "critical":
		return v1.Severity_SEVERITY_CRITICAL
	default:
		return v1.Severity_SEVERITY_UNSPECIFIED
	}
}

func parseGithubCWEsFromTags(tags []string) []int32 {
	// example input: ["security", "external/cwe/cwe-022"]
	cwePrefix := "external/cwe/cwe-"

	cwes := []int32{}
	for _, item := range tags {
		if strings.HasPrefix(item, cwePrefix) {
			cweString := strings.TrimPrefix(item, cwePrefix)

			cwe, err := strconv.Atoi(cweString)
			if err != nil {
				slog.Warn("Failed to parse CWE from tag", "tag", item, "error", err)
				continue
			}

			cwes = append(cwes, int32(cwe))
		}
	}
	return cwes
}
