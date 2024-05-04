package pipelines

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	tektonv1beta1api "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	kustomizetypes "sigs.k8s.io/kustomize/api/types"

	"github.com/ocurity/dracon/pkg/manifests"
)

func TestResolveKustomizationResourceBase(t *testing.T) {
	testCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	basePipelineFromFile, err := manifests.LoadTektonV1Beta1Pipeline(testCtx, ".", "../../components/base")
	require.NoError(t, err)

	BasePipeline.Labels = map[string]string{"TestResolveKustomizationResourceBases": "bla"}
	defer func() {
		BasePipeline.Labels = nil
		BaseTask.Labels = nil
	}()

	testCases := []struct {
		name             string
		kustomization    *kustomizetypes.Kustomization
		expectedPipeline *tektonv1beta1api.Pipeline
		expectedErr      error
	}{
		{
			name: "no base pipeline in the kustomization",
			kustomization: &kustomizetypes.Kustomization{
				NameSuffix: "-cyberdyne-card-processing",
				Resources:  []string{},
			},
			expectedPipeline: BasePipeline,
		},
		{
			name: "success",
			kustomization: &kustomizetypes.Kustomization{
				NameSuffix: "-cyberdyne-card-processing",
				Resources: []string{
					"../../components/base/pipeline.yaml",
				},
			},
			expectedPipeline: basePipelineFromFile,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			runCtx, cancel := context.WithCancel(testCtx)
			defer cancel()

			basePipeline, err := ResolveBase(runCtx, ".", *testCase.kustomization)
			require.ErrorIs(t, err, testCase.expectedErr)
			require.Equal(t, testCase.expectedPipeline, basePipeline)
		})
	}
}
