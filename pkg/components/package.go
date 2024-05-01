package components

import (
	"context"
	stderrors "errors"
	"os"
	"path"
	"path/filepath"

	"github.com/go-errors/errors"
	tektonv1beta1api "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"sigs.k8s.io/yaml"

	"github.com/ocurity/dracon/pkg/manifests"
)

var (
	// ErrCouldNotStatPath is returned when the path to the components doesn't
	// exist or is inaccesbile
	ErrCouldNotStatPath = errors.New("could not stat path")
	// ErrNotADirectory is returned when the path is not pointing to a
	// directory
	ErrNotADirectory = errors.New("path is not a directory")
)

// Package explores the components folder provided and gathers all the Tekton
// Tasks into one Helm chart.
func Package(ctx context.Context, name, componentFolder string, draconVersion string, chartVersion string) (err error) {
	fs, err := os.Stat(componentFolder)
	if err != nil {
		return errors.Errorf("%s: could not stat: %w", componentFolder, err)
	}

	if !fs.IsDir() {
		return errors.Errorf("%s: path is not a directory", componentFolder)
	}

	tempFolder, err := os.MkdirTemp("/tmp", "dracon-helm")
	if err != nil {
		return errors.Errorf("there was an error while trying to create temp directory: %w", err)
	}
	defer func() {
		err = stderrors.Join(err, os.RemoveAll(tempFolder))
	}()

	taskPaths, err := gatherTasks(componentFolder)
	if err != nil {
		return errors.Errorf("could not discover tasks: %w", err)
	}
	if err = constructPackage(ctx, tempFolder, name, chartVersion, draconVersion, taskPaths); err != nil {
		return errors.Errorf("could not generate Helm manifests: %w", err)
	}

	taskChart, err := loader.LoadDir(tempFolder)
	if err != nil {
		return errors.Errorf("could not load tasks into a Helm chart: %w", err)
	}

	_, err = chartutil.Save(taskChart, ".")
	if err != nil {
		return errors.Errorf("could not create package tar file: %w", err)
	}

	return nil
}

// constructPackage creates a templates folder with all the discovered tasks
// in at the designated Path, along with a Chart file
//
//revive:disable:cyclomatic High complexity score but easy to understand
//revive:disable:cognitive-complexity High complexity score but easy to understand
func constructPackage(ctx context.Context, helmFolder, name, version, appVersion string, taskPaths []string) error {
	taskList := []*tektonv1beta1api.Task{}
	for _, taskFile := range taskPaths {
		task, err := manifests.LoadTektonV1Beta1Task(ctx, ".", taskFile)
		if err != nil {
			return errors.Errorf("%s: not a valid manifest: %w", taskFile, err)
		}
		taskList = append(taskList, task)
	}

	if err := os.Mkdir(path.Join(helmFolder, "templates"), os.ModePerm); err != nil {
		return errors.Errorf("could not create templates folder")
	}

	tasksFilePath := path.Join(helmFolder, "templates", "tasks.yaml")
	tasksFile, err := os.OpenFile(tasksFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return errors.Errorf("could not create tasks file: %w", err)
	}

	for _, task := range taskList {
		addAnchorParameter(task)
		addAnchorResult(task)

		if err = manifests.TektonV1Beta1ObjEncoder.Encode(task, tasksFile); err != nil {
			return errors.Errorf("could not store task %s: %w", task.Name, err)
		}

		if _, err = tasksFile.WriteString("---\n"); err != nil {
			return errors.Errorf("%s: could not write to file: %w", tasksFilePath, err)
		}
	}

	chartFile := path.Join(helmFolder, "Chart.yaml")
	_, err = os.Create(chartFile)
	if err != nil {
		return errors.Errorf("could not create Chart file: %w", err)
	}

	helmChart := chart.Metadata{
		APIVersion: chart.APIVersionV2,
		Name:       name,
		Version:    version,
		AppVersion: appVersion,
	}

	if err = helmChart.Validate(); err != nil {
		return errors.Errorf("there was an issue generating the Helm metadata: %w", err)
	}

	helmMetadataBytes, err := yaml.Marshal(helmChart)
	if err != nil {
		return errors.Errorf("could not marshal Helm matadata: %w", err)
	}

	err = os.WriteFile(chartFile, helmMetadataBytes, os.ModeAppend)
	if err != nil {
		return errors.Errorf("could not write ")
	}

	return nil
}

// gatherTasks returns the paths of all the Tekton Tasks discovered
//
//revive:disable:cognitive-complexity High complexity score but easy to understand
func gatherTasks(folder string) ([]string, error) {
	taskPaths := []string{}

	for _, componentType := range []string{"sources", "producers", "enrichers", "consumers"} {
		componentsFolder := path.Join(folder, componentType)
		stat, err := os.Stat(componentsFolder)
		if err != nil {
			continue
		}

		if !stat.IsDir() {
			return nil, errors.Errorf("%s: %w", componentsFolder, ErrNotADirectory)
		}

		err = filepath.Walk(componentsFolder, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return errors.Errorf("%s: %w", path, err)
			}

			if info == nil {
				return errors.Errorf("path %s doesn't exist", path)
			}

			if !info.IsDir() && info.Name() == "task.yaml" {
				taskPaths = append(taskPaths, path)
			}
			return nil
		})

		if err != nil {
			return nil, errors.Errorf("there was an error while gathering tasks: %w", err)
		}
	}
	return taskPaths, nil
}

// TODO(ptzianos): these are copies of the functino in the pipelines package.
// eventually we need to merge the two and clean them up
// addAnchorResult adds an `anchor` entry to the results section of a Task.
// This helps reduce the amount of boilerplate needed to be written by a user
// to introduce a component. The base task doesn't need an anchor because its
// output it a dependency for the consumer tasks.
func addAnchorResult(task *tektonv1beta1api.Task) {
	if task.Labels[LabelKey] == Consumer.String() || task.Labels[LabelKey] == Base.String() {
		return
	}

	task.Spec.Results = append(task.Spec.Results, tektonv1beta1api.TaskResult{
		Name:        "anchor",
		Description: "An anchor to allow other tasks to depend on this task.",
	})

	task.Spec.Steps = append(task.Spec.Steps, tektonv1beta1api.Step{
		Name:   "anchor",
		Image:  "docker.io/busybox",
		Script: "echo \"$(context.task.name)\" > \"$(results.anchor.path)\"",
	})
}

// addAnchorParameter adds an `anchors` entry to the parameters of a Task. This entry will then be
// filled in the pipeline with the anchors of the tasks that this task depends on.
func addAnchorParameter(task *tektonv1beta1api.Task) {
	componentType, err := ToComponentType(task.Labels[LabelKey])
	if err != nil {
		panic(errors.Errorf("%s: %w", task.Name, err))
	}
	if componentType < Producer {
		return
	}

	for _, param := range task.Spec.Params {
		if param.Name == "anchors" {
			return
		}
	}

	task.Spec.Params = append(task.Spec.Params, tektonv1beta1api.ParamSpec{
		Name:        "anchors",
		Description: "A list of tasks that this task depends on",
		Type:        "array",
		Default: &tektonv1beta1api.ParamValue{
			Type: tektonv1beta1api.ParamTypeArray,
		},
	})
}
