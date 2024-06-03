// Package consumers provides helper functions for working with Dracon compatible outputs as a Consumer.
// Subdirectories in this package have more complete example usages of this package.
package consumers

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	apiv1 "github.com/ocurity/dracon/api/proto/v1"

	draconLogger "github.com/ocurity/dracon/pkg/log"
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
	// Debug flag initializes the logger with a debug level
	Debug bool
)

func init() {
	flag.StringVar(&inResults, "in", "", "the directory where dracon producer/enricher outputs are")
	flag.BoolVar(&Raw, "raw", false, "if the non-enriched results should be used")
	flag.BoolVar(&Debug, "debug", false, "turn on debug logging")

}

// ParseFlags will parse the input flags for the consumer and perform simple validation.
func ParseFlags() error {
	flag.Parse()
	if Debug {
		draconLogger.SetDefault(slog.LevelDebug, os.Getenv(EnvDraconScanID), true)
	} else {
		draconLogger.SetDefault(slog.LevelInfo, os.Getenv(EnvDraconScanID), true)
	}

	if len(inResults) < 1 {
		return fmt.Errorf("in is undefined")
	}
	return nil
}

// LoadToolResponse loads raw results from producers.
func LoadToolResponse() ([]*apiv1.LaunchToolResponse, error) {
	return putil.LoadToolResponse(inResults)
}

// LoadEnrichedToolResponse loads enriched results from the enricher.
func LoadEnrichedToolResponse() ([]*apiv1.EnrichedLaunchToolResponse, error) {
	return putil.LoadEnrichedToolResponse(inResults)
}
