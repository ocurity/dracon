// Package consumers provides helper functions for working with Dracon compatible outputs as a Consumer.
// Subdirectories in this package have more complete example usages of this package.
package consumers

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	draconapiv1 "github.com/ocurity/dracon/api/proto/v1"
	"github.com/ocurity/dracon/components"
	"github.com/ocurity/dracon/pkg/putil"
)

var (
	inResults string
	// Raw represents if the non-enriched results should be used.
	Raw bool
	// debug flag initializes the logger with a debug level
	debug bool
)

func init() {
	flag.StringVar(&inResults, "in", "", "the directory where dracon producer/enricher outputs are")
	flag.BoolVar(&Raw, "raw", false, "if the non-enriched results should be used")
	flag.BoolVar(&debug, "debug", false, "turn on debug logging")
}

// ParseFlags will parse the input flags for the consumer and perform simple validation.
func ParseFlags() error {
	flag.Parse()

	logLevel := slog.LevelInfo
	if debug {
		logLevel = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})).With("scanID", os.Getenv(components.EnvDraconScanID)))
	if len(inResults) < 1 {
		return fmt.Errorf("in is undefined")
	}
	return nil
}

// LoadToolResponse loads raw results from producers.
func LoadToolResponse() ([]*draconapiv1.LaunchToolResponse, error) {
	return putil.LoadToolResponse(inResults)
}

// LoadEnrichedToolResponse loads enriched results from the enricher.
func LoadEnrichedToolResponse() ([]*draconapiv1.EnrichedLaunchToolResponse, error) {
	return putil.LoadEnrichedToolResponse(inResults)
}
