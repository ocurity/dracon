package enrichers

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/timestamppb"

	draconapiv1 "github.com/ocurity/dracon/api/proto/v1"
	v1 "github.com/ocurity/dracon/api/proto/v1"
	"github.com/ocurity/dracon/pkg/putil"
)

func createObjects() *draconapiv1.EnrichedLaunchToolResponse {
	scanID := "7c78f6c9-b4b0-493c-a912-0bb0a4aaaaa0"
	times, _ := time.Parse(time.RFC3339, "2023-01-19T18:09:06.370037788Z")
	timestamp := timestamppb.New(times)
	si := v1.ScanInfo{
		ScanUuid:      scanID,
		ScanStartTime: timestamp,
	}
	toolName := "SAT-Tool"
	response := v1.LaunchToolResponse{
		ToolName: toolName,
		ScanInfo: &si,
	}
	enrichedResponse := v1.EnrichedLaunchToolResponse{}

	var issues []*v1.Issue
	var enrichedIssues []*v1.EnrichedIssue
	for j := 0; j < 10; j++ {
		id := uuid.New()
		x := v1.Issue{
			Target:     fmt.Sprintf("target-%d", j),
			Type:       fmt.Sprintf("type-%d", j),
			Title:      fmt.Sprintf("title-%d", j),
			Severity:   v1.Severity_SEVERITY_INFO,
			Confidence: v1.Confidence_CONFIDENCE_INFO,
			Uuid:       id.String(),
		}
		y := v1.EnrichedIssue{
			RawIssue:      &x,
			FirstSeen:     response.ScanInfo.ScanStartTime,
			Count:         uint64(j),
			FalsePositive: false,
			UpdatedAt:     response.ScanInfo.ScanStartTime,
			Hash:          fmt.Sprintf("d41d8cd98f00b204e9800998ecf842%d", j),
		}
		issues = append(issues, &x)
		enrichedIssues = append(enrichedIssues, &y)
	}
	response.Issues = issues
	enrichedResponse.OriginalResults = &response
	enrichedResponse.Issues = enrichedIssues
	return &enrichedResponse
}

func TestWriteDataNormalOperation(t *testing.T) {
	enricherName := "tests-enricher"
	// prepare
	workdir, err := os.MkdirTemp("/tmp", "")
	require.NoError(t, err)
	require.NoError(t, os.Mkdir(filepath.Join(workdir, "raw"), 0755))

	// test errors first
	enrichedResponse := createObjects()
	require.Error(t, WriteData(nil, enricherName))

	enrichedResponse.Issues = []*v1.EnrichedIssue{}
	require.Error(t, WriteData(enrichedResponse, enricherName))

	enrichedResponse = createObjects() // reset

	enrichedResponse.OriginalResults.Issues = []*v1.Issue{}
	require.Error(t, WriteData(enrichedResponse, enricherName))

	enrichedResponse.OriginalResults = nil
	require.Error(t, WriteData(enrichedResponse, enricherName))

	// happy path
	enrichedResponse = createObjects()
	SetWritePathForTests(workdir)
	require.NoError(t, WriteData(enrichedResponse, enricherName))

	require.NoError(t, err)
	er, err := putil.LoadEnrichedNonAggregatedToolResponse(workdir)
	require.NoError(t, err)
	opt := cmp.Comparer(func(x, y timestamppb.Timestamp) bool {
		return x.Nanos == y.Nanos
	})

	require.True(t, cmp.Equal([]*v1.EnrichedLaunchToolResponse{enrichedResponse}, er, protocmp.Transform(), opt),
		cmp.Diff([]*v1.EnrichedLaunchToolResponse{enrichedResponse}, er, protocmp.Transform()))

	r, err := putil.LoadToolResponse(filepath.Join(workdir, fmt.Sprintf("%s.raw.pb", enrichedResponse.GetOriginalResults().GetToolName())))
	require.NoError(t, err)

	require.True(t, cmp.Equal([]*v1.LaunchToolResponse{enrichedResponse.OriginalResults}, r, protocmp.Transform(), opt),
		cmp.Diff([]*v1.LaunchToolResponse{enrichedResponse.OriginalResults}, r, protocmp.Transform()))
}
