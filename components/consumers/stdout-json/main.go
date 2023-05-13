package main

import (
	"encoding/json"
	"log"
	"time"

	v1 "github.com/ocurity/dracon/api/proto/v1"
	"github.com/ocurity/dracon/components/consumers"
	"github.com/ocurity/dracon/pkg/enumtransformers"
)

func main() {
	if err := consumers.ParseFlags(); err != nil {
		log.Fatal(err)
	}
	if consumers.Raw {
		responses, err := consumers.LoadToolResponse()
		if err != nil {
			log.Fatal("could not load raw results, file malformed: ", err)
		}
		for _, res := range responses {
			scanStartTime := res.GetScanInfo().GetScanStartTime().AsTime()
			for _, iss := range res.GetIssues() {
				b, err := getRawIssue(scanStartTime, res, iss)
				if err != nil {
					log.Fatal("Could not parse raw issue", err)
				}
				log.Printf("%s", string(b))
			}
		}
	} else {
		responses, err := consumers.LoadEnrichedToolResponse()
		if err != nil {
			log.Fatal("could not load enriched results, file malformed: ", err)
		}
		for _, res := range responses {
			scanStartTime := res.GetOriginalResults().GetScanInfo().GetScanStartTime().AsTime()
			for _, iss := range res.GetIssues() {
				b, err := getEnrichedIssue(scanStartTime, res, iss)
				if err != nil {
					log.Fatal("Could not parse enriched issue", err)
				}
				log.Printf("%s", string(b))
			}
		}
	}
}

func getRawIssue(scanStartTime time.Time, res *v1.LaunchToolResponse, iss *v1.Issue) ([]byte, error) {
	jBytes, err := json.Marshal(&draconDocument{
		ScanStartTime: scanStartTime,
		ScanID:        res.GetScanInfo().GetScanUuid(),
		ToolName:      res.GetToolName(),
		Source:        iss.GetSource(),
		Title:         iss.GetTitle(),
		Target:        iss.GetTarget(),
		Type:          iss.GetType(),
		Severity:      iss.GetSeverity(),
		CVSS:          iss.GetCvss(),
		Confidence:    iss.GetConfidence(),
		Description:   iss.GetDescription(),
		FirstFound:    scanStartTime,
		Count:         1,
		FalsePositive: false,
		CVE:           iss.GetCve(),
	})
	if err != nil {
		return []byte{}, err
	}
	return jBytes, nil
}

func getEnrichedIssue(scanStartTime time.Time, res *v1.EnrichedLaunchToolResponse, iss *v1.EnrichedIssue) ([]byte, error) {
	firstSeenTime := iss.GetFirstSeen().AsTime()
	jBytes, err := json.Marshal(&draconDocument{
		ScanStartTime:  scanStartTime,
		ScanID:         res.GetOriginalResults().GetScanInfo().GetScanUuid(),
		ToolName:       res.GetOriginalResults().GetToolName(),
		Source:         iss.GetRawIssue().GetSource(),
		Title:          iss.GetRawIssue().GetTitle(),
		Target:         iss.GetRawIssue().GetTarget(),
		Type:           iss.GetRawIssue().GetType(),
		Severity:       iss.GetRawIssue().GetSeverity(),
		CVSS:           iss.GetRawIssue().GetCvss(),
		Confidence:     iss.GetRawIssue().GetConfidence(),
		Description:    iss.GetRawIssue().GetDescription(),
		FirstFound:     firstSeenTime,
		Count:          iss.GetCount(),
		FalsePositive:  iss.GetFalsePositive(),
		SeverityText:   enumtransformers.SeverityToText(iss.GetRawIssue().GetSeverity()),
		ConfidenceText: enumtransformers.ConfidenceToText(iss.GetRawIssue().GetConfidence()),
		CVE:            iss.GetRawIssue().GetCve(),
	})
	if err != nil {
		return []byte{}, err
	}
	return jBytes, nil
}

type draconDocument struct {
	ScanStartTime  time.Time     `json:"scan_start_time"`
	ScanID         string        `json:"scan_id"`
	ToolName       string        `json:"tool_name"`
	Source         string        `json:"source"`
	Target         string        `json:"target"`
	Type           string        `json:"type"`
	Title          string        `json:"title"`
	Severity       v1.Severity   `json:"severity"`
	SeverityText   string        `json:"severity_text"`
	CVSS           float64       `json:"cvss"`
	Confidence     v1.Confidence `json:"confidence"`
	ConfidenceText string        `json:"confidence_text"`
	Description    string        `json:"description"`
	FirstFound     time.Time     `json:"first_found"`
	Count          uint64        `json:"count"`
	FalsePositive  bool          `json:"false_positive"`
	CVE            string        `json:"cve"`
}
