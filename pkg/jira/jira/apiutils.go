package jira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	jira "github.com/andygrunwald/go-jira"
	"github.com/trivago/tgo/tcontainer"

	v1 "github.com/ocurity/dracon/api/proto/v1"
	"github.com/ocurity/dracon/pkg/enumtransformers"
	"github.com/ocurity/dracon/pkg/jira/config"
	"github.com/ocurity/dracon/pkg/jira/document"
	"github.com/ocurity/dracon/pkg/templating"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type defaultJiraFields struct {
	Project         jira.Project
	IssueType       jira.IssueType
	Components      []*jira.Component
	AffectsVersions []*jira.AffectsVersion
	Labels          []string
	CustomFields    tcontainer.MarshalMap
}

// getDefaultFields creates the fields for Project, IssueType, Components, AffectsVersions, Labels and CustomFields
// with the default values specified in config.yaml and serializes them into Jira Fields.
func getDefaultFields(config config.Config) defaultJiraFields {
	defaultFields := defaultJiraFields{}
	defaultFields.Project = jira.Project{
		Key: config.DefaultValues.Project,
	}

	defaultFields.IssueType = jira.IssueType{
		Name: config.DefaultValues.IssueType,
	}

	components := []*jira.Component{}
	for _, v := range config.DefaultValues.Components {
		components = append(components, &jira.Component{Name: v})
	}
	defaultFields.Components = components

	affectsVersions := []*jira.AffectsVersion{}
	for _, v := range config.DefaultValues.AffectsVersions {
		affectsVersions = append(affectsVersions, &jira.AffectsVersion{Name: v})
	}
	defaultFields.AffectsVersions = affectsVersions

	defaultFields.Labels = config.DefaultValues.Labels

	customFields := tcontainer.NewMarshalMap()
	for _, cf := range config.DefaultValues.CustomFields {
		customFields[cf.ID] = makeCustomField(cf.FieldType, cf.Values)
	}
	defaultFields.CustomFields = customFields

	return defaultFields
}

// makeCustomField returns the appropriate interface for a jira CustomField given it's type and values
// :param fieldType: the type of the field in Jira (single-value, multi-value, float)
// :param values: list of values to be filled in
// :return the appropriate interface for a CustomField, given the corresponding fieldType and value(s).
func makeCustomField(fieldType string, values []string) interface{} {
	switch fieldType {
	case "single-value":
		return map[string]string{"value": values[0]}
	case "multi-value":
		cf := []map[string]string{}
		for _, v := range values {
			cf = append(cf, map[string]string{"value": v})
		}
		return cf
	case "float":
		f, err := strconv.ParseFloat(values[0], 64)
		if err != nil {
			log.Fatalf("Error parsing float field-type: %v", err)
		}
		return f
	case "simple-value":
		return values[0]
	default:
		log.Printf("Warning: Field type %s is not supported. Edit your config.yaml file, as this field will not be displayed correctly.", fieldType)
		return nil
	}
}

func draconResultToSTRMaps(draconResult document.Document) (map[string]string, string) {
	var strMap map[string]string

	annotations, err := json.Marshal(draconResult.Annotations)
	if err != nil {
		log.Fatalf("could not marshal annotations: %s", err)
	}
	draconResult.Annotations = nil
	tmp, err := json.Marshal(draconResult)
	if err != nil {
		log.Fatalf("could not marshal result: %s", err)
	}
	if err := json.Unmarshal(tmp, &strMap); err != nil {
		log.Fatalf("could not unmarshal result: %s", err)
	}
	return strMap, string(annotations)
}

// makeDescription creates the description of an issue's enhanced with extra information from the Dracon Result.
func makeDescription(draconResult document.Document, extras []string, template string) string {
	if draconResult.Count == "" {
		draconResult.Count = "0"
	}
	count, err := strconv.Atoi(draconResult.Count)
	if err != nil {
		log.Fatal("could not template enriched issue ", err)
	}
	fp := false
	if strings.ToLower(draconResult.FalsePositive) == "true" {
		fp = true
	}
	if draconResult.CVSS == "" {
		draconResult.CVSS = "0.0"
	}
	cvss, err := strconv.ParseFloat(draconResult.CVSS, 64)
	if err != nil {
		log.Fatal("could not template enriched issue ", err)
	}
	description, err := templating.TemplateStringEnriched(template,
		&v1.EnrichedIssue{
			Annotations:   draconResult.Annotations,
			Count:         uint64(count),
			FalsePositive: fp,
			FirstSeen:     timestamppb.New(draconResult.FirstFound),
			Hash:          draconResult.Hash,
			RawIssue: &v1.Issue{
				Confidence:  enumtransformers.TextToConfidence(draconResult.ConfidenceText),
				Cve:         draconResult.CVE,
				Cvss:        cvss,
				Description: draconResult.Description,
				Severity:    enumtransformers.TextToSeverity(draconResult.SeverityText),
				Source:      draconResult.Source,
				Target:      draconResult.Target,
				Title:       draconResult.Title,
				Type:        draconResult.Type,
			},
		},
	)
	if err != nil {
		log.Fatal("Could not template enriched issue ", err)
	}
	desc := *description
	// Append the extra fields to the description
	strMap, annotations := draconResultToSTRMaps(draconResult)
	if len(extras) > 0 {
		desc += "{code:}" + "\n"
		for _, s := range extras {
			if s == "annotations" {
				desc += fmt.Sprintf("%s: %*s\n", s, 25-len(s)+len(annotations), annotations)
			} else {
				desc += fmt.Sprintf("%s: %*s\n", s, 25-len(s)+len(strMap[s]), strMap[s])
			}
		}
		desc += "{code}" + "\n"
	}
	return desc
}

// makeSummary creates the Summary/Title of an issue.
func makeSummary(draconResult document.Document) (string, string) {
	summary := filepath.Base(draconResult.Target) + " " + draconResult.Title

	if len(summary) > 255 { // jira summary field supports up to 255 chars
		tobytes := bytes.Runes([]byte(summary))
		summary = string(tobytes[:254])
		extra := string(tobytes[255:])
		return summary, extra
	}
	return summary, ""
}
