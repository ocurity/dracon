package pipelines

import (
	"github.com/go-errors/errors"
	tektonv1beta1api "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"

	"github.com/ocurity/dracon/pkg/components"
)

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
func addAnchorResult(task *tektonv1beta1api.Task) {
	noResultAnchorNeeded, err := components.LabelValueOneOf(task.Labels, components.Consumer, components.Base)
	if err != nil {
		panic(err)
	} else if noResultAnchorNeeded {
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
	componentType, err := components.ToComponentType(task.Labels[components.LabelKey])
	if err != nil {
		panic(errors.Errorf("%s: %w", task.Name, err))
	}
	if componentType < components.Producer {
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
