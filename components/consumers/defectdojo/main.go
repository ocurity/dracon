package main

import (
	"encoding/json"
	"flag"
	"log"
	"strconv"
	"time"

	v1 "github.com/ocurity/dracon/api/proto/v1"
	"github.com/ocurity/dracon/components/consumers"
	"github.com/ocurity/dracon/components/consumers/defectdojo/client"
)

// DojoTimeFormat is the time format accepted by defect dojo.
const DojoTimeFormat = "2006-01-02"

// DojoTestTimeFormat is the time format expected by defect dojo when creating a test.
const DojoTestTimeFormat = "2006-01-02T03:03"

var (
	authUser               string
	authToken              string
	authURL                string
	newEngagementEveryScan bool
	productID              string
)

func init() {
	// envUser := os.Getenv(EnvDojoUser)
	// envToken := os.Getenv(EnvDojoToken)
	// envURL := os.Getenv(EnvDojoURL)

	flag.StringVar(&authUser, "dojoUser", "", "defect dojo user")
	flag.StringVar(&authToken, "dojoToken", "", "defect dojo api token")
	flag.StringVar(&authURL, "dojoURL", "", "defect dojo api base url")
	flag.StringVar(&productID, "dojoProductID", "", "defect dojo product ID if you want to create an engagement")
	flag.BoolVar(&newEngagementEveryScan, "createEngagement", false, "for every dracon scan id, create a different engagement")
	flag.Parse()
}

func handleRawResults(product int, dojoClient *client.Client, responses []*v1.LaunchToolResponse) error {
	if len(responses) == 0 {
		log.Println("No tool responses provided")
		return nil
	}
	scanUUID := responses[0].GetScanInfo().GetScanUuid()
	if scanUUID == "" {
		log.Fatalln("Non-uuid scan", responses)
	}
	tags := []string{"DraconScan", "RawScan", scanUUID}

	engagement, err := dojoClient.CreateEngagement( // with current architecture, all responses should have the same scaninfo
		scanUUID, responses[0].GetScanInfo().GetScanStartTime().AsTime().Format(DojoTimeFormat), tags, int32(product))
	if err != nil {
		return err
	}
	for _, res := range responses {
		log.Println("handling response for tool", res.GetToolName(), "with", len(res.GetIssues()), "findings")
		startTime := res.GetScanInfo().GetScanStartTime().AsTime()
		test, err := dojoClient.CreateTest(startTime.Format(DojoTestTimeFormat), res.GetToolName(), "", []string{"DraconScan", "RawTest", scanUUID}, engagement.ID)
		if err != nil {
			log.Printf("could not create test in defectdojo, err: %#v", err)
			return err
		}
		for _, iss := range res.GetIssues() {
			b, err := getRawIssue(res.GetScanInfo().GetScanStartTime().AsTime(), res, iss)
			if err != nil {
				log.Fatal("Could not parse raw issue", err)
			}
			description := string(b)
			finding, err := dojoClient.CreateFinding(
				iss.GetTitle(),
				description,
				severityToText(iss.GetSeverity()),
				iss.GetTarget(),
				startTime.Format(DojoTimeFormat),
				severityToDojoSeverity(iss.Severity),
				[]string{"DraconScan", "RawFinding", scanUUID, res.GetToolName()},
				test.ID,
				0,
				0,
				dojoClient.UserID,
				false,
				false,
				true,
				iss.GetCvss())
			if err != nil {
				log.Fatalf("Could not create raw finding error: %v\n", err)
			} else {
				log.Println("Created finding successfully", finding)
			}
		}
	}
	return nil
}

func handleEnrichedResults(product int, dojoClient *client.Client, responses []*v1.EnrichedLaunchToolResponse) error {
	if len(responses) == 0 {
		log.Println("No tool responses provided")
		return nil
	}
	scanUUID := responses[0].GetOriginalResults().GetScanInfo().GetScanUuid()
	if scanUUID == "" {
		log.Fatalln("Non-uuid scan", responses)
	}
	tags := []string{"DraconScan", "EnrichedScan", scanUUID}
	engagement, err := dojoClient.CreateEngagement( // with current architecture, all responses should have the same scaninfo
		scanUUID,
		responses[0].GetOriginalResults().GetScanInfo().GetScanStartTime().AsTime().Format(DojoTimeFormat), tags, int32(product))
	if err != nil {
		log.Println("could not create Engagement, err:", err)
		return err
	}
	for _, res := range responses {
		log.Println("handling response for tool", res.GetOriginalResults().GetToolName(), "with", len(res.GetIssues()), "findings")

		scanStartTime := res.GetOriginalResults().GetScanInfo().GetScanStartTime().AsTime()
		test, err := dojoClient.CreateTest(scanStartTime.Format(DojoTestTimeFormat), res.GetOriginalResults().GetToolName(), "", []string{"DraconScan", "EnrichedTest", scanUUID}, engagement.ID)
		if err != nil {
			log.Println("could not create test in defectdojo, err:", err)
			return err
		}
		for _, iss := range res.GetIssues() {
			b, err := getEnrichedIssue(scanStartTime, res, iss)
			if err != nil {
				log.Fatal("Could not parse enriched issue", err)
			}
			description := string(b)
			duplicate := false
			if iss.GetFirstSeen().AsTime().Before(scanStartTime) || iss.GetCount() > 1 {
				duplicate = true
			}
			rawIss := iss.GetRawIssue()
			finding, err := dojoClient.CreateFinding(
				rawIss.GetTitle(),
				description,
				severityToText(rawIss.GetSeverity()),
				rawIss.GetTarget(),
				scanStartTime.Format(DojoTimeFormat),
				severityToDojoSeverity(rawIss.Severity),
				[]string{"DraconScan", "EnrichedFinding", scanUUID, res.GetOriginalResults().GetToolName()},
				test.ID, 0, 0, dojoClient.UserID,
				iss.GetFalsePositive(),
				duplicate,
				true,
				rawIss.GetCvss())
			if err != nil {
				log.Fatalf("Could not create enriched finding error: %v\n", err)
			} else {
				log.Printf("Created enriched finding successfully %v\n", finding)
			}
		}
	}
	return nil
}

func main() {
	if err := consumers.ParseFlags(); err != nil {
		log.Fatal(err)
	}
	product, err := strconv.Atoi(productID)
	if err != nil {
		log.Fatalf("productID %s is not a number, err: %#v\n", productID, err)
	}
	client, err := client.DojoClient(authURL, authToken, authUser)
	if err != nil {
		log.Panicf("could not instantiate dojo client err: %#v\n", err)
	}
	if consumers.Raw {
		responses, err := consumers.LoadToolResponse()
		if err != nil {
			log.Fatal("could not load raw results, file malformed: ", err)
		}
		if err = handleRawResults(product, client, responses); err != nil {
			log.Fatalf("Could not handle raw results, err %v", err)
		}
	} else {
		responses, err := consumers.LoadEnrichedToolResponse()
		if err != nil {
			log.Fatal("could not load enriched results, file malformed: ", err)
		}
		if err = handleEnrichedResults(product, client, responses); err != nil {
			log.Fatalf("Could not handle enriched results, err %v", err)
		}

	}
}

func getRawIssue(scanStartTime time.Time, res *v1.LaunchToolResponse, iss *v1.Issue) ([]byte, error) {
	jBytes, err := json.Marshal(&draconDocument{
		Confidence:     iss.GetConfidence(),
		ConfidenceText: confidenceToText(iss.GetConfidence()),
		Count:          1,
		CVE:            iss.GetCve(),
		CVSS:           iss.GetCvss(),
		Description:    iss.GetDescription(),
		FalsePositive:  false,
		FirstFound:     scanStartTime.Format(DojoTimeFormat),
		ScanID:         res.GetScanInfo().GetScanUuid(),
		ScanStartTime:  scanStartTime.Format(DojoTimeFormat),
		Severity:       iss.GetSeverity(),
		SeverityText:   severityToText(iss.GetSeverity()),
		Source:         iss.GetSource(),
		Target:         iss.GetTarget(),
		Title:          iss.GetTitle(),
		ToolName:       res.GetToolName(),
		Type:           iss.GetType(),
	})
	if err != nil {
		return []byte{}, err
	}
	return jBytes, nil
}

func severityToDojoSeverity(severity v1.Severity) string {
	switch severity {
	case v1.Severity_SEVERITY_INFO:
		return "S:I"
	case v1.Severity_SEVERITY_LOW:
		return "S:L"
	case v1.Severity_SEVERITY_MEDIUM:
		return "S:M"
	case v1.Severity_SEVERITY_HIGH:
		return "S:H"
	case v1.Severity_SEVERITY_CRITICAL:
		return "S:C"
	default:
		return "S:I"
	}
}

func severityToText(severity v1.Severity) string {
	switch severity {
	case v1.Severity_SEVERITY_INFO:
		return "Info"
	case v1.Severity_SEVERITY_LOW:
		return "Low"
	case v1.Severity_SEVERITY_MEDIUM:
		return "Medium"
	case v1.Severity_SEVERITY_HIGH:
		return "High"
	case v1.Severity_SEVERITY_CRITICAL:
		return "Critical"
	default:
		return "N/A"
	}
}

func confidenceToText(confidence v1.Confidence) string {
	switch confidence {
	case v1.Confidence_CONFIDENCE_INFO:
		return "Info"
	case v1.Confidence_CONFIDENCE_LOW:
		return "Low"
	case v1.Confidence_CONFIDENCE_MEDIUM:
		return "Medium"
	case v1.Confidence_CONFIDENCE_HIGH:
		return "High"
	case v1.Confidence_CONFIDENCE_CRITICAL:
		return "Critical"
	default:
		return "N/A"
	}
}

func getEnrichedIssue(scanStartTime time.Time, res *v1.EnrichedLaunchToolResponse, iss *v1.EnrichedIssue) ([]byte, error) {
	firstSeenTime := iss.GetFirstSeen().AsTime()
	jBytes, err := json.Marshal(&draconDocument{
		Confidence:     iss.GetRawIssue().GetConfidence(),
		ConfidenceText: confidenceToText(iss.GetRawIssue().GetConfidence()),
		Count:          iss.GetCount(),
		CVE:            iss.GetRawIssue().GetCve(),
		CVSS:           iss.GetRawIssue().GetCvss(),
		Description:    iss.GetRawIssue().GetDescription(),
		FalsePositive:  iss.GetFalsePositive(),
		FirstFound:     firstSeenTime.Format(DojoTimeFormat),
		ScanID:         res.GetOriginalResults().GetScanInfo().GetScanUuid(),
		ScanStartTime:  scanStartTime.Format(DojoTimeFormat),
		Severity:       iss.GetRawIssue().GetSeverity(),
		SeverityText:   severityToText(iss.GetRawIssue().GetSeverity()),
		Source:         iss.GetRawIssue().GetSource(),
		Target:         iss.GetRawIssue().GetTarget(),
		Title:          iss.GetRawIssue().GetTitle(),
		ToolName:       res.GetOriginalResults().GetToolName(),
		Type:           iss.GetRawIssue().GetType(),
	})
	if err != nil {
		return []byte{}, err
	}
	return jBytes, nil
}

type draconDocument struct {
	ScanStartTime  string        `json:"scan_start_time"`
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
	FirstFound     string        `json:"first_found"`
	Count          uint64        `json:"count"`
	FalsePositive  bool          `json:"false_positive"`
	CVE            string        `json:"cve"`
}
