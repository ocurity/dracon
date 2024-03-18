// Package consumers provides helper functions for working with Dracon compatible outputs as a Consumer.
// Subdirectories in this package have more complete example usages of this package.
package consumers

import (
	"flag"
	"fmt"
	"time"

	v1 "github.com/ocurity/dracon/api/proto/v1"

	"github.com/ocurity/dracon/pkg/putil"
)

const (
	// EnvDraconStartTime Start Time of Dracon Scan in RFC3339.
	EnvDraconStartTime = "DRACON_SCAN_TIME"
	// EnvDraconScanID the ID of the dracon scan.
	EnvDraconScanID = "DRACON_SCAN_ID"
	// EnvDraconScanTags the tags of the dracon scan.
	EnvDraconScanTags = "DRACON_SCAN_TAGS"
)

var (
	inResults string
	// Raw represents if the non-enriched results should be used.
	Raw bool
)

func init() {
	flag.StringVar(&inResults, "in", "", "the directory where dracon producer/enricher outputs are")
	flag.BoolVar(&Raw, "raw", false, "if the non-enriched results should be used")
}

// ParseFlags will parse the input flags for the consumer and perform simple validation.
func ParseFlags() error {
	flag.Parse()
	if len(inResults) < 1 {
		return fmt.Errorf("in is undefined")
	}
	return nil
}

// LoadToolResponse loads raw results from producers.
func LoadToolResponse() ([]*v1.LaunchToolResponse, error) {
	return putil.LoadToolResponse(inResults)
}

// LoadEnrichedToolResponse loads enriched results from the enricher.
func LoadEnrichedToolResponse() ([]*v1.EnrichedLaunchToolResponse, error) {
	return putil.LoadEnrichedToolResponse(inResults)
}

// FlatenLaunchToolResponse returns an array of map[string]string with each element containing a flattened version of each issue, useful for writing data in datalakes
func FlatenLaunchToolResponse(response *v1.LaunchToolResponse) []map[string]string {
	result := []map[string]string{}
	for _, iss := range response.Issues {
		flat := map[string]string{
			"ScanStartTime": response.GetScanInfo().GetScanStartTime().AsTime().Format(time.RFC3339),
			"ScanID":        response.GetScanInfo().GetScanUuid(),
			"ToolName":      response.GetToolName(),
			"Source":        iss.GetSource(),
			"Title":         iss.GetTitle(),
			"Target":        iss.GetTarget(),
			"Type":          iss.GetType(),
			"Severity":      iss.GetSeverity().String(),
			"CVSS":          fmt.Sprintf("%f", iss.GetCvss()),
			"Confidence":    iss.GetConfidence().Enum().String(),
			"Description":   iss.GetDescription(),
			"CVE":           iss.GetCve(),
			"CycloneDXSBOM": iss.GetCycloneDXSBOM(),
		}
		for k, v := range response.GetScanInfo().GetScanTags() {
			flat[fmt.Sprintf("ScanTag:%s", k)] = v
		}
		result = append(result, flat)
	}
	return result
}

// FlatenLaunchToolResponse returns an array of map[string]string with each element containing a flattened version of each issue, useful for writing data in datalakes
func FlatenEnrichedLaunchToolResponse(response *v1.EnrichedLaunchToolResponse) []map[string]string {
	result := []map[string]string{}
	for _, iss := range response.GetIssues() {
		flat := map[string]string{
			"ScanStartTime": response.GetOriginalResults().GetScanInfo().GetScanStartTime().AsTime().Format(time.RFC3339),
			"ScanID":        response.GetOriginalResults().GetScanInfo().GetScanUuid(),
			"ToolName":      response.GetOriginalResults().GetToolName(),
			"Source":        iss.GetRawIssue().GetSource(),
			"Title":         iss.GetRawIssue().GetTitle(),
			"Target":        iss.GetRawIssue().GetTarget(),
			"Type":          iss.GetRawIssue().GetType(),
			"Severity":      iss.GetRawIssue().GetSeverity().String(),
			"CVSS":          fmt.Sprintf("%f", iss.GetRawIssue().GetCvss()),
			"Confidence":    iss.GetRawIssue().GetConfidence().Enum().String(),
			"Description":   iss.GetRawIssue().GetDescription(),
			"CVE":           iss.GetRawIssue().GetCve(),
			"CycloneDXSBOM": iss.GetRawIssue().GetCycloneDXSBOM(),
			"FirstSeen":     iss.GetFirstSeen().AsTime().Format(time.RFC3339),
			"Count":         fmt.Sprintf("%d", iss.GetCount()),
			"FalsePositive": fmt.Sprintf("%t", iss.GetFalsePositive()),
			"UpdatedAt":     iss.GetUpdatedAt().AsTime().Format(time.RFC3339),
		}
		for k, v := range response.GetOriginalResults().GetScanInfo().GetScanTags() {
			flat[fmt.Sprintf("ScanTag:%s", k)] = v
		}
		for k, v := range iss.GetAnnotations() {
			flat[fmt.Sprintf("Annotation:%s", k)] = v
		}

		result = append(result, flat)
	}
	return result
}
