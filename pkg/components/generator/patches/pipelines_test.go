package patches

import (
	"testing"

	"github.com/ocurity/dracon/pkg/types/kubernetes"
	tektonV1Beta1 "github.com/ocurity/dracon/pkg/types/tekton.dev/v1beta1"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func getMockTask(labels map[string]string) *tektonV1Beta1.Task {
	return &tektonV1Beta1.Task{
		Metadata: &kubernetes.Metadata{
			Name:   "testTask",
			Labels: labels,
		},
		Spec: &tektonV1Beta1.TaskSpec{
			Workspaces: []*tektonV1Beta1.TaskSpecWorkspace{
				{
					Name:        "source-code-ws",
					Description: "The workspace containing the source-code to scan.",
				},
			},
			Parameters: []*tektonV1Beta1.TaskSpecParameter{
				{
					Name:        "producer-aggregator-anchors",
					Description: "A list of tasks that this task depends on",
					Type:        "array",
				},
			},
		},
	}
}

var expectedRes = `apiVersion: kustomize.config.k8s.io/v1alpha1
kind: Component
resources:
  - task.yaml
patches:
  - target:
      kind: Pipeline
    patch: |
      - path: /spec/workspaces/-
        op: add
        value:
          name: source-code-ws
      - path: /spec/params/-
        op: add
        value:
          name: producer-aggregator-anchors
          type: array
          description: A list of tasks that this task depends on.
          default: []
      - path: /spec/tasks/-
        op: add
        value:
          name: producer-aggregator
          taskRef:
            name: producer-aggregator
          workspaces:
            - name: source-code-ws
              workspace: source-code-ws
          params:
            - name: producer-aggregator-anchors
              value:
                - $(params.producer-aggregator-anchors)
`

func TestAddTaskToPipeline(t *testing.T) {
	kustomizePatches, err := AddTaskToPipeline(getMockTask(nil))
	require.NoError(t, err)
	out, err := yaml.Marshal(kustomizePatches)
	t.Log(string(out))
	require.Error(t, err)
}
