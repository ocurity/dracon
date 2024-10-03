package enrichers

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	draconv1 "github.com/ocurity/dracon/api/proto/v1"
)

// SetupIODirs creates temporary directories for input and output files
func SetupIODirs(t *testing.T) (indir, outdir string) {
	indir, err := os.MkdirTemp("/tmp", "")
	require.NoError(t, err)

	outdir, err = os.MkdirTemp("/tmp", "")
	require.NoError(t, err)

	return indir, outdir
}

// GetEmptyLaunchToolResponse returns a slice of LaunchToolResponse with no issues
func GetEmptyLaunchToolResponse(_ *testing.T) []*draconv1.LaunchToolResponse {
	return []*draconv1.LaunchToolResponse{
		{
			ToolName: "tool1",
			Issues:   []*draconv1.Issue{},
		},
		{
			ToolName: "tool2",
			Issues:   []*draconv1.Issue{},
		},
	}
}

// GetEmptyLaunchToolResponse returns a slice of LaunchToolResponse with no issues
func GetLaunchToolResponse(_ *testing.T) []*draconv1.LaunchToolResponse {
	code := `this
					is
					some
					code`
	return []*draconv1.LaunchToolResponse{
		{
			ToolName: "tool1",
			Issues: []*draconv1.Issue{
				{
					Target:         "file:/a/b/c/d.php:1-2",
					Type:           "sometype",
					Title:          "this is a title",
					Severity:       draconv1.Severity_SEVERITY_CRITICAL,
					Cvss:           1.0,
					Confidence:     draconv1.Confidence_CONFIDENCE_CRITICAL,
					Description:    "this is a handy dandy description",
					Source:         "this is a source",
					Cve:            "CVE-2020-123",
					Uuid:           "d9681ae9-223b-4df8-a422-7b29bb917a36",
					Cwe:            []int32{123},
					ContextSegment: &code,
				},
				{
					Target:         "file:/a/b/c/d.go:2-3",
					Type:           "sometype1",
					Title:          "this is a title1",
					Severity:       draconv1.Severity_SEVERITY_CRITICAL,
					Cvss:           1.0,
					Confidence:     draconv1.Confidence_CONFIDENCE_CRITICAL,
					Description:    "this is a handy dandy description1",
					Source:         "this is a source1",
					Cve:            "CVE-2020-124",
					Uuid:           "a9681ae9-223b-4df8-a422-7b29bb917a36",
					Cwe:            []int32{123},
					ContextSegment: &code,
				},
			},
		},
		{
			ToolName: "tool2",
			Issues: []*draconv1.Issue{
				{
					Target:         "file:/a/b/c/d.py:1-2",
					Type:           "sometype",
					Title:          "this is a title",
					Severity:       draconv1.Severity_SEVERITY_CRITICAL,
					Cvss:           1.0,
					Confidence:     draconv1.Confidence_CONFIDENCE_CRITICAL,
					Description:    "this is a handy dandy description",
					Source:         "this is a source",
					Cve:            "CVE-2020-123",
					Uuid:           "q9681ae9-223b-4df8-a422-7b29bb917a36",
					Cwe:            []int32{123},
					ContextSegment: &code,
				},
				{
					Target:         "file:/a/b/c/d.py:2-3",
					Type:           "sometype1",
					Title:          "this is a title1",
					Severity:       draconv1.Severity_SEVERITY_CRITICAL,
					Cvss:           1.0,
					Confidence:     draconv1.Confidence_CONFIDENCE_CRITICAL,
					Description:    "this is a handy dandy description1",
					Source:         "this is a source1",
					Cve:            "CVE-2020-124",
					Uuid:           "w9681ae9-223b-4df8-a422-7b29bb917a36",
					Cwe:            []int32{123},
					ContextSegment: &code,
				},
			},
		},
	}
}
