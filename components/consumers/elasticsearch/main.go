package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"strings"
	"time"

	v1 "github.com/ocurity/dracon/api/proto/v1"
	"github.com/ocurity/dracon/components/consumers"
	"github.com/ocurity/dracon/pkg/enumtransformers"
	"github.com/ocurity/dracon/pkg/templating"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
)

var (
	esUrls  string
	esAddrs []string
	esIndex string

	esAPIKey  string
	esCloudID string

	basicAuthUser string
	basicAuthPass string
	issueTemplate string
)

func parseFlags() error {
	if err := consumers.ParseFlags(); err != nil {
		return err
	}
	if len(esIndex) == 0 {
		return fmt.Errorf("esIndex '%s' is undefined", esIndex)
	}
	if len(esUrls) > 0 {
		for _, u := range strings.Split(esUrls, ",") {
			esAddrs = append(esAddrs, strings.TrimSpace(u))
		}
	}
	return nil
}

func main() {
	flag.StringVar(&esIndex, "esIndex", "", "the index in elasticsearch to push results to")
	flag.StringVar(&issueTemplate, "descriptionTemplate", "", "a Go Template string describing how to show Raw or Enriched issues")

	// es SaaS options
	flag.StringVar(&esAPIKey, "esAPIKey", "", "the api key in elasticsearch to contact results to")
	flag.StringVar(&esCloudID, "esCloudID", "", "the cloud id in elasticsearch to contact results to")

	// es self-hosted options
	flag.StringVar(&esUrls, "esURL", "", "[OPTIONAL] URLs to connect to elasticsearch comma separated. Can also use env variable ELASTICSEARCH_URL")
	flag.StringVar(&basicAuthUser, "basic-auth-user", "", "[OPTIONAL] the basic auth username")
	flag.StringVar(&basicAuthPass, "basic-auth-pass", "", "[OPTIONAL] the basic auth password")
	flag.Parse()

	if err := parseFlags(); err != nil {
		log.Fatal(err)
	}

	es, err := getESClient()
	if err != nil {
		log.Fatal("could not contact remote Elasticsearch: ", err)
	}

	if consumers.Raw {
		slog.Debug("Parsing Raw results")
		responses, err := consumers.LoadToolResponse()
		if err != nil {
			log.Fatal("could not load raw results, file malformed: ", err)
		}
		numIssues := 0
		for _, res := range responses {
			scanStartTime := res.GetScanInfo().GetScanStartTime().AsTime()
			numIssues += len(res.GetIssues())
			for _, iss := range res.GetIssues() {
				b, err := getRawIssue(scanStartTime, res, iss)
				if err != nil {
					log.Fatal("Could not parse raw issue", err)
				}
				res, err := es.Index(esIndex, bytes.NewBuffer(b))
				log.Printf("%+v", res)
				if err != nil {
					log.Fatal("Could not push raw issue", err)
				}
			}
		}
		slog.Info("Pushed", "numIssues", numIssues, "issues to Elasticsearch", "")
	} else {
		log.Print("Parsing Enriched results")
		responses, err := consumers.LoadEnrichedToolResponse()
		if err != nil {
			log.Fatal("could not load enriched results, file malformed: ", err)
		}
		numIssues := 0
		for _, res := range responses {
			scanStartTime := res.GetOriginalResults().GetScanInfo().GetScanStartTime().AsTime()
			numIssues += len(res.GetIssues())
			for _, iss := range res.GetIssues() {
				b, err := getEnrichedIssue(scanStartTime, res, iss)
				if err != nil {
					log.Fatal("Could not parse enriched issue", err)
				}
				res, err := es.Index(esIndex, bytes.NewBuffer(b))
				if err != nil || res.IsError() {
					log.Fatal("Could not push enriched issue", err, "received", res.StatusCode)
				}
			}
		}
		slog.Info("Pushed", "numIssues", numIssues, "issues to Elasticsearch", "")
	}
}

func getRawIssue(scanStartTime time.Time, res *v1.LaunchToolResponse, iss *v1.Issue) ([]byte, error) {
	description, err := templating.TemplateStringRaw(issueTemplate, iss)
	if err != nil {
		log.Fatal("Could not template raw issue ", err)
	}
	jBytes, err := json.Marshal(&esDocument{
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
		Description:   *description,
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
	description, err := templating.TemplateStringEnriched(issueTemplate, iss)
	if err != nil {
		log.Fatal("Could not template raw issue ", err)
	}
	firstSeenTime := iss.GetFirstSeen().AsTime()
	jBytes, err := json.Marshal(&esDocument{
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
		Description:    *description,
		FirstFound:     firstSeenTime,
		Count:          iss.GetCount(),
		FalsePositive:  iss.GetFalsePositive(),
		SeverityText:   enumtransformers.SeverityToText(iss.GetRawIssue().GetSeverity()),
		ConfidenceText: enumtransformers.ConfidenceToText(iss.GetRawIssue().GetConfidence()),
		CVE:            iss.GetRawIssue().GetCve(),
		Annotations:    iss.GetAnnotations(),
	})
	if err != nil {
		return []byte{}, err
	}
	return jBytes, nil
}

type esDocument struct {
	ScanStartTime  time.Time         `json:"scan_start_time"`
	ScanID         string            `json:"scan_id"`
	ToolName       string            `json:"tool_name"`
	Source         string            `json:"source"`
	Target         string            `json:"target"`
	Type           string            `json:"type"`
	Title          string            `json:"title"`
	Severity       v1.Severity       `json:"severity"`
	SeverityText   string            `json:"severity_text"`
	CVSS           float64           `json:"cvss"`
	Confidence     v1.Confidence     `json:"confidence"`
	ConfidenceText string            `json:"confidence_text"`
	Description    string            `json:"description"`
	FirstFound     time.Time         `json:"first_found"`
	Count          uint64            `json:"count"`
	FalsePositive  bool              `json:"false_positive"`
	CVE            string            `json:"cve"`
	Annotations    map[string]string `json:"annotations"`
}

func getESClient() (*elasticsearch.Client, error) {
	var es *elasticsearch.Client
	var err error
	var esConfig elasticsearch.Config = elasticsearch.Config{}

	if basicAuthUser != "" && basicAuthPass != "" {
		esConfig.Username = basicAuthUser
		esConfig.Password = basicAuthPass
	}
	if esAPIKey != "" {
		esConfig.APIKey = esAPIKey
	}
	if esCloudID != "" {
		esConfig.CloudID = esCloudID
	}
	if len(esAddrs) > 0 {
		esConfig.Addresses = esAddrs
	}

	es, err = elasticsearch.NewClient(esConfig)
	if err != nil {
		return nil, err
	}
	type esInfo struct {
		Version struct {
			Number string `json:"number"`
		} `json:"version"`
	}

	res, err := es.Info()
	if err != nil {
		return nil, err
	}
	var info esInfo
	if err := json.NewDecoder(res.Body).Decode(&info); err != nil {
		return nil, err
	}
	if info.Version.Number[0] != '8' {
		return nil, fmt.Errorf("unsupported ES Server version %s", info.Version.Number)
	}
	return es, err
}
