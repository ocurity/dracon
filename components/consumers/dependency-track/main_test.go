package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	dtrack "github.com/DependencyTrack/client-go"
	"github.com/google/uuid"
	v1 "github.com/ocurity/dracon/api/proto/v1"
	cyclonedx "github.com/ocurity/dracon/pkg/cyclonedx"
	"github.com/stretchr/testify/assert"
)

func TestUploadBomsFromRaw(t *testing.T) {
	projUUID := uuid.MustParse("7c78f6c9-b4b0-493c-a912-0bb0a4f221f1")
	expectedRequest := dtrack.BOMUploadRequest{
		ProjectName:    "test",
		ProjectUUID:    &projUUID,
		ProjectVersion: "2022-1",
		AutoCreate:     false,
		BOM:            "eyJib21Gb3JtYXQiOiJDeWNsb25lRFgiLCJzcGVjVmVyc2lvbiI6IjEuNCIsInNlcmlhbE51bWJlciI6InVybjp1dWlkOjNlNjcxNjg3LTM5NWItNDFmNS1hMzBmLWE1ODkyMWE2OWI3OSIsInZlcnNpb24iOjEsIm1ldGFkYXRhIjp7InRpbWVzdGFtcCI6IjIwMjEtMDEtMTBUMTI6MDA6MDBaIiwiY29tcG9uZW50Ijp7ImJvbS1yZWYiOiJhY21lLWFwcGxpY2F0aW9uIiwidHlwZSI6ImFwcGxpY2F0aW9uIiwibmFtZSI6IkFjbWUgQ2xvdWQgRXhhbXBsZSIsInZlcnNpb24iOiIyMDIyLTEifX0sInNlcnZpY2VzIjpbeyJib20tcmVmIjoiYXBpLWdhdGV3YXkiLCJwcm92aWRlciI6eyJuYW1lIjoiQWNtZSBJbmMiLCJ1cmwiOlsiaHR0cHM6Ly9leGFtcGxlLmNvbSJdfSwiZ3JvdXAiOiJjb20uZXhhbXBsZSIsIm5hbWUiOiJBUEkgR2F0ZXdheSIsInZlcnNpb24iOiIyMDIyLTEiLCJkZXNjcmlwdGlvbiI6IkV4YW1wbGUgQVBJIEdhdGV3YXkiLCJlbmRwb2ludHMiOlsiaHR0cHM6Ly9leGFtcGxlLmNvbS8iLCJodHRwczovL2V4YW1wbGUuY29tL2FwcCJdLCJhdXRoZW50aWNhdGVkIjpmYWxzZSwieC10cnVzdC1ib3VuZGFyeSI6dHJ1ZSwiZGF0YSI6W3siZmxvdyI6ImJpLWRpcmVjdGlvbmFsIiwiY2xhc3NpZmljYXRpb24iOiJQSUkifSx7ImZsb3ciOiJiaS1kaXJlY3Rpb25hbCIsImNsYXNzaWZpY2F0aW9uIjoiUElGSSJ9LHsiZmxvdyI6ImJpLWRpcmVjdGlvbmFsIiwiY2xhc3NpZmljYXRpb24iOiJQdWJsaWMifV0sImV4dGVybmFsUmVmZXJlbmNlcyI6W3sidXJsIjoiaHR0cDovL2V4YW1wbGUuY29tL2FwcC9zd2FnZ2VyIiwidHlwZSI6ImRvY3VtZW50YXRpb24ifV0sInNlcnZpY2VzIjpbeyJib20tcmVmIjoibXMtMS5leGFtcGxlLmNvbSIsInByb3ZpZGVyIjp7Im5hbWUiOiJBY21lIEluYyIsInVybCI6WyJodHRwczovL2V4YW1wbGUuY29tIl19LCJncm91cCI6ImNvbS5leGFtcGxlIiwibmFtZSI6Ik1pY3Jvc2VydmljZSAxIiwidmVyc2lvbiI6IjIwMjItMSIsImRlc2NyaXB0aW9uIjoiRXhhbXBsZSBNaWNyb3NlcnZpY2UiLCJlbmRwb2ludHMiOlsiaHR0cHM6Ly9tcy0xLmV4YW1wbGUuY29tIl0sImF1dGhlbnRpY2F0ZWQiOnRydWUsIngtdHJ1c3QtYm91bmRhcnkiOmZhbHNlLCJkYXRhIjpbeyJmbG93IjoiYmktZGlyZWN0aW9uYWwiLCJjbGFzc2lmaWNhdGlvbiI6IlBJSSJ9XSwiZXh0ZXJuYWxSZWZlcmVuY2VzIjpbeyJ1cmwiOiJodHRwczovL21zLTEuZXhhbXBsZS5jb20vc3dhZ2dlciIsInR5cGUiOiJkb2N1bWVudGF0aW9uIn1dfSx7ImJvbS1yZWYiOiJtcy0yLmV4YW1wbGUuY29tIiwicHJvdmlkZXIiOnsibmFtZSI6IkFjbWUgSW5jIiwidXJsIjpbImh0dHBzOi8vZXhhbXBsZS5jb20iXX0sImdyb3VwIjoiY29tLmV4YW1wbGUiLCJuYW1lIjoiTWljcm9zZXJ2aWNlIDIiLCJ2ZXJzaW9uIjoiMjAyMi0xIiwiZGVzY3JpcHRpb24iOiJFeGFtcGxlIE1pY3Jvc2VydmljZSIsImVuZHBvaW50cyI6WyJodHRwczovL21zLTIuZXhhbXBsZS5jb20iXSwiYXV0aGVudGljYXRlZCI6dHJ1ZSwieC10cnVzdC1ib3VuZGFyeSI6ZmFsc2UsImRhdGEiOlt7ImZsb3ciOiJiaS1kaXJlY3Rpb25hbCIsImNsYXNzaWZpY2F0aW9uIjoiUElGSSJ9XSwiZXh0ZXJuYWxSZWZlcmVuY2VzIjpbeyJ1cmwiOiJodHRwczovL21zLTIuZXhhbXBsZS5jb20vc3dhZ2dlciIsInR5cGUiOiJkb2N1bWVudGF0aW9uIn1dfSx7ImJvbS1yZWYiOiJtcy0zLmV4YW1wbGUuY29tIiwicHJvdmlkZXIiOnsibmFtZSI6IkFjbWUgSW5jIiwidXJsIjpbImh0dHBzOi8vZXhhbXBsZS5jb20iXX0sImdyb3VwIjoiY29tLmV4YW1wbGUiLCJuYW1lIjoiTWljcm9zZXJ2aWNlIDMiLCJ2ZXJzaW9uIjoiMjAyMi0xIiwiZGVzY3JpcHRpb24iOiJFeGFtcGxlIE1pY3Jvc2VydmljZSIsImVuZHBvaW50cyI6WyJodHRwczovL21zLTMuZXhhbXBsZS5jb20iXSwiYXV0aGVudGljYXRlZCI6dHJ1ZSwieC10cnVzdC1ib3VuZGFyeSI6ZmFsc2UsImRhdGEiOlt7ImZsb3ciOiJiaS1kaXJlY3Rpb25hbCIsImNsYXNzaWZpY2F0aW9uIjoiUHVibGljIn1dLCJleHRlcm5hbFJlZmVyZW5jZXMiOlt7InVybCI6Imh0dHBzOi8vbXMtMy5leGFtcGxlLmNvbS9zd2FnZ2VyIiwidHlwZSI6ImRvY3VtZW50YXRpb24ifV19LHsiYm9tLXJlZiI6Im1zLTEtcGdzcWwuZXhhbXBsZS5jb20iLCJncm91cCI6Im9yZy5wb3N0Z3Jlc3FsIiwibmFtZSI6IlBvc3RncmVzIiwidmVyc2lvbiI6IjE0LjEiLCJkZXNjcmlwdGlvbiI6IlBvc3RncmVzIGRhdGFiYXNlIGZvciBNaWNyb3NlcnZpY2UgIzEiLCJlbmRwb2ludHMiOlsiaHR0cHM6Ly9tcy0xLXBnc3FsLmV4YW1wbGUuY29tOjU0MzIiXSwiYXV0aGVudGljYXRlZCI6dHJ1ZSwieC10cnVzdC1ib3VuZGFyeSI6ZmFsc2UsImRhdGEiOlt7ImZsb3ciOiJiaS1kaXJlY3Rpb25hbCIsImNsYXNzaWZpY2F0aW9uIjoiUElJIn1dfSx7ImJvbS1yZWYiOiJzMy1leGFtcGxlLmFtYXpvbi5jb20iLCJncm91cCI6ImNvbS5hbWF6b24iLCJuYW1lIjoiUzMiLCJkZXNjcmlwdGlvbiI6IlMzIGJ1Y2tldCIsImVuZHBvaW50cyI6WyJodHRwczovL3MzLWV4YW1wbGUuYW1hem9uLmNvbSJdLCJhdXRoZW50aWNhdGVkIjp0cnVlLCJ4LXRydXN0LWJvdW5kYXJ5Ijp0cnVlLCJkYXRhIjpbeyJmbG93IjoiYmktZGlyZWN0aW9uYWwiLCJjbGFzc2lmaWNhdGlvbiI6IlB1YmxpYyJ9XX1dfV0sImRlcGVuZGVuY2llcyI6W3sicmVmIjoiYWNtZS1hcHBsaWNhdGlvbiIsImRlcGVuZHNPbiI6WyJhcGktZ2F0ZXdheSJdfSx7InJlZiI6ImFwaS1nYXRld2F5IiwiZGVwZW5kc09uIjpbIm1zLTEuZXhhbXBsZS5jb20iLCJtcy0yLmV4YW1wbGUuY29tIiwibXMtMy5leGFtcGxlLmNvbSJdfSx7InJlZiI6Im1zLTEuZXhhbXBsZS5jb20iLCJkZXBlbmRzT24iOlsibXMtMS1wZ3NxbC5leGFtcGxlLmNvbSJdfSx7InJlZiI6Im1zLTIuZXhhbXBsZS5jb20iLCJkZXBlbmRzT24iOltdfSx7InJlZiI6Im1zLTMuZXhhbXBsZS5jb20iLCJkZXBlbmRzT24iOlsiczMtZXhhbXBsZS5hbWF6b24uY29tIl19XX0=",
	}

	//nolint:gosec
	expectedToken := "7c78f6c9-token"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		var actualRequest dtrack.BOMUploadRequest
		json.Unmarshal(body, &actualRequest)
		assert.Equal(t, expectedRequest, actualRequest)
		w.Write([]byte("{\"Token\":\"" + expectedToken + "\"}"))
	}))
	defer ts.Close()
	projectUUID = projUUID.String()
	apiKey = "test"
	projectName = "test"
	c, err := dtrack.NewClient(ts.URL, dtrack.WithAPIKey(apiKey))
	assert.Nil(t, err)

	client = c
	issues, err := cyclonedx.ToDracon([]byte(saasBOM), "json")

	assert.Nil(t, err)
	ltr := v1.LaunchToolResponse{
		ToolName: "SAT",
		Issues:   issues,
	}
	tokens, err := uploadBOMsFromRaw([]*v1.LaunchToolResponse{&ltr})
	assert.Nil(t, err)
	assert.Equal(t, tokens, []string{expectedToken})
}

func TestUploadBomsFromEnriched(t *testing.T) {
	projUUID := uuid.MustParse("7c78f6c9-b4b0-493c-a912-0bb0a4f221f1")
	expectedRequest := dtrack.BOMUploadRequest{
		ProjectName:    "test",
		ProjectUUID:    &projUUID,
		ProjectVersion: "2022-1",
		AutoCreate:     false,
		BOM:            "eyJib21Gb3JtYXQiOiJDeWNsb25lRFgiLCJzcGVjVmVyc2lvbiI6IjEuNCIsInNlcmlhbE51bWJlciI6InVybjp1dWlkOjNlNjcxNjg3LTM5NWItNDFmNS1hMzBmLWE1ODkyMWE2OWI3OSIsInZlcnNpb24iOjEsIm1ldGFkYXRhIjp7InRpbWVzdGFtcCI6IjIwMjEtMDEtMTBUMTI6MDA6MDBaIiwiY29tcG9uZW50Ijp7ImJvbS1yZWYiOiJhY21lLWFwcGxpY2F0aW9uIiwidHlwZSI6ImFwcGxpY2F0aW9uIiwibmFtZSI6IkFjbWUgQ2xvdWQgRXhhbXBsZSIsInZlcnNpb24iOiIyMDIyLTEifX0sInNlcnZpY2VzIjpbeyJib20tcmVmIjoiYXBpLWdhdGV3YXkiLCJwcm92aWRlciI6eyJuYW1lIjoiQWNtZSBJbmMiLCJ1cmwiOlsiaHR0cHM6Ly9leGFtcGxlLmNvbSJdfSwiZ3JvdXAiOiJjb20uZXhhbXBsZSIsIm5hbWUiOiJBUEkgR2F0ZXdheSIsInZlcnNpb24iOiIyMDIyLTEiLCJkZXNjcmlwdGlvbiI6IkV4YW1wbGUgQVBJIEdhdGV3YXkiLCJlbmRwb2ludHMiOlsiaHR0cHM6Ly9leGFtcGxlLmNvbS8iLCJodHRwczovL2V4YW1wbGUuY29tL2FwcCJdLCJhdXRoZW50aWNhdGVkIjpmYWxzZSwieC10cnVzdC1ib3VuZGFyeSI6dHJ1ZSwiZGF0YSI6W3siZmxvdyI6ImJpLWRpcmVjdGlvbmFsIiwiY2xhc3NpZmljYXRpb24iOiJQSUkifSx7ImZsb3ciOiJiaS1kaXJlY3Rpb25hbCIsImNsYXNzaWZpY2F0aW9uIjoiUElGSSJ9LHsiZmxvdyI6ImJpLWRpcmVjdGlvbmFsIiwiY2xhc3NpZmljYXRpb24iOiJQdWJsaWMifV0sImV4dGVybmFsUmVmZXJlbmNlcyI6W3sidXJsIjoiaHR0cDovL2V4YW1wbGUuY29tL2FwcC9zd2FnZ2VyIiwidHlwZSI6ImRvY3VtZW50YXRpb24ifV0sInNlcnZpY2VzIjpbeyJib20tcmVmIjoibXMtMS5leGFtcGxlLmNvbSIsInByb3ZpZGVyIjp7Im5hbWUiOiJBY21lIEluYyIsInVybCI6WyJodHRwczovL2V4YW1wbGUuY29tIl19LCJncm91cCI6ImNvbS5leGFtcGxlIiwibmFtZSI6Ik1pY3Jvc2VydmljZSAxIiwidmVyc2lvbiI6IjIwMjItMSIsImRlc2NyaXB0aW9uIjoiRXhhbXBsZSBNaWNyb3NlcnZpY2UiLCJlbmRwb2ludHMiOlsiaHR0cHM6Ly9tcy0xLmV4YW1wbGUuY29tIl0sImF1dGhlbnRpY2F0ZWQiOnRydWUsIngtdHJ1c3QtYm91bmRhcnkiOmZhbHNlLCJkYXRhIjpbeyJmbG93IjoiYmktZGlyZWN0aW9uYWwiLCJjbGFzc2lmaWNhdGlvbiI6IlBJSSJ9XSwiZXh0ZXJuYWxSZWZlcmVuY2VzIjpbeyJ1cmwiOiJodHRwczovL21zLTEuZXhhbXBsZS5jb20vc3dhZ2dlciIsInR5cGUiOiJkb2N1bWVudGF0aW9uIn1dfSx7ImJvbS1yZWYiOiJtcy0yLmV4YW1wbGUuY29tIiwicHJvdmlkZXIiOnsibmFtZSI6IkFjbWUgSW5jIiwidXJsIjpbImh0dHBzOi8vZXhhbXBsZS5jb20iXX0sImdyb3VwIjoiY29tLmV4YW1wbGUiLCJuYW1lIjoiTWljcm9zZXJ2aWNlIDIiLCJ2ZXJzaW9uIjoiMjAyMi0xIiwiZGVzY3JpcHRpb24iOiJFeGFtcGxlIE1pY3Jvc2VydmljZSIsImVuZHBvaW50cyI6WyJodHRwczovL21zLTIuZXhhbXBsZS5jb20iXSwiYXV0aGVudGljYXRlZCI6dHJ1ZSwieC10cnVzdC1ib3VuZGFyeSI6ZmFsc2UsImRhdGEiOlt7ImZsb3ciOiJiaS1kaXJlY3Rpb25hbCIsImNsYXNzaWZpY2F0aW9uIjoiUElGSSJ9XSwiZXh0ZXJuYWxSZWZlcmVuY2VzIjpbeyJ1cmwiOiJodHRwczovL21zLTIuZXhhbXBsZS5jb20vc3dhZ2dlciIsInR5cGUiOiJkb2N1bWVudGF0aW9uIn1dfSx7ImJvbS1yZWYiOiJtcy0zLmV4YW1wbGUuY29tIiwicHJvdmlkZXIiOnsibmFtZSI6IkFjbWUgSW5jIiwidXJsIjpbImh0dHBzOi8vZXhhbXBsZS5jb20iXX0sImdyb3VwIjoiY29tLmV4YW1wbGUiLCJuYW1lIjoiTWljcm9zZXJ2aWNlIDMiLCJ2ZXJzaW9uIjoiMjAyMi0xIiwiZGVzY3JpcHRpb24iOiJFeGFtcGxlIE1pY3Jvc2VydmljZSIsImVuZHBvaW50cyI6WyJodHRwczovL21zLTMuZXhhbXBsZS5jb20iXSwiYXV0aGVudGljYXRlZCI6dHJ1ZSwieC10cnVzdC1ib3VuZGFyeSI6ZmFsc2UsImRhdGEiOlt7ImZsb3ciOiJiaS1kaXJlY3Rpb25hbCIsImNsYXNzaWZpY2F0aW9uIjoiUHVibGljIn1dLCJleHRlcm5hbFJlZmVyZW5jZXMiOlt7InVybCI6Imh0dHBzOi8vbXMtMy5leGFtcGxlLmNvbS9zd2FnZ2VyIiwidHlwZSI6ImRvY3VtZW50YXRpb24ifV19LHsiYm9tLXJlZiI6Im1zLTEtcGdzcWwuZXhhbXBsZS5jb20iLCJncm91cCI6Im9yZy5wb3N0Z3Jlc3FsIiwibmFtZSI6IlBvc3RncmVzIiwidmVyc2lvbiI6IjE0LjEiLCJkZXNjcmlwdGlvbiI6IlBvc3RncmVzIGRhdGFiYXNlIGZvciBNaWNyb3NlcnZpY2UgIzEiLCJlbmRwb2ludHMiOlsiaHR0cHM6Ly9tcy0xLXBnc3FsLmV4YW1wbGUuY29tOjU0MzIiXSwiYXV0aGVudGljYXRlZCI6dHJ1ZSwieC10cnVzdC1ib3VuZGFyeSI6ZmFsc2UsImRhdGEiOlt7ImZsb3ciOiJiaS1kaXJlY3Rpb25hbCIsImNsYXNzaWZpY2F0aW9uIjoiUElJIn1dfSx7ImJvbS1yZWYiOiJzMy1leGFtcGxlLmFtYXpvbi5jb20iLCJncm91cCI6ImNvbS5hbWF6b24iLCJuYW1lIjoiUzMiLCJkZXNjcmlwdGlvbiI6IlMzIGJ1Y2tldCIsImVuZHBvaW50cyI6WyJodHRwczovL3MzLWV4YW1wbGUuYW1hem9uLmNvbSJdLCJhdXRoZW50aWNhdGVkIjp0cnVlLCJ4LXRydXN0LWJvdW5kYXJ5Ijp0cnVlLCJkYXRhIjpbeyJmbG93IjoiYmktZGlyZWN0aW9uYWwiLCJjbGFzc2lmaWNhdGlvbiI6IlB1YmxpYyJ9XX1dfV0sImRlcGVuZGVuY2llcyI6W3sicmVmIjoiYWNtZS1hcHBsaWNhdGlvbiIsImRlcGVuZHNPbiI6WyJhcGktZ2F0ZXdheSJdfSx7InJlZiI6ImFwaS1nYXRld2F5IiwiZGVwZW5kc09uIjpbIm1zLTEuZXhhbXBsZS5jb20iLCJtcy0yLmV4YW1wbGUuY29tIiwibXMtMy5leGFtcGxlLmNvbSJdfSx7InJlZiI6Im1zLTEuZXhhbXBsZS5jb20iLCJkZXBlbmRzT24iOlsibXMtMS1wZ3NxbC5leGFtcGxlLmNvbSJdfSx7InJlZiI6Im1zLTIuZXhhbXBsZS5jb20iLCJkZXBlbmRzT24iOltdfSx7InJlZiI6Im1zLTMuZXhhbXBsZS5jb20iLCJkZXBlbmRzT24iOlsiczMtZXhhbXBsZS5hbWF6b24uY29tIl19XX0=",
	}
	expectedToken := "7c78f6c9-token"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		var actualRequest dtrack.BOMUploadRequest
		json.Unmarshal(body, &actualRequest)
		assert.Equal(t, expectedRequest, actualRequest)
		w.Write([]byte("{\"Token\":\"" + expectedToken + "\"}"))
	}))
	defer ts.Close()

	projectUUID = projUUID.String()
	apiKey = "test"
	projectName = "test"
	c, err := dtrack.NewClient(ts.URL, dtrack.WithAPIKey(apiKey))
	assert.Nil(t, err)

	client = c
	issues, err := cyclonedx.ToDracon([]byte(saasBOM), "json")

	assert.Nil(t, err)
	ltr := v1.LaunchToolResponse{
		ToolName: "SAT",
		Issues:   issues,
	}
	eltr := v1.EnrichedLaunchToolResponse{
		OriginalResults: &ltr,
		Issues: []*v1.EnrichedIssue{
			{RawIssue: issues[0]},
		},
	}
	tokens, err := uploadBOMSFromEnriched([]*v1.EnrichedLaunchToolResponse{&eltr})
	assert.Nil(t, err)
	assert.Equal(t, tokens, []string{expectedToken})
}



func TestUploadBomsFromEnrichedWithOwners(t *testing.T) {
	projUUID := uuid.MustParse("7c78f6c9-b4b0-493c-a912-0bb0a4f221f1")
	expectedRequest := dtrack.BOMUploadRequest{
		ProjectName:    "test",
		ProjectUUID:    &projUUID,
		ProjectVersion: "2022-1",
		AutoCreate:     false,
		BOM:            "eyJib21Gb3JtYXQiOiJDeWNsb25lRFgiLCJzcGVjVmVyc2lvbiI6IjEuNCIsInNlcmlhbE51bWJlciI6InVybjp1dWlkOjNlNjcxNjg3LTM5NWItNDFmNS1hMzBmLWE1ODkyMWE2OWI3OSIsInZlcnNpb24iOjEsIm1ldGFkYXRhIjp7InRpbWVzdGFtcCI6IjIwMjEtMDEtMTBUMTI6MDA6MDBaIiwiY29tcG9uZW50Ijp7ImJvbS1yZWYiOiJhY21lLWFwcGxpY2F0aW9uIiwidHlwZSI6ImFwcGxpY2F0aW9uIiwibmFtZSI6IkFjbWUgQ2xvdWQgRXhhbXBsZSIsInZlcnNpb24iOiIyMDIyLTEifX0sInNlcnZpY2VzIjpbeyJib20tcmVmIjoiYXBpLWdhdGV3YXkiLCJwcm92aWRlciI6eyJuYW1lIjoiQWNtZSBJbmMiLCJ1cmwiOlsiaHR0cHM6Ly9leGFtcGxlLmNvbSJdfSwiZ3JvdXAiOiJjb20uZXhhbXBsZSIsIm5hbWUiOiJBUEkgR2F0ZXdheSIsInZlcnNpb24iOiIyMDIyLTEiLCJkZXNjcmlwdGlvbiI6IkV4YW1wbGUgQVBJIEdhdGV3YXkiLCJlbmRwb2ludHMiOlsiaHR0cHM6Ly9leGFtcGxlLmNvbS8iLCJodHRwczovL2V4YW1wbGUuY29tL2FwcCJdLCJhdXRoZW50aWNhdGVkIjpmYWxzZSwieC10cnVzdC1ib3VuZGFyeSI6dHJ1ZSwiZGF0YSI6W3siZmxvdyI6ImJpLWRpcmVjdGlvbmFsIiwiY2xhc3NpZmljYXRpb24iOiJQSUkifSx7ImZsb3ciOiJiaS1kaXJlY3Rpb25hbCIsImNsYXNzaWZpY2F0aW9uIjoiUElGSSJ9LHsiZmxvdyI6ImJpLWRpcmVjdGlvbmFsIiwiY2xhc3NpZmljYXRpb24iOiJQdWJsaWMifV0sImV4dGVybmFsUmVmZXJlbmNlcyI6W3sidXJsIjoiaHR0cDovL2V4YW1wbGUuY29tL2FwcC9zd2FnZ2VyIiwidHlwZSI6ImRvY3VtZW50YXRpb24ifV0sInNlcnZpY2VzIjpbeyJib20tcmVmIjoibXMtMS5leGFtcGxlLmNvbSIsInByb3ZpZGVyIjp7Im5hbWUiOiJBY21lIEluYyIsInVybCI6WyJodHRwczovL2V4YW1wbGUuY29tIl19LCJncm91cCI6ImNvbS5leGFtcGxlIiwibmFtZSI6Ik1pY3Jvc2VydmljZSAxIiwidmVyc2lvbiI6IjIwMjItMSIsImRlc2NyaXB0aW9uIjoiRXhhbXBsZSBNaWNyb3NlcnZpY2UiLCJlbmRwb2ludHMiOlsiaHR0cHM6Ly9tcy0xLmV4YW1wbGUuY29tIl0sImF1dGhlbnRpY2F0ZWQiOnRydWUsIngtdHJ1c3QtYm91bmRhcnkiOmZhbHNlLCJkYXRhIjpbeyJmbG93IjoiYmktZGlyZWN0aW9uYWwiLCJjbGFzc2lmaWNhdGlvbiI6IlBJSSJ9XSwiZXh0ZXJuYWxSZWZlcmVuY2VzIjpbeyJ1cmwiOiJodHRwczovL21zLTEuZXhhbXBsZS5jb20vc3dhZ2dlciIsInR5cGUiOiJkb2N1bWVudGF0aW9uIn1dfSx7ImJvbS1yZWYiOiJtcy0yLmV4YW1wbGUuY29tIiwicHJvdmlkZXIiOnsibmFtZSI6IkFjbWUgSW5jIiwidXJsIjpbImh0dHBzOi8vZXhhbXBsZS5jb20iXX0sImdyb3VwIjoiY29tLmV4YW1wbGUiLCJuYW1lIjoiTWljcm9zZXJ2aWNlIDIiLCJ2ZXJzaW9uIjoiMjAyMi0xIiwiZGVzY3JpcHRpb24iOiJFeGFtcGxlIE1pY3Jvc2VydmljZSIsImVuZHBvaW50cyI6WyJodHRwczovL21zLTIuZXhhbXBsZS5jb20iXSwiYXV0aGVudGljYXRlZCI6dHJ1ZSwieC10cnVzdC1ib3VuZGFyeSI6ZmFsc2UsImRhdGEiOlt7ImZsb3ciOiJiaS1kaXJlY3Rpb25hbCIsImNsYXNzaWZpY2F0aW9uIjoiUElGSSJ9XSwiZXh0ZXJuYWxSZWZlcmVuY2VzIjpbeyJ1cmwiOiJodHRwczovL21zLTIuZXhhbXBsZS5jb20vc3dhZ2dlciIsInR5cGUiOiJkb2N1bWVudGF0aW9uIn1dfSx7ImJvbS1yZWYiOiJtcy0zLmV4YW1wbGUuY29tIiwicHJvdmlkZXIiOnsibmFtZSI6IkFjbWUgSW5jIiwidXJsIjpbImh0dHBzOi8vZXhhbXBsZS5jb20iXX0sImdyb3VwIjoiY29tLmV4YW1wbGUiLCJuYW1lIjoiTWljcm9zZXJ2aWNlIDMiLCJ2ZXJzaW9uIjoiMjAyMi0xIiwiZGVzY3JpcHRpb24iOiJFeGFtcGxlIE1pY3Jvc2VydmljZSIsImVuZHBvaW50cyI6WyJodHRwczovL21zLTMuZXhhbXBsZS5jb20iXSwiYXV0aGVudGljYXRlZCI6dHJ1ZSwieC10cnVzdC1ib3VuZGFyeSI6ZmFsc2UsImRhdGEiOlt7ImZsb3ciOiJiaS1kaXJlY3Rpb25hbCIsImNsYXNzaWZpY2F0aW9uIjoiUHVibGljIn1dLCJleHRlcm5hbFJlZmVyZW5jZXMiOlt7InVybCI6Imh0dHBzOi8vbXMtMy5leGFtcGxlLmNvbS9zd2FnZ2VyIiwidHlwZSI6ImRvY3VtZW50YXRpb24ifV19LHsiYm9tLXJlZiI6Im1zLTEtcGdzcWwuZXhhbXBsZS5jb20iLCJncm91cCI6Im9yZy5wb3N0Z3Jlc3FsIiwibmFtZSI6IlBvc3RncmVzIiwidmVyc2lvbiI6IjE0LjEiLCJkZXNjcmlwdGlvbiI6IlBvc3RncmVzIGRhdGFiYXNlIGZvciBNaWNyb3NlcnZpY2UgIzEiLCJlbmRwb2ludHMiOlsiaHR0cHM6Ly9tcy0xLXBnc3FsLmV4YW1wbGUuY29tOjU0MzIiXSwiYXV0aGVudGljYXRlZCI6dHJ1ZSwieC10cnVzdC1ib3VuZGFyeSI6ZmFsc2UsImRhdGEiOlt7ImZsb3ciOiJiaS1kaXJlY3Rpb25hbCIsImNsYXNzaWZpY2F0aW9uIjoiUElJIn1dfSx7ImJvbS1yZWYiOiJzMy1leGFtcGxlLmFtYXpvbi5jb20iLCJncm91cCI6ImNvbS5hbWF6b24iLCJuYW1lIjoiUzMiLCJkZXNjcmlwdGlvbiI6IlMzIGJ1Y2tldCIsImVuZHBvaW50cyI6WyJodHRwczovL3MzLWV4YW1wbGUuYW1hem9uLmNvbSJdLCJhdXRoZW50aWNhdGVkIjp0cnVlLCJ4LXRydXN0LWJvdW5kYXJ5Ijp0cnVlLCJkYXRhIjpbeyJmbG93IjoiYmktZGlyZWN0aW9uYWwiLCJjbGFzc2lmaWNhdGlvbiI6IlB1YmxpYyJ9XX1dfV0sImRlcGVuZGVuY2llcyI6W3sicmVmIjoiYWNtZS1hcHBsaWNhdGlvbiIsImRlcGVuZHNPbiI6WyJhcGktZ2F0ZXdheSJdfSx7InJlZiI6ImFwaS1nYXRld2F5IiwiZGVwZW5kc09uIjpbIm1zLTEuZXhhbXBsZS5jb20iLCJtcy0yLmV4YW1wbGUuY29tIiwibXMtMy5leGFtcGxlLmNvbSJdfSx7InJlZiI6Im1zLTEuZXhhbXBsZS5jb20iLCJkZXBlbmRzT24iOlsibXMtMS1wZ3NxbC5leGFtcGxlLmNvbSJdfSx7InJlZiI6Im1zLTIuZXhhbXBsZS5jb20iLCJkZXBlbmRzT24iOltdfSx7InJlZiI6Im1zLTMuZXhhbXBsZS5jb20iLCJkZXBlbmRzT24iOlsiczMtZXhhbXBsZS5hbWF6b24uY29tIl19XX0=",
	}
	expectedProjectUpdate := dtrack.Project{
		UUID:       projUUID,
		Name:       "fooProj",
		PURL:       "pkg://npm/xyz/asdf@v1.2.2",
		Properties: []dtrack.ProjectProperty(nil),
		Tags: []dtrack.Tag{
			{Name: "foo:bar"},
			{Name: "Owner:foo"},
			{Name: "Owner:bar"},
		},
	}
	expectedToken := "7c78f6c9-token"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() == "/api/v1/bom" {
			body, _ := io.ReadAll(r.Body)
			var actualRequest dtrack.BOMUploadRequest
			json.Unmarshal(body, &actualRequest)
			assert.Equal(t, expectedRequest, actualRequest)
			w.Write([]byte("{\"Token\":\"" + expectedToken + "\"}"))
		} else if r.URL.String() == "/api/v1/project/7c78f6c9-b4b0-493c-a912-0bb0a4f221f1" {
			project := dtrack.Project{
				UUID: projUUID,
				Name: "fooProj",
				PURL: "pkg://npm/xyz/asdf@v1.2.2",
				Tags: []dtrack.Tag{{Name: "foo:bar"}, {Name: "Owner:foo"}},
			}
			res, _ := json.Marshal(project)
			w.Write(res)
			w.WriteHeader(http.StatusOK)
		} else if r.URL.String() == "/api/v1/project" && r.Method == http.MethodPost {
			body, _ := io.ReadAll(r.Body)
			var req dtrack.Project
			json.Unmarshal(body, &req)
			assert.Equal(t, req.Tags, expectedProjectUpdate.Tags)
		} else {
			assert.Fail(t, r.URL.String())
		} // if r.URL.String() == ""
	}))
	defer ts.Close()

	projectUUID = projUUID.String()
	apiKey = "test"
	projectName = "test"
	c, err := dtrack.NewClient(ts.URL, dtrack.WithAPIKey(apiKey))
	assert.Nil(t, err)

	client = c
	issues, err := cyclonedx.ToDracon([]byte(saasBOM), "json")

	assert.Nil(t, err)
	ltr := v1.LaunchToolResponse{
		ToolName: "SAT",
		Issues:   issues,
	}
	eltr := v1.EnrichedLaunchToolResponse{
		OriginalResults: &ltr,
		Issues: []*v1.EnrichedIssue{
			{
				RawIssue: issues[0],
				Annotations: map[string]string{
					"Owner-0": "foo",
					"Owner-1": "bar",
				},
			},
		},
	}
	ownerAnnotation = "Owner"
	tokens, err := uploadBOMSFromEnriched([]*v1.EnrichedLaunchToolResponse{&eltr})
	assert.Nil(t, err)
	assert.Equal(t, tokens, []string{expectedToken})
}

const saasBOM = `{
	"bomFormat": "CycloneDX",
	"specVersion": "1.4",
	"serialNumber": "urn:uuid:3e671687-395b-41f5-a30f-a58921a69b79",
	"version": 1,
	"metadata": {
	  "timestamp": "2021-01-10T12:00:00Z",
	  "component": {
		"bom-ref": "acme-application",
		"type": "application",
		"name": "Acme Cloud Example",
		"version": "2022-1"
	  }
	},
	"services": [
	  {
		"bom-ref": "api-gateway",
		"provider": {
		  "name": "Acme Inc",
		  "url": [ "https://example.com" ]
		},
		"group": "com.example",
		"name": "API Gateway",
		"version": "2022-1",
		"description": "Example API Gateway",
		"endpoints": [
		  "https://example.com/",
		  "https://example.com/app"
		],
		"authenticated": false,
		"x-trust-boundary": true,
		"data": [
		  {
			"classification": "PII",
			"flow": "bi-directional"
		  },
		  {
			"classification": "PIFI",
			"flow": "bi-directional"
		  },
		  {
			"classification": "Public",
			"flow": "bi-directional"
		  }
		],
		"externalReferences": [
		  {
			"type": "documentation",
			"url": "http://example.com/app/swagger"
		  }
		],
		"services": [
		  {
			"bom-ref": "ms-1.example.com",
			"provider": {
			  "name": "Acme Inc",
			  "url": [ "https://example.com" ]
			},
			"group": "com.example",
			"name": "Microservice 1",
			"version": "2022-1",
			"description": "Example Microservice",
			"endpoints": [
			  "https://ms-1.example.com"
			],
			"authenticated": true,
			"x-trust-boundary": false,
			"data": [
			  {
				"classification": "PII",
				"flow": "bi-directional"
			  }
			],
			"externalReferences": [
			  {
				"type": "documentation",
				"url": "https://ms-1.example.com/swagger"
			  }
			]
		  },
		  {
			"bom-ref": "ms-2.example.com",
			"provider": {
			  "name": "Acme Inc",
			  "url": [ "https://example.com" ]
			},
			"group": "com.example",
			"name": "Microservice 2",
			"version": "2022-1",
			"description": "Example Microservice",
			"endpoints": [
			  "https://ms-2.example.com"
			],
			"authenticated": true,
			"x-trust-boundary": false,
			"data": [
			  {
				"classification": "PIFI",
				"flow": "bi-directional"
			  }
			],
			"externalReferences": [
			  {
				"type": "documentation",
				"url": "https://ms-2.example.com/swagger"
			  }
			]
		  },
		  {
			"bom-ref": "ms-3.example.com",
			"provider": {
			  "name": "Acme Inc",
			  "url": [ "https://example.com" ]
			},
			"group": "com.example",
			"name": "Microservice 3",
			"version": "2022-1",
			"description": "Example Microservice",
			"endpoints": [
			  "https://ms-3.example.com"
			],
			"authenticated": true,
			"x-trust-boundary": false,
			"data": [
			  {
				"classification": "Public",
				"flow": "bi-directional"
			  }
			],
			"externalReferences": [
			  {
				"type": "documentation",
				"url": "https://ms-3.example.com/swagger"
			  }
			]
		  },
		  {
			"bom-ref": "ms-1-pgsql.example.com",
			"group": "org.postgresql",
			"name": "Postgres",
			"version": "14.1",
			"description": "Postgres database for Microservice #1",
			"endpoints": [
			  "https://ms-1-pgsql.example.com:5432"
			],
			"authenticated": true,
			"x-trust-boundary": false,
			"data": [
			  {
				"classification": "PII",
				"flow": "bi-directional"
			  }
			]
		  },
		  {
			"bom-ref": "s3-example.amazon.com",
			"group": "com.amazon",
			"name": "S3",
			"description": "S3 bucket",
			"endpoints": [
			  "https://s3-example.amazon.com"
			],
			"authenticated": true,
			"x-trust-boundary": true,
			"data": [
			  {
				"classification": "Public",
				"flow": "bi-directional"
			  }
			]
		  }
		]
	  }
	],
	"dependencies": [
	  {
		"ref": "acme-application",
		"dependsOn": [ "api-gateway" ]
	  },
	  {
		"ref": "api-gateway",
		"dependsOn": [
		  "ms-1.example.com",
		  "ms-2.example.com",
		  "ms-3.example.com"
		]
	  },
	  {
		"ref": "ms-1.example.com",
		"dependsOn": [ "ms-1-pgsql.example.com" ]
	  },
	  {
		"ref": "ms-2.example.com",
		"dependsOn": [ ]
	  },
	  {
		"ref": "ms-3.example.com",
		"dependsOn": [ "s3-example.amazon.com" ]
	  }
	]
  }
  `
