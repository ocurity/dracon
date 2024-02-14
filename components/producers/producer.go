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
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	v1 "github.com/ocurity/dracon/api/proto/v1"

	"github.com/ocurity/dracon/pkg/putil"
)

var (
	// InResults represents incoming tool output.
	InResults string
	// OutFile points to the protobuf file where dracon results will be written.
	OutFile string
	// Append flag will append to the outfile instead of overwriting, useful when there's multiple inresults.
	Append bool
)

// migrate from flags to viper
//  use viper to autoenv
// if tags env var exists, automap to viper scanTags
// make sure no producers touch scan results and if they do, overwrite with tags

const (
	sourceDir = "/workspace/source-code-ws"

	// EnvDraconStartTime Start Time of Dracon Scan in RFC3339.
	EnvDraconStartTime = "DRACON_SCAN_TIME"
	// EnvDraconScanID the ID of the dracon scan.
	EnvDraconScanID = "DRACON_SCAN_ID"
	// EnvDraconScanTags the tags of the dracon scan.
	EnvDraconScanTags = "DRACON_SCAN_TAGS"
)

// ParseFlags will parse the input flags for the producer and perform simple validation.
func ParseFlags() error {
	flag.StringVar(&InResults, "in", "", "")
	flag.StringVar(&OutFile, "out", "", "")
	flag.BoolVar(&Append, "append", false, "Append to output file instead of overwriting it")

	if InResults == "" {
		return fmt.Errorf("in is undefined")
	}
	if OutFile == "" {
		return fmt.Errorf("out is undefined")
	}
	return nil
}

// ReadLines returns the lines of the contents of the file given by InResults.
func ReadLines() ([][]byte, error) {
	file, err := os.Open(InResults)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var result [][]byte

	for scanner.Scan() {
		result = append(result, scanner.Bytes())
	}

	return result, nil
}

// ReadInFile returns the contents of the file given by InResults.
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

// ParseJSON provides a generic method to parse JSON input (e.g. the results
// provided by a tool) into a given struct.
func ParseJSON(in []byte, structure interface{}) error {
	if err := json.Unmarshal(in, &structure); err != nil {
		return err
	}
	return nil
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
	issues []*v1.Issue,
) error {
	source := getSource()
	cleanIssues := []*v1.Issue{}
	for _, iss := range issues {
		iss.Description = strings.ReplaceAll(iss.Description, sourceDir, ".")
		iss.Title = strings.ReplaceAll(iss.Title, sourceDir, ".")
		iss.Target = strings.ReplaceAll(iss.Target, sourceDir, ".")
		iss.Source = source
		cleanIssues = append(cleanIssues, iss)
		log.Printf("found issue: %+v\n", iss)
	}
	scanStartTime := strings.TrimSpace(os.Getenv(EnvDraconStartTime))
	if scanStartTime == "" {
		scanStartTime = time.Now().UTC().Format(time.RFC3339)
	}
	scanUUUID := strings.TrimSpace(os.Getenv(EnvDraconScanID))
	scanTagsStr := strings.TrimSpace(os.Getenv(EnvDraconScanTags))
	scanTags := map[string]string{}
	err := json.Unmarshal([]byte(scanTagsStr), &scanTags)
	if err != nil {
		log.Println("scan with uuid", scanUUUID, "does not have any tags, err: '", err, "'")
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

	dat, err := ioutil.ReadFile(sourceMetaPath)
	if err != nil {
		log.Println(err)
	}
	return strings.TrimSpace(string(dat))
}
