package pipelines

import (
	"context"
	"testing"
	"time"

	"github.com/ocurity/dracon/pkg/manifests"
	"github.com/stretchr/testify/require"
	tektonV1Beta1API "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kustomizeTypes "sigs.k8s.io/kustomize/api/types"
)

func TestResolveKustomizationResourceBases(t *testing.T) {
	testCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	baseTask, err := manifests.LoadTektonV1Beta1Task(testCtx, ".", "../../components/base")
	require.NoError(t, err)

	gitTask, err := manifests.LoadTektonV1Beta1Task(testCtx, ".", "../../components/sources/git")
	require.NoError(t, err)

	testCases := []struct {
		name             string
		kustomization    *Kustomization
		expectedPipeline *tektonV1Beta1API.Pipeline
		expectedTasks    []*tektonV1Beta1API.Task
		expectedErr      error
	}{
		{
			name: "no base pipeline in the kustomization",
			kustomization: &Kustomization{
				Kustomization: &kustomizeTypes.Kustomization{
					NameSuffix: "-cyberdyne-card-processing",
					Resources:  []string{},
					Components: []string{},
				},
			},
			expectedErr: ErrKustomizationMissingBaseResources,
		},
		{
			name: "no components in the kustomization",
			kustomization: &Kustomization{
				Kustomization: &kustomizeTypes.Kustomization{
					NameSuffix: "-cyberdyne-card-processing",
					Resources: []string{
						"https://raw.githubusercontent.com/ocurity/dracon/main/components/base/pipeline.yaml",
						"https://raw.githubusercontent.com/ocurity/dracon/main/components/base/task.yaml",
					},
					Components: []string{},
				},
			},
			expectedErr: ErrNoComponentsInKustomization,
		},
		{
			name: "success",
			kustomization: &Kustomization{
				Kustomization: &kustomizeTypes.Kustomization{
					NameSuffix: "-cyberdyne-card-processing",
					Resources: []string{
						"https://raw.githubusercontent.com/ocurity/dracon/main/components/base/pipeline.yaml",
						"https://raw.githubusercontent.com/ocurity/dracon/main/components/base/task.yaml",
					},
					Components: []string{"../../components/sources/git"},
				},
			},
			expectedPipeline: &tektonV1Beta1API.Pipeline{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Pipeline",
					APIVersion: "tekton.dev/v1beta1",
				},
				ObjectMeta: metav1.ObjectMeta{Name: "dracon"},
				Spec: tektonV1Beta1API.PipelineSpec{
					Tasks:      []tektonV1Beta1API.PipelineTask{},
					Workspaces: []tektonV1Beta1API.PipelineWorkspaceDeclaration{},
				},
			},
			expectedTasks: []*tektonV1Beta1API.Task{baseTask, gitTask},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			runCtx, cancel := context.WithCancel(testCtx)
			defer cancel()

			basePipeline, taskList, err := testCase.kustomization.ResolveKustomizationResources(runCtx)
			require.ErrorIs(t, err, testCase.expectedErr)
			require.Equal(t, testCase.expectedPipeline, basePipeline)
			require.EqualValues(t, testCase.expectedTasks, taskList)
		})
	}
}
