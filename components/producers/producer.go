// Package producers provides helper functions for writing Dracon compatible producers that parse tool outputs.
// Subdirectories in this package have more complete example usages of this package.
package producers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	draconapiv1 "github.com/ocurity/dracon/api/proto/v1"
	"github.com/ocurity/dracon/components"

	"github.com/ocurity/dracon/pkg/putil"
)

var (
	// InResults represents incoming tool output.
	InResults string
	// OutFile points to the protobuf file where dracon results will be written.
	OutFile string
	// Append flag will append to the outfile instead of overwriting, useful when there's multiple inresults.
	Append bool
	// debug flag initializes the logger with a debug level
	debug bool
)

const (
	sourceDir = "/workspace/output"
)

// ParseFlags will parse the input flags for the producer and perform simple validation.
func ParseFlags() error {
	flag.StringVar(&InResults, "in", "", "")
	flag.StringVar(&OutFile, "out", "", "")
	flag.BoolVar(&debug, "debug", false, "turn on debug logging")
	flag.BoolVar(&Append, "append", false, "Append to output file instead of overwriting it")

	flag.Parse()
	logLevel := slog.LevelInfo
	if debug {
		logLevel = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})).With("scanID", os.Getenv(components.EnvDraconScanID)))

	if InResults == "" {
		return fmt.Errorf("in is undefined")
	}
	if OutFile == "" {
		return fmt.Errorf("out is undefined")
	}
	return nil
}

// ReadLines returns the lines of the contents of the file given by InResults.
func ReadLines() (result [][]byte, err error) {
	var file *os.File

	file, err = os.Open(InResults)
	if err != nil {
		return nil, err
	}
	defer func() { err = errors.Join(err, file.Close()) }()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		result = append(result, scanner.Bytes())
	}

	return result, err
}

// ReadInFile returns the contents of the file given by InResults.
// TODO: replace with os.ReadFile
func ReadInFile() ([]byte, error) {
	file, err := os.Open(InResults)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buffer := new(bytes.Buffer)
	if _, err := buffer.ReadFrom(file); err != nil {
		return nil, fmt.Errorf("could not read from buffer: %w", err)
	}
	return buffer.Bytes(), nil
}

// ParseMultiJSONMessages provides method to parse tool results in JSON format.
// It allows for parsing single JSON files with multiple JSON messages in them.
func ParseMultiJSONMessages(in []byte) ([]interface{}, error) {
	dec := json.NewDecoder(strings.NewReader(string(in)))
	var res []interface{}
	for {
		var item interface{}
		err := dec.Decode(&item)
		if errors.Is(err, io.EOF) {
			res = append(res, item)
			break
		} else if err != nil {
			return res, err
		}
		res = append(res, item)
	}
	return res, nil
}

// WriteDraconOut provides a generic method to write the resulting protobuf to the output file.
func WriteDraconOut(
	toolName string,
	issues []*draconapiv1.Issue,
) error {
	source := getSource()
	cleanIssues := []*draconapiv1.Issue{}
	for _, iss := range issues {
		iss.Description = strings.ReplaceAll(iss.Description, sourceDir, ".")
		iss.Title = strings.ReplaceAll(iss.Title, sourceDir, ".")
		iss.Target = strings.ReplaceAll(iss.Target, sourceDir, ".")
		iss.Source = source
		cleanIssues = append(cleanIssues, iss)
		slog.Info(fmt.Sprintf("found issue: %+v\n", iss))
	}
	scanStartTime := strings.TrimSpace(os.Getenv(components.EnvDraconStartTime))
	if scanStartTime == "" {
		scanStartTime = time.Now().UTC().Format(time.RFC3339)
	}
	scanUUUID := strings.TrimSpace(os.Getenv(components.EnvDraconScanID))
	scanTagsStr := strings.TrimSpace(os.Getenv(components.EnvDraconScanTags))
	scanTags := map[string]string{}
	err := json.Unmarshal([]byte(scanTagsStr), &scanTags)
	if err != nil {
		slog.Error(fmt.Sprintf("scan does not have any tags, err:%s", err))
	}

	stat, err := os.Stat(OutFile)
	if Append && err == nil && stat.Size() > 0 {
		return putil.AppendResults(cleanIssues, OutFile)
	}
	return putil.WriteResults(toolName, cleanIssues, OutFile, scanUUUID, scanStartTime, scanTags)
}

func getSource() string {
	sourceMetaPath := filepath.Join(sourceDir, ".source.dracon")
	_, err := os.Stat(sourceMetaPath)
	if os.IsNotExist(err) {
		return "unknown"
	}

	dat, err := os.ReadFile(sourceMetaPath)
	if err != nil {
		slog.Error(err.Error())
	}
	return strings.TrimSpace(string(dat))
}
