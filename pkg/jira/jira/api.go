// Package jira holds Jira API utilities, used to get data in/out of the target jira instance
package jira

import (
	"fmt"
	"io/ioutil"
	"log"

	jira "github.com/andygrunwald/go-jira"

	"github.com/ocurity/dracon/pkg/jira/config"
	"github.com/ocurity/dracon/pkg/jira/document"
)

// Client is a wrapper of a go-jira client with our config on top.
type Client struct {
	JiraClient    *jira.Client
	DryRunMode    bool
	Config        config.Config
	DefaultFields defaultJiraFields
}

// NewClient returns a client containing the authentication details and the configuration settings.
func NewClient(user, token, url string, dryRun bool, config config.Config) *Client {
	return &Client{
		JiraClient:    authJiraClient(user, token, url),
		DryRunMode:    dryRun,
		Config:        config,
		DefaultFields: getDefaultFields(config),
	}
}

// authJiraClient authenticates the client with the given Username, API token, and URL domain.
func authJiraClient(user, token, url string) *jira.Client {
	tp := jira.BasicAuthTransport{
		Username: user,
		Password: token,
	}
	JiraClientlient, err := jira.NewClient(tp.Client(), url)
	if err != nil {
		log.Fatalf("Unable to contact Jira: %s", err)
	}
	return JiraClientlient
}

// assembleIssue parses the Dracon message and serializes it into a Jira Issue object.
func (c Client) assembleIssue(draconResult document.Document) *jira.Issue {
	// Mappings the Dracon Result fields to their corresponding Jira fields specified in the configuration
	customFields := c.DefaultFields.CustomFields.Clone()

	for _, m := range c.Config.Mappings {
		strMap, _ := draconResultToSTRMaps(draconResult)
		if _, ok := draconResult.Annotations[m.DraconField]; ok {
			customFields[m.JiraField] = makeCustomField(m.FieldType, []string{draconResult.Annotations[m.DraconField]})
		} else {
			customFields[m.JiraField] = makeCustomField(m.FieldType, []string{strMap[m.DraconField]})
		}
	}
	summary, extra := makeSummary(draconResult)
	description := makeDescription(draconResult, c.Config.DescriptionExtras)
	if extra != "" {
		description = fmt.Sprintf(".... %s\n%s", extra, description)
	}
	iss := &jira.Issue{
		Fields: &jira.IssueFields{
			Project:     c.DefaultFields.Project,
			Type:        c.DefaultFields.IssueType,
			Description: description,
			Summary:     summary,
			Unknowns:    customFields,
		},
	}
	if len(c.DefaultFields.Components) != 0 {
		iss.Fields.Components = c.DefaultFields.Components
	}
	if len(c.DefaultFields.AffectsVersions) != 0 {
		iss.Fields.AffectsVersions = c.DefaultFields.AffectsVersions
	}
	if len(c.DefaultFields.Labels) != 0 {
		iss.Fields.Labels = c.DefaultFields.Labels
	}
	return iss
}

// CreateIssue creates a new issue in Jira.
func (c Client) CreateIssue(draconResult document.Document) error {
	issue := c.assembleIssue(draconResult)

	if c.DryRunMode {
		log.Printf("Dry run mode. The following issue would have been created: '%s'", issue.Fields.Summary)
		return nil
	}

	ri, resp, err := c.JiraClient.Issue.Create(issue)
	if err != nil {
		if resp != nil {
			body, _ := ioutil.ReadAll(resp.Body)
			log.Printf("Error occurred posting to Jira. Response body:\n%s", body)
		} else {
			log.Println("create issue jira response is nil, should be anything but that")
		}
		return err
	}
	log.Printf("Created Jira Issue ID %s. jira_key=%s", ri.ID, string(ri.Key))
	return nil
}

// SearchByJQL searches jira instance by JQL and returns results with history.
func (c Client) SearchByJQL(jql string) ([]jira.Issue, error) {
	var results []jira.Issue
	startAt := 0
	maxresults := 100
	expand := "names,schema,operations,editmeta,changelog,renderedFields"
	issues, response, err := c.JiraClient.Issue.Search(jql, &jira.SearchOptions{Expand: expand, StartAt: startAt, MaxResults: maxresults}) // maxresults is capped to 100 by attlasian
	if err != nil {
		log.Print(response)
		return nil, err
	}
	results = append(results, issues...)
	startAt = len(results)
	log.Print("The query returned ", response.Total, " results")
	for len(results) < response.Total {
		issues, response, err = c.JiraClient.Issue.Search(jql, &jira.SearchOptions{Expand: expand, StartAt: startAt, MaxResults: maxresults}) // maxresults is capped to 100 by attlasian
		if err != nil {
			log.Print(response)
			return nil, err
		}
		results = append(results, issues...)
		startAt = len(results)
	}
	return results, nil
}
