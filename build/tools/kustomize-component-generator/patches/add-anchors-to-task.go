package patches

import (
	kustomize "github.com/ocurity/dracon/build/tools/kustomize-component-generator/types/kustomize.config.k8s.io/v1alpha1"
	tekton "github.com/ocurity/dracon/build/tools/kustomize-component-generator/types/tekton.dev/v1beta1"
)

// AddAnchorsToTask implements Patch for adding DAG anchors to every Task.
type AddAnchorsToTask struct {
	task *tekton.Task
}

// NewAddAnchorsToTask returns a new implementation of Patch for adding DAG
// anchors to every Task.
func NewAddAnchorsToTask(task *tekton.Task) *AddAnchorsToTask {
	return &AddAnchorsToTask{
		task: task,
	}
}

// GeneratePatch implements Patch.GeneratePatch.
func (p *AddAnchorsToTask) GeneratePatch() *kustomize.TargetPatch {
	task := &tekton.Task{
		GVK:      p.task.GVK,
		Metadata: p.task.Metadata,
		Spec: &tekton.TaskSpec{
			Parameters: []*tekton.TaskSpecParameter{
				{
					Name:        "anchors",
					Description: "A list of tasks that this task depends on using their anchors.",
					Type:        "array",
					Default:     []string{},
				},
			},
			Results: []*tekton.TaskSpecResult{
				{
					Name:        "anchor",
					Description: "An anchor to allow other tasks to depend on this task.",
				},
			},
			Steps: []*tekton.TaskSpecStep{
				{
					Name:   "anchor",
					Image:  "docker.io/busybox:1.35.0",
					Script: `echo "$(context.task.name)" > "$(results.anchor.path)"`,
				},
			},
		},
	}

	return &kustomize.TargetPatch{
		DescriptiveComment: "Add anchors to Task.",
		Target: &kustomize.Target{
			Kind: "Task",
			Name: p.task.Metadata.Name,
		},
		Patch: mustYAMLString(task),
	}
}
