package patches

import (
	"fmt"

	"github.com/ocurity/dracon/pkg/types/json6902"
	kustomizeV1Alpha1 "github.com/ocurity/dracon/pkg/types/kustomize.config.k8s.io/v1alpha1"
	tektonV1Beta1 "github.com/ocurity/dracon/pkg/types/tekton.dev/v1beta1"
)

// AddAnchorsToTask generates kustomize pathces that add an anchor to both the parameters and the
// results of the Task, so that it can receive the result of another task and its result can be
// used in other tasks. These fields are are pretty much guaranteed to be present in pretty much
// all out tasks, so this generator is here for convenience to reduce duplication..
func AddAnchorsToTask(task *tektonV1Beta1.Task) ([]kustomizeV1Alpha1.TargetPatch, error) {
	return []kustomizeV1Alpha1.TargetPatch{
		{
			DescriptiveComment: "Add anchors to Task.",
			Target: &kustomizeV1Alpha1.Target{
				Kind: "Task",
				Name: task.Metadata.Name,
			},
			Patch: mustYAMLString([]json6902.Operation{
				{
					Type: json6902.Add,
					Path: "/spec/params/-",
					Value: tektonV1Beta1.TaskSpecParameter{
						Name:        "anchors",
						Description: "A list of tasks that this task depends on using their anchors.",
						Type:        "array",
						Default:     []string{},
					},
				},
				{
					Type: json6902.Add,
					Path: "/spec/results/-",
					Value: tektonV1Beta1.TaskSpecResult{
						Name:        "anchor",
						Description: "An anchor to allow other tasks to depend on this task.",
					},
				},
				{
					Type: json6902.Add,
					Path: "/spec/steps/-",
					Value: tektonV1Beta1.TaskSpecStep{
						Name:   "anchor",
						Image:  "docker.io/busybox:1.35.0",
						Script: `echo "$(context.task.name)" > "$(results.anchor.path)"`,
					},
				},
			}),
		},
	}, nil
}

// Adds scan UUID and start time related parameters to Task that is being processed.
func AddScanUUIDAndStartTimeToTask(task *tektonV1Beta1.Task) ([]kustomizeV1Alpha1.TargetPatch, error) {
	if !hasLabel(task.Metadata, "v1.dracon.ocurity.com/component", "producer") {
		return nil, nil
	}

	var ops []json6902.Operation
	for i := range task.Spec.Steps {
		ops = append(ops, []json6902.Operation{
			{
				Type: json6902.Add,
				Path: fmt.Sprintf("/spec/steps/%d/env/-", i),
				Value: tektonV1Beta1.TaskSpecStepEnv{
					Name:  "DRACON_SCAN_TIME",
					Value: "$(params.dracon_scan_start_time)",
				},
			},
			{
				Type: json6902.Add,
				Path: fmt.Sprintf("/spec/steps/%d/env/-", i),
				Value: tektonV1Beta1.TaskSpecStepEnv{
					Name:  "DRACON_SCAN_ID",
					Value: "$(params.dracon_scan_id)",
				},
			},
		}...)
	}

	return []kustomizeV1Alpha1.TargetPatch{
		{
			DescriptiveComment: "Add scan information to Task.",
			Target: &kustomizeV1Alpha1.Target{
				Kind: "Task",
				Name: task.Metadata.Name,
			},
			Patch: mustYAMLString(append(
				[]json6902.Operation{
					{
						Type: json6902.Add,
						Path: "/spec/params/-",
						Value: tektonV1Beta1.TaskSpecParameter{
							Name: "dracon_scan_id",
							Type: "string",
						},
					},
					{
						Type: json6902.Add,
						Path: "/spec/params/-",
						Value: tektonV1Beta1.TaskSpecParameter{
							Name: "dracon_scan_start_time",
							Type: "string",
						},
					},
				}, ops...)),
		},
	}, nil
}
