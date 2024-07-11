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
