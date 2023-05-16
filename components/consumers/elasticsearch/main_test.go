package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	v1 "github.com/ocurity/dracon/api/proto/v1"

	"github.com/stretchr/testify/assert"
)

var (
	want             = "OK"
	info             = `{"Version":{"Number":"8.1.0"}}`
	scanUUID         = "test-uuid"
	scanStartTime, _ = time.Parse(time.RFC3339, "2020-04-13 11:51:53+01:00")

	esIn, _ = json.Marshal(&esDocument{
		ScanStartTime:  scanStartTime,
		ScanID:         scanUUID,
		ToolName:       "es-unit-tests",
		Source:         "es-tests-source",
		Title:          "es-tests-title",
		Target:         "es-tests-target",
		Type:           "es-tests-type",
		Severity:       v1.Severity_SEVERITY_INFO,
		SeverityText:   "Info",
		CVSS:           0.01,
		Confidence:     v1.Confidence_CONFIDENCE_INFO,
		ConfidenceText: "Info",
		Description:    "es-tests-description",
		FirstFound:     scanStartTime,
		Count:          2,
		FalsePositive:  false,
		CVE:            "CVE-0000-99999",
	})
)

func TestEsPushBasicAuth(t *testing.T) {
	esIndex = "dracon-es-test"

	esStub := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		w.Header().Set("X-Elastic-Product", "Elasticsearch")
		w.WriteHeader(http.StatusOK)
		if r.Method == http.MethodGet {
			uname, pass, ok := r.BasicAuth()
			assert.Equal(t, uname, "foo")
			assert.Equal(t, pass, "bar")
			assert.Equal(t, ok, true)

			w.Write([]byte(info))
		} else if r.Method == http.MethodPost {
			// assert non authed operation (write results to index)
			assert.Equal(t, buf.String(), string(esIn))
			assert.Equal(t, r.RequestURI, "/"+esIndex+"/_doc")

			uname, pass, ok := r.BasicAuth()
			assert.Equal(t, uname, "foo")
			assert.Equal(t, pass, "bar")
			assert.Equal(t, ok, true)

			w.Write([]byte(want))
		}
	}))
	defer esStub.Close()
	os.Setenv("ELASTICSEARCH_URL", esStub.URL)

	// basic auth ops
	basicAuthUser = "foo"
	basicAuthPass = "bar"
	client, err := getESClient()
	assert.Nil(t, err)
	client.Index(esIndex, bytes.NewBuffer(esIn))
}

func TestEsPush(t *testing.T) {
	esStub := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		w.Header().Set("X-Elastic-Product", "Elasticsearch")
		w.WriteHeader(http.StatusOK)
		if r.Method == http.MethodGet {
			w.Write([]byte(info))
		} else if r.Method == http.MethodPost {
			// assert non authed operation (write results to index)
			assert.Equal(t, buf.String(), string(esIn))
			assert.Equal(t, r.RequestURI, "/"+esIndex+"/_doc")
			w.Write([]byte(want))
		}
	}))
	defer esStub.Close()
	os.Setenv("ELASTICSEARCH_URL", esStub.URL)
	client, err := getESClient()
	assert.Nil(t, err)
	client.Index(esIndex, bytes.NewBuffer(esIn))
}
