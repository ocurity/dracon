package main

import (
	"encoding/json"
	"log"
	"time"

	v1 "github.com/ocurity/dracon/api/proto/v1"
	"github.com/ocurity/dracon/components/consumers"
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
	var sbom map[string]interface{}
	if iss.GetCycloneDXSBOM() != "" {
		if err := json.Unmarshal([]byte(iss.GetCycloneDXSBOM()), &sbom); err != nil {
			log.Fatalf("error unmarshaling cyclonedx sbom, err:%s", err)
		}
	}
	jBytes, err := json.Marshal(consumers.FlatenLaunchToolResponse(res))
	if err != nil {
		return []byte{}, err
	}
	return jBytes, nil
}

func getEnrichedIssue(scanStartTime time.Time, res *v1.EnrichedLaunchToolResponse, iss *v1.EnrichedIssue) ([]byte, error) {
	var sbom map[string]interface{}
	if iss.GetRawIssue().GetCycloneDXSBOM() != "" {
		if err := json.Unmarshal([]byte(iss.GetRawIssue().GetCycloneDXSBOM()), &sbom); err != nil {
			log.Fatalf("error unmarshaling cyclonedx sbom, err:%s", err)
		}
	}
	jBytes, err := json.Marshal(consumers.FlatenEnrichedLaunchToolResponse(res))
	if err != nil {
		return []byte{}, err
	}
	return jBytes, nil
}
