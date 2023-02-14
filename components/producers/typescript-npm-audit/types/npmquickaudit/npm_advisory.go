package npmquickaudit

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
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

// NewAdvisoryData constructs an AdvisoryData from the npm Registry advisory at
// the given URL.
func NewAdvisoryData(url string) (*AdvisoryData, error) {
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

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("npm Registry request failed: " + resp.Status)
	}

	if resp.Header.Get("Content-Type") != "application/json" {
		return nil, errors.New("npm Registry did not respond with JSON content")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
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
