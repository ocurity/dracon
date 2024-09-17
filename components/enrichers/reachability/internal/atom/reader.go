package atom

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/jmespath/go-jmespath"

	"github.com/ocurity/dracon/components/enrichers/reachability/internal/atom/purl"
	"github.com/ocurity/dracon/components/enrichers/reachability/internal/logging"
)

type (
	// Reader is responsible for managing how to read atom files and understand their contents.
	Reader struct {
		atomFilePath string
		purlParser   *purl.Parser
	}

	// ReachablePurls maps reachable purls based on the report.
	ReachablePurls map[string]struct{}

	// Reachables is a slice of Reachable.
	Reachables []Reachable

	// Response maps the content of an atom reachability report in json format.
	Response struct {
		Reachables Reachables `json:"reachables"`
	}

	// Reachable represents an atom reachable result.
	Reachable struct {
		Flows []Flows  `json:"flows"`
		Purls []string `json:"purls"`
	}

	// Flows describes the flows on how to reach such reachable.
	Flows struct {
		ID                    int    `json:"id"`
		Label                 string `json:"label"`
		Name                  string `json:"name"`
		FullName              string `json:"fullName"`
		Signature             string `json:"signature"`
		IsExternal            bool   `json:"isExternal"`
		Code                  string `json:"code"`
		TypeFullName          string `json:"typeFullName"`
		ParentMethodName      string `json:"parentMethodName"`
		ParentMethodSignature string `json:"parentMethodSignature"`
		ParentFileName        string `json:"parentFileName"`
		ParentPackageName     string `json:"parentPackageName"`
		ParentClassName       string `json:"parentClassName"`
		LineNumber            int    `json:"lineNumber"`
		ColumnNumber          int    `json:"columnNumber"`
		Tags                  string `json:"tags"`
	}
)

// NewReader returns a new atom file reader.
func NewReader(atomFilePath string, purlParser *purl.Parser) (*Reader, error) {
	switch {
	case atomFilePath == "":
		return nil, errors.New("invalid empty atom file path")
	case purlParser == nil:
		return nil, errors.New("invalid nil purl parser")
	}

	return &Reader{
		atomFilePath: atomFilePath,
		purlParser:   purlParser,
	}, nil
}

// Read deserialises the json content of the provided atom file into Reachables format.
func (r *Reader) Read(ctx context.Context) (*Response, error) {
	b, err := os.ReadFile(r.atomFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read atom file: %w", err)
	}

	logging.FromContext(ctx).Debug("sample atom file contents", slog.String("payload", string(b)))

	var res Response
	if err := json.Unmarshal(b, &res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal atom response: %w", err)
	}

	return &res, nil
}

// ReachablePurls finds all the reachable purls presents in the atom reachability result.
func (r *Reader) ReachablePurls(ctx context.Context, reachables *Response) (ReachablePurls, error) {
	logger := logging.FromContext(ctx)

	rawPurls, err := jmespath.Search("reachables[].purls[]", reachables)
	if err != nil {
		return nil, fmt.Errorf("failed to search reachable purls: %w", err)
	}

	purls, ok := rawPurls.([]any)
	if !ok {
		logger.Error(
			"invalid raw reachable purl. Expected an array",
			slog.Any("raw_purls", rawPurls),
		)
		return nil, errors.New("invalid raw reachable purl. Expected an array")
	}

	uniquePurls := make(map[string]struct{})
	for idx, p := range purls {
		ps, ok := p.(string)
		if !ok {
			logger.Error(
				"unexpected purl type, expected a string. Continuing...",
				slog.Any("purl", p),
				slog.Int("index", idx),
			)
			continue
		}
		uniquePurls[ps] = struct{}{}
	}

	finalPurls := make(ReachablePurls)
	for p := range uniquePurls {
		parsedPurls, err := r.purlParser.ParsePurl(p)
		if err != nil {
			logger.Error(
				"could not parse purl. Continuing...",
				slog.Any("purl", p),
			)
			continue
		}

		for _, pp := range parsedPurls {
			finalPurls[pp] = struct{}{}
		}
	}

	return finalPurls, nil
}

func (rp ReachablePurls) IsPurlReachable(purl string) bool {
	purl = strings.ToLower(purl)
	_, ok := rp[purl]
	return ok
}
