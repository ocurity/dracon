package components

import (
	"context"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	tektonv1beta1api "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/ocurity/dracon/pkg/manifests"
)

func TestGatherTasks(t *testing.T) {
	taskPaths, err := gatherTasks("testdata")
	require.NoError(t, err)
	require.EqualValues(t,
		[]string{
			"testdata/base/task.yaml",
			"testdata/sources/git/task.yaml",
			"testdata/producers/aggregator/task.yaml",
			"testdata/enrichers/aggregator/task.yaml",
			"testdata/consumers/arangodb/task.yaml",
		},
		taskPaths,
	)

	taskPaths, err = gatherTasks("..")
	require.NoError(t, err)
	require.Empty(t, taskPaths)
}

func TestCreateHelmPackage(t *testing.T) {
	helmFolder := t.TempDir()
	err := constructPackage(
		context.Background(),
		helmFolder,
		"dracon-oss-components",
		"0.1.0",
		"0.10.0",
		[]string{
			"testdata/base/task.yaml",
			"testdata/sources/git/task.yaml",
			"testdata/producers/aggregator/task.yaml",
			"testdata/enrichers/aggregator/task.yaml",
			"testdata/consumers/arangodb/task.yaml",
		},
	)
	require.NoError(t, err)
	require.FileExists(t, path.Join(helmFolder, "Chart.yaml"))
	chartFileContents, err := os.ReadFile(path.Join(helmFolder, "Chart.yaml"))
	require.NoError(t, err)
	require.Equal(
		t,
		`apiVersion: v2
appVersion: 0.10.0
name: dracon-oss-components
version: 0.1.0
`,
		string(chartFileContents),
	)
	require.FileExists(t, path.Join(path.Join(helmFolder, "templates", "tasks.yaml")))
	taskTemplateContents, err := os.ReadFile(path.Join(helmFolder, "templates", "tasks.yaml"))
	require.NoError(t, err)

	for _, rawTask := range strings.Split(string(taskTemplateContents), "---\n") {
		// the Split function returns an empty string at the end
		if rawTask == "" {
			continue
		}
		obj, gKV, err := manifests.K8sObjDecoder.Decode([]byte(rawTask), nil, nil)
		require.NoError(t, err)
		require.Equal(t, *gKV, tektonv1beta1api.SchemeGroupVersion.WithKind("Task"))

		task := obj.(*tektonv1beta1api.Task)
		if task.Labels[LabelKey] == Base.String() || task.Labels[LabelKey] == Source.String() {
			goto checkSteps
		}

		for _, param := range task.Spec.Params {
			if param.Name == "anchors" {
				goto checkSteps
			}
		}
		t.Fatalf("task %s has no anchor parameter", task.Name)

	checkSteps:
		if task.Labels[LabelKey] == Consumer.String() || task.Labels[LabelKey] == Base.String() {
			continue
		}

		require.Equal(t, "anchor", task.Spec.Steps[len(task.Spec.Steps)-1].Name, "task %s has no anchor step", task.Name)
	}
}

func TestImagePinning(t *testing.T) {
	task := &tektonv1beta1api.Task{
		ObjectMeta: metav1.ObjectMeta{
			Name: "testImageFixing",
		},
		Spec: tektonv1beta1api.TaskSpec{
			Steps: []tektonv1beta1api.Step{
				{
					Name:  "test_pinned_to_latest",
					Image: "{{ default 'ghcr.io/ocurity/dracon' .Values.container_registry }}/components/enrichers/aggregator:latest",
				},
				{
					Name:  "test_pinned_to_some_version",
					Image: "docker.io/library/buildpack-deps:stable-curl@sha256:3d5e59c47d5f82a769ad3f372cc9f86321e2e2905141bba974b75d3c08a53e8e",
				},
				{
					Name:  "test_image_is_a_parameter",
					Image: "$(taskName.param.some-image)",
				},
			},
		},
	}
	fixImageVersion(task, "1.0.1")
	require.Equal(t, "{{ default 'ghcr.io/ocurity/dracon' .Values.container_registry }}/components/enrichers/aggregator:1.0.1", task.Spec.Steps[0].Image)
	require.Equal(t, "docker.io/library/buildpack-deps:stable-curl@sha256:3d5e59c47d5f82a769ad3f372cc9f86321e2e2905141bba974b75d3c08a53e8e", task.Spec.Steps[1].Image)
	require.Equal(t, "$(taskName.param.some-image)", task.Spec.Steps[2].Image)
}
