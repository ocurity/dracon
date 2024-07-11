package enrichers

import (
	"os"
	"testing"

	v1 "github.com/ocurity/dracon/api/proto/v1"
	"github.com/stretchr/testify/require"
)

func SetupIODirs(t *testing.T) (indir, outdir string) {
	indir, err := os.MkdirTemp("/tmp", "")
	require.NoError(t, err)

	outdir, err = os.MkdirTemp("/tmp", "")
	require.NoError(t, err)

	return indir, outdir
}

func GetEmptyLaunchToolResponse(t *testing.T) []*v1.LaunchToolResponse {
	return []*v1.LaunchToolResponse{
		{
			ToolName: "tool1",
			Issues:   []*v1.Issue{},
		},
		{
			ToolName: "tool2",
			Issues:   []*v1.Issue{},
		},
	}
}
