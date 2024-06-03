// Package enrichers provides helper functions for writing Dracon compatible enrichers that enrich dracon outputs.
package enrichers

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/go-errors/errors"

	draconapiv1 "github.com/ocurity/dracon/api/proto/v1"
	"github.com/ocurity/dracon/components"
	"github.com/ocurity/dracon/pkg/putil"
)

var (
	readPath  string
	writePath string
	debug     bool
)

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
func LoadData() ([]*draconapiv1.LaunchToolResponse, error) {
	return putil.LoadTaggedToolResponse(readPath)
}

func WriteData(enrichedIssues []*draconapiv1.EnrichedIssue, originalResults *draconapiv1.LaunchToolResponse, enricherName string) error {
	if len(enrichedIssues) > 0 {
		if err := putil.WriteEnrichedResults(originalResults, enrichedIssues,
			filepath.Join(writePath, fmt.Sprintf("%s.%s.enriched.pb", originalResults.GetToolName(), enricherName)),
		); err != nil {
			return err
		}
	} else {
		log.Println("no enriched issues were created for", originalResults.GetToolName())
	}
	if len(originalResults.GetIssues()) > 0 {
		scanStartTime := originalResults.GetScanInfo().GetScanStartTime().AsTime()
		if err := putil.WriteResults(
			originalResults.GetToolName(),
			originalResults.GetIssues(),
			filepath.Join(writePath, fmt.Sprintf("%s.raw.pb", originalResults.GetToolName())),
			originalResults.GetScanInfo().GetScanUuid(),
			scanStartTime.Format(time.RFC3339),
			originalResults.GetScanInfo().GetScanTags(),
		); err != nil {
			return errors.Errorf("could not write results: %s", err)
		}
	}
	return nil
}

func SetReadPathForTests(readFromPath string) {
	readPath = readFromPath
}

func SetWritePathForTests(writeToPath string) {
	writePath = writeToPath
}
