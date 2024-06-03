// Package enrichers provides helper functions for writing Dracon compatible enrichers that enrich dracon outputs.
package enrichers

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/go-errors/errors"

	draconV1 "github.com/ocurity/dracon/api/proto/v1"
	"github.com/ocurity/dracon/components"
	"github.com/ocurity/dracon/pkg/putil"
)

var (
	readPath  string
	writePath string
	debug     bool
)

// LookupEnvOrString will return the value of the environment variable
// if it exists, otherwise it will return the default value.
func LookupEnvOrString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

// ParseFlags will parse the input flags for the producer and perform simple validation.
func ParseFlags() error {
	flag.StringVar(&readPath, "read_path", LookupEnvOrString("READ_PATH", ""), "where to find producer results")
	flag.StringVar(&writePath, "write_path", LookupEnvOrString("WRITE_PATH", ""), "where to put enriched results")
	flag.BoolVar(&debug, "debug", false, "turn on debug logging")

	flag.Parse()
	logLevel := slog.LevelInfo
	if debug {
		logLevel = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})).With("scanID", os.Getenv(components.EnvDraconScanID)))
	if readPath == "" {
		return fmt.Errorf("read_path is undefined")
	}
	if writePath == "" {
		return fmt.Errorf("write_path is undefined")
	}
	return nil
}

// LoadData returns the LaunchToolResponses meant for this enricher.
func LoadData() ([]*draconV1.LaunchToolResponse, error) {
	return putil.LoadTaggedToolResponse(readPath)
}

// WriteData will write the enriched results to the write path.
func WriteData(enrichedLaunchToolResponse *draconV1.EnrichedLaunchToolResponse, enricherName string) error {
	if enrichedLaunchToolResponse == nil || len(enrichedLaunchToolResponse.Issues) == 0 {
		return errors.Errorf("no enriched issues were created for %s", enrichedLaunchToolResponse.GetOriginalResults().GetToolName())
	}
	if err := putil.WriteEnrichedResults(enrichedLaunchToolResponse.GetOriginalResults(), enrichedLaunchToolResponse.GetIssues(),
		filepath.Join(writePath, fmt.Sprintf("%s.%s.enriched.pb", enrichedLaunchToolResponse.GetOriginalResults().GetToolName(), enricherName)),
	); err != nil {
		return err
	}
	if enrichedLaunchToolResponse.OriginalResults == nil || len(enrichedLaunchToolResponse.OriginalResults.GetIssues()) == 0 {
		return errors.Errorf("original results is empty for %s", enrichedLaunchToolResponse.GetOriginalResults().GetToolName())
	}
	scanStartTime := enrichedLaunchToolResponse.GetOriginalResults().GetScanInfo().GetScanStartTime().AsTime()
	if err := putil.WriteResults(
		enrichedLaunchToolResponse.GetOriginalResults().GetToolName(),
		enrichedLaunchToolResponse.GetOriginalResults().GetIssues(),
		filepath.Join(writePath, fmt.Sprintf("%s.raw.pb", enrichedLaunchToolResponse.GetOriginalResults().GetToolName())),
		enrichedLaunchToolResponse.GetOriginalResults().GetScanInfo().GetScanUuid(),
		scanStartTime,
		enrichedLaunchToolResponse.GetOriginalResults().GetScanInfo().GetScanTags(),
	); err != nil {
		return errors.Errorf("could not write results: %s", err)
	}
	return nil
}

// SetReadPathForTests sets the read path for tests.
func SetReadPathForTests(readFromPath string) {
	readPath = readFromPath
}

// SetWritePathForTests sets the write path for tests.
func SetWritePathForTests(writeToPath string) {
	writePath = writeToPath
}
