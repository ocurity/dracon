package components

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	tektonv1beta1api "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"

	"github.com/ocurity/dracon/pkg/manifests"
)

func TestGatherTasks(t *testing.T) {
	taskPaths, err := GatherTasks("testdata")
	require.NoError(t, err)
	require.EqualValues(t,
		[]string{
			"testdata/base/task.yaml",
			"testdata/sources/git/task.yaml",
			"testdata/producers/aggregator/task.yaml",
			"testdata/producers/cdxgen/task.yaml",
			"testdata/enrichers/aggregator/task.yaml",
			"testdata/consumers/arangodb/task.yaml",
		},
		taskPaths,
	)

	taskPaths, err = GatherTasks("..")
	require.NoError(t, err)
	require.Empty(t, taskPaths)
}

func TestCreateHelmPackage(t *testing.T) {
	draconVersion := "v0.10.0"
	semVer := "0.10.0"
	taskList, err := LoadTasks(
		context.Background(),
		[]string{
			"testdata/base/task.yaml",
			"testdata/sources/git/task.yaml",
			"testdata/producers/cdxgen/task.yaml",
			"testdata/producers/aggregator/task.yaml",
			"testdata/enrichers/aggregator/task.yaml",
			"testdata/consumers/arangodb/task.yaml",
		},
	)
	require.NoError(t, err)

	require.NoError(t, ProcessTasks(taskList...))
	require.Len(t, taskList[2].Spec.Steps[1].Env, 3)
	require.Equal(t, "DRACON_SCAN_TIME", taskList[2].Spec.Steps[1].Env[0].Name)
	require.Equal(t, "DRACON_SCAN_ID", taskList[2].Spec.Steps[1].Env[1].Name)
	require.Equal(t, "DRACON_SCAN_TAGS", taskList[2].Spec.Steps[1].Env[2].Name)
	paramLen := len(taskList[2].Spec.Params)
	require.Equal(t, taskList[2].Spec.Params[paramLen-4].Name, "anchors")
	require.Equal(t, taskList[2].Spec.Params[paramLen-3].Name, "dracon_scan_id")
	require.Equal(t, taskList[2].Spec.Params[paramLen-2].Name, "dracon_scan_start_time")
	require.Equal(t, taskList[2].Spec.Params[paramLen-1].Name, "dracon_scan_tags")

	helmFolder := t.TempDir()
	require.NoError(t, constructPackage(helmFolder, "dracon-oss-components", semVer, draconVersion, taskList))
	require.FileExists(t, path.Join(helmFolder, "Chart.yaml"))
	chartFileContents, err := os.ReadFile(path.Join(helmFolder, "Chart.yaml"))
	require.NoError(t, err)
	require.Equal(
		t,
		fmt.Sprintf(`apiVersion: v2
appVersion: %s
name: dracon-oss-components
version: %s
`, draconVersion, semVer),
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

		task, isATask := obj.(*tektonv1beta1api.Task)
		require.True(t, isATask)

		baseOrSource, err := LabelValueOneOf(task.Labels, Base, Source)
		require.NoError(t, err)

		if baseOrSource {
			goto checkSteps
		}

		for _, param := range task.Spec.Params {
			if param.Name == "anchors" {
				goto checkSteps
			}
		}
		t.Fatalf("task %s has no anchor parameter", task.Name)

	checkSteps:
		baseOrConsumer, err := LabelValueOneOf(task.Labels, Base, Consumer)
		require.NoError(t, err)
		if baseOrConsumer {
			continue
		}

		require.Equal(t, "anchor", task.Spec.Steps[len(task.Spec.Steps)-1].Name, "task %s has no anchor step", task.Name)
	}
}
