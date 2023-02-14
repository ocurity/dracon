package consumers

import (
	"log"

	"github.com/golang/protobuf/ptypes"
)

func ExampleParseFlags() {
	if err := ParseFlags(); err != nil {
		log.Fatal(err)
	}
}

func ExampleLoadToolResponse() {
	responses, err := LoadToolResponse()
	if err != nil {
		log.Fatal(err)
	}
	for _, res := range responses {
		scanStartTime, _ := ptypes.Timestamp(res.GetScanInfo().GetScanStartTime())
		_ = scanStartTime
		for _, iss := range res.GetIssues() {
			// Do your own logic with issues here
			_ = iss
		}
	}
}

func ExampleLoadEnrichedToolResponse() {
	responses, err := LoadEnrichedToolResponse()
	if err != nil {
		log.Fatal(err)
	}
	for _, res := range responses {
		scanStartTime, _ := ptypes.Timestamp(res.GetOriginalResults().GetScanInfo().GetScanStartTime())
		_ = scanStartTime
		for _, iss := range res.GetIssues() {
			// Do your own logic with issues here
			_ = iss
		}
	}
}
