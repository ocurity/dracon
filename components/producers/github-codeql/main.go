package main

import (
	"context"
	"flag"
	"io"
	"log"
	"log/slog"
	"strconv"
	"strings"

	"github.com/go-errors/errors"
	"github.com/google/go-github/v65/github"

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

	// toolName is internal
	toolName string
)

func main() {
	flag.StringVar(&RepositoryOwner, "repository-owner", "", "The owner of the GitHub repository")
	flag.StringVar(&RepositoryName, "repository-name", "", "The name of the GitHub repository")
	flag.StringVar(&GitHubToken, "github-token", "", "The GitHub token used to authenticate with the API")
	flag.StringVar(&Ref, "reference", "", "The Ref/branch to get alerts for")
	flag.StringVar(&Severity, "severity", "", "If specified, only code scanning alerts with this severity will be returned. Possible values are: critical, high, medium, low, warning, note, error")
	toolName = "CodeQL"
	if err := producers.ParseFlags(); err != nil {
		log.Fatal(err)
	}
	apiClient := wrapper.NewClient(GitHubToken)
	alerts, err := listAlertsForRepo(apiClient, RepositoryOwner, RepositoryName, toolName)
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

func listAlertsForRepo(apiClient wrapper.Wrapper, owner, repo, toolName string) ([]*github.Alert, error) {
	var severity string
	if Severity != "" {
		severity = Severity
	}
	opt := &github.AlertListOptions{
		State:    "open",
		Ref:      Ref,
		Severity: severity,
		ToolName: toolName,
		ListOptions: github.ListOptions{
			PerPage: 30,
		},
	}

	var allAlerts []*github.Alert
	for {
		alerts, resp, err := apiClient.ListRepoAlerts(context.Background(), owner, repo, opt)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode > 299 {
			body, bodyReadErr := io.ReadAll(resp.Body)
			if bodyReadErr != nil {
				return nil, errors.Errorf("could not read response error from github, err:%w", bodyReadErr)
			}
			return nil, errors.Errorf("did not receive a valid status code from github, status code: %d, body:%s", resp.StatusCode, body)
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

func parseIssues(alerts []*github.Alert) []*v1proto.Issue {
	issues := []*v1proto.Issue{}
	for _, alert := range alerts {

		issue := &v1proto.Issue{
			Target: producers.GetFileTarget(
				alert.GetMostRecentInstance().GetLocation().GetPath(),
				alert.GetMostRecentInstance().GetLocation().GetStartLine(),
				alert.GetMostRecentInstance().GetLocation().GetEndLine(),
			),
			Type:        strconv.Itoa(alert.GetNumber()),
			Title:       *alert.GetRule().Description,
			Severity:    parseGitHubSeverity(*alert.GetRule().Severity),
			Cvss:        0,
			Confidence:  v1proto.Confidence_CONFIDENCE_UNSPECIFIED,
			Description: alert.GetMostRecentInstance().GetMessage().GetText(),
			Source:      alert.GetHTMLURL(),
			Cwe:         parseGithubCWEsFromTags(alert.Rule.Tags),
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
