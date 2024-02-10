package pipelines

import (
	"errors"
	"fmt"
	"slices"

	"github.com/ocurity/dracon/pkg/components"
	tektonV1Beta1 "github.com/ocurity/dracon/pkg/types/tekton.dev/v1beta1"
)

var _ Backend[*tektonV1Beta1.Pipeline] = (*tektonBackend)(nil)

type tektonBackend struct {
	pipeline *tektonV1Beta1.Pipeline
	tasks    []tektonV1Beta1.Task
	prefix   string
	suffix   string
}

// func renameParameterRef(oldVal string, newParameterNames map[string]string) string {
// 	if !strings.HasPrefix(oldVal, "$(params.") {
// 		return oldVal
// 	}
// 	oldParamRef := strings.Split(oldVal, "$(params.")[1]
// 	// remove last parentheses
// 	oldParamRef = oldParamRef[:len(oldParamRef)-1]
// 	return fmt.Sprintf("$(params.%s)", newParameterNames[oldParamRef])
// }

// fixTaskPrefixSuffix adds a prefix and a suffix to the name of the task and all the task
// parameters. Having task parameters prefixed with the same name as the task itself, helps
// users figure out more easily which parameters configure what.
// func fixTaskPrefixSuffix(task tektonV1Beta1.Task, prefix, suffix string) {
// 	// keep track of renamings so that we can also fix the environment variables
// 	// referencing the parameters
// 	newParameterNames := map[string]string{}
// 	for _, param := range task.Spec.Parameters {
// 		oldParamName := param.Name
// 		paramNameChunks := strings.Split(param.Name, task.Name)
// 		param.Name = prefix + task.Name + suffix + paramNameChunks[1]
// 		newParameterNames[oldParamName] = param.Name
// 	}
// 	// fix references to parameters in step env vars and images
// 	for _, step := range task.Spec.Steps {
// 		for _, env := range step.Env {
// 			env.Value = renameParameterRef(env.Value, newParameterNames)
// 		}
// 		step.Image = renameParameterRef(step.Image, newParameterNames)
// 	}
// 	task.Name = prefix + task.Name + suffix
// }

// addAnchorResult adds an `anchor` entry to the results section of a Task. This helps reduce the
// amount of boilerplate needed to be written by a user to introduce a component.
func addAnchorResult(task tektonV1Beta1.Task) {
	if task.Metadata.Labels[components.LabelKey] == components.Consumer.String() {
		return
	}

	for _, result := range task.Spec.Results {
		if result.Name == "anchor" {
			return
		}
	}

	task.Spec.Results = append(task.Spec.Results, &tektonV1Beta1.TaskSpecResult{
		Name:        "anchor",
		Description: "An anchor to allow other tasks to depend on this task.",
	})
}

// addAnchorParameter adds an `anchors` entry to the parameters of a Task. This entry will then be
// filled in the pipeline with the anchors of the tasks that this task depends on.
func addAnchorParameter(task tektonV1Beta1.Task) {
	if task.Metadata.Labels[components.LabelKey] == components.Source.String() {
		return
	}

	for _, param := range task.Spec.Parameters {
		if param.Name == "anchors" {
			return
		}
	}

	task.Spec.Parameters = append(task.Spec.Parameters, &tektonV1Beta1.TaskSpecParameter{
		Name:        "anchors",
		Description: "A list of tasks that this task depends on",
		Type:        "array",
	})
}

// NewTektonBackend returns an implementation of the Backend interface that will produce a Tekton
// Pipeline object with all the configured tasks.
func NewTektonBackend(basePipeline *tektonV1Beta1.Pipeline, tasks []tektonV1Beta1.Task, prefix, suffix string) (Backend[*tektonV1Beta1.Pipeline], error) {
	if len(tasks) == 0 {
		return nil, errors.New("no tasks provided")
	}

	tektonBackend := &tektonBackend{pipeline: basePipeline, tasks: tasks[:], prefix: prefix, suffix: suffix}
	for _, task := range tasks {
		// TODO(?): revisit if we need this in the future
		// fixTaskPrefixSuffix(task, prefix, suffix)
		addAnchorParameter(task)
		addAnchorResult(task)
	}

	// Sort tasks based on their component type
	slices.SortFunc(tektonBackend.tasks, func(a tektonV1Beta1.Task, b tektonV1Beta1.Task) int {
		componentTypeA := components.MustGetComponentType(a.Metadata.Labels[components.LabelKey])
		componentTypeB := components.MustGetComponentType(b.Metadata.Labels[components.LabelKey])
		return int(componentTypeA) - int(componentTypeB)
	})

	return tektonBackend, nil
}

func (tb *tektonBackend) Generate() (*tektonV1Beta1.Pipeline, error) {
	tb.pipeline.Metadata.Name = tb.prefix + tb.pipeline.Metadata.Name + tb.suffix
	pipelineWorkspaces := map[string]struct{}{}
	anchors := map[string][]string{}

	for _, task := range tb.tasks {
		componentType := task.Metadata.Labels[components.LabelKey]
		anchors[componentType] = append(anchors[componentType], task.Metadata.Name)

		// add task to pipeline tasks
		pipelineTask := &tektonV1Beta1.PipelineSpecTask{
			Name: task.Metadata.Name,
			TaskRef: &tektonV1Beta1.PipelineSpecTaskTaskRef{
				Name: task.Metadata.Name,
			},
		}

		// add task to pipeline's tasks
		tb.pipeline.Spec.Tasks = append(tb.pipeline.Spec.Tasks, pipelineTask)

		if componentType == components.Base.String() {
			continue
		}

		// add task's workspaces to pipeline workspaces
		for _, ws := range task.Spec.Workspaces {
			if _, inserted := pipelineWorkspaces[ws.Name]; !inserted {
				tb.pipeline.Spec.Workspaces = append(tb.pipeline.Spec.Workspaces, &tektonV1Beta1.PipelineSpecWorkspace{
					Name: ws.Name,
				})
				pipelineWorkspaces[ws.Name] = struct{}{}
			}
			pipelineTask.Workspaces = append(pipelineTask.Workspaces, &tektonV1Beta1.PipelineSpecTaskWorkspace{
				Name:      ws.Name,
				Workspace: ws.Name,
			})
		}

		// add task's parameters to pipeline's parameters
		for _, param := range task.Spec.Parameters {
			pipelineSpecTaskParameter := &tektonV1Beta1.PipelineSpecTaskParameter{
				Name: param.Name,
			}

			if param.Name == "anchors" {
				anchorTargetComponentType := components.MustGetComponentType(componentType) - 1
				if len(anchors[anchorTargetComponentType.String()]) > 0 {
					values := []string{}
					for _, anchorTarget := range anchors[anchorTargetComponentType.String()] {
						values = append(values, fmt.Sprintf("$(tasks.%s.results.anchor)", anchorTarget))
					}

					pipelineSpecTaskParameter.Value = values
				}
			} else {
				switch param.Type {
				case "array":
					pipelineSpecTaskParameter.Value = []string{fmt.Sprintf("$(params.%s)", param.Name)}
				case "string":
					pipelineSpecTaskParameter.Value = fmt.Sprintf("$(params.%s)", param.Name)
				}
			}

			pipelineTask.Parameters = append(pipelineTask.Parameters, pipelineSpecTaskParameter)

			// do not add anchor to pipeline params
			if param.Name == "anchors" {
				continue
			}

			tb.pipeline.Spec.Parameters = append(tb.pipeline.Spec.Parameters, &tektonV1Beta1.PipelineSpecParameter{
				Name:        param.Name,
				Type:        param.Type,
				Description: param.Description,
				Default:     param.Default,
			})
		}

		// add scan ID and scan time to all  producers
		if task.Metadata.Labels[components.LabelKey] != components.Producer.String() {
			continue
		}

		pipelineTask.Parameters = append(pipelineTask.Parameters, []*tektonV1Beta1.PipelineSpecTaskParameter{
			{
				Name:  "dracon_scan_id",
				Value: fmt.Sprintf("$(tasks.%s.results.dracon-scan-id)", anchors[components.Base.String()][0]),
			},
			{
				Name:  "dracon_scan_start_time",
				Value: fmt.Sprintf("$(tasks.%s.results.dracon-scan-start-time)", anchors[components.Base.String()][0]),
			},
			{
				Name:  "dracon_scan_tags",
				Value: fmt.Sprintf("$(tasks.%s.results.dracon-scan-tags)", anchors[components.Base.String()][0]),
			},
		}...)

		task.Spec.Parameters = append(task.Spec.Parameters, []*tektonV1Beta1.TaskSpecParameter{
			{
				Name: "dracon_scan_id",
				Type: "string",
			},
			{
				Name: "dracon_scan_start_time",
				Type: "string",
			},
			{
				Name: "dracon_scan_tags",
				Type: "string",
			},
		}...)

		// add environment variables for every step of a task to set the DRACON_SCAN_TIME and DRACON_SCAN_ID
		for _, step := range task.Spec.Steps {
			step.Env = append(step.Env, []*tektonV1Beta1.TaskSpecStepEnv{
				{
					Name:  "DRACON_SCAN_TIME",
					Value: "$(params.dracon_scan_start_time)",
				},
				{
					Name:  "DRACON_SCAN_ID",
					Value: "$(params.dracon_scan_id)",
				},
				{
					Name:  "DRACON_SCAN_TAGS",
					Value: "$(params.dracon_scan_tags)",
				},
			}...)
		}
	}

	tb.pipeline.Spec.Parameters = append(tb.pipeline.Spec.Parameters, []*tektonV1Beta1.PipelineSpecParameter{
		{
			Name: "dracon_scan_id",
			Type: "string",
		},
		{
			Name: "dracon_scan_start_time",
			Type: "string",
		},
		{
			Name: "dracon_scan_tags",
			Type: "string",
		},
	}...)

	return tb.pipeline, nil
}
