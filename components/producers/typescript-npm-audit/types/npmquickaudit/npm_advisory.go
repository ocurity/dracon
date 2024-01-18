package npmquickaudit

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// AdvisoryData represents a subset of the data returned in an advisoryData
// object in an npm Registry advisory. Only the data relevant to Dracon issue
// reports is retained.
type AdvisoryData struct {
	CVEs               []string `json:"cves"`
	CWE                string   `json:"cwe"`
	Overview           string   `json:"overview"`
	PatchedVersions    string   `json:"patched_versions"`
	Recommendation     string   `json:"recommendation"`
	References         string   `json:"recommendations"`
	VulnerableVersions string   `json:"vulnerable_versions"`
}

// FetchAdvisoryData constructs an AdvisoryData from the npm registry advisory
// at the given URL.
// If the advisory redirects to Github, then the code will automatically fetch
// the GH advisory.
func FetchAdvisoryData(url string) (*AdvisoryData, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// "X-Spiferack: 1" as a request header provides a JSON-encoded
	// response from an endpoint that usually responds with HTML:
	// https://npm.community/t/can-i-query-npm-for-all-advisory-information/2096/5
	req.Header.Add("X-Spiferack", "1")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, errors.New("npm Registry request failed: " + resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json") {
		return nil, errors.New("npm Registry did not respond with JSON content")
	}

	if resp.StatusCode > 300 {
		if resp.Header.Get("Location") == "" {
			return nil, fmt.Errorf("%s: no location header in redirect response", url)
		}
		return FetchAdvisoryData(resp.Header.Get("Location"))
	} else if resp.StatusCode == 200 && strings.Contains(string(body), "Redirect") && resp.Header.Get("Location") != "" {
		// under normal circumstances we should not get a response like this but the  NPM API currently has a bug
		return FetchAdvisoryData(resp.Header.Get("Location"))
	}

	var data map[string]json.RawMessage
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	if _, ok := data["advisoryData"]; !ok {
		return nil, errors.New("npm Registry response did not contain an advisoryData key")
	}

	var advisory *AdvisoryData
	if err := json.Unmarshal(data["advisoryData"], &advisory); err != nil {
		return nil, err
	}

	return advisory, nil
}
