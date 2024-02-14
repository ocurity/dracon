package patches

import (
	kustomize "github.com/ocurity/dracon/build/tools/kustomize-component-generator/types/kustomize.config.k8s.io/v1alpha1"
	tekton "github.com/ocurity/dracon/build/tools/kustomize-component-generator/types/tekton.dev/v1beta1"
)

// AddScanUUIDAndStartTimeToTask implements Patch for adding DAG anchors to every Task.
type AddScanUUIDAndStartTimeToTask struct {
	task *tekton.Task
}

// NewAddScanUUIDAndStartTimeToTask returns a new implementation of Patch for adding DAG
// anchors to every Task.
func NewAddScanUUIDAndStartTimeToTask(task *tekton.Task) *AddScanUUIDAndStartTimeToTask {
	return &AddScanUUIDAndStartTimeToTask{
		task: task,
	}
}

// GeneratePatch implements Patch.GeneratePatch.
func (p *AddScanUUIDAndStartTimeToTask) GeneratePatch() *kustomize.TargetPatch {
	if !hasLabel(p.task.Metadata, "v1.dracon.ocurity.com/component", "producer") {
		return nil
	}
	task := &tekton.Task{
		GVK:      p.task.GVK,
		Metadata: p.task.Metadata,
		Spec: &tekton.TaskSpec{
			Parameters: []*tekton.TaskSpecParameter{
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
			},
		},
	}
	requiredEnv := []*tekton.TaskSpecStepEnv{
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
	}
	var steps []*tekton.TaskSpecStep
	for _, step := range p.task.Spec.Steps {
		s := step
		s.Env = append(s.Env, requiredEnv...)
		steps = append(steps, s)
	}
	task.Spec.Steps = steps
	return &kustomize.TargetPatch{
		DescriptiveComment: "Add scan information to Task.",
		Target: &kustomize.Target{
			Kind: "Task",
			Name: p.task.Metadata.Name,
		},
		Patch: mustYAMLString(task),
	}
}
