package patches

import (
	"fmt"
	"log"
	"strings"

	"github.com/ocurity/dracon/pkg/types/json6902"
	kustomizeV1Alpha1 "github.com/ocurity/dracon/pkg/types/kustomize.config.k8s.io/v1alpha1"
	tektonV1Beta1 "github.com/ocurity/dracon/pkg/types/tekton.dev/v1beta1"
)

// AddTaskToPipeline generates JSON 6902 operations that add the task's parameters and workspaces
// to the Pipeline definition. It also generates a JSON 6902 operation that will add a task
// reference to the task in the Pipeline's steps.
func AddTaskToPipeline(task *tektonV1Beta1.Task) ([]kustomizeV1Alpha1.TargetPatch, error) {
	patches := []kustomizeV1Alpha1.TargetPatch{}
	pipelineTask := tektonV1Beta1.PipelineSpecTask{
		Name: task.Metadata.Name,
		TaskRef: &tektonV1Beta1.PipelineSpecTaskTaskRef{
			Name: task.Metadata.Name,
		},
	}
	jsonOps := []json6902.Operation{}

	for _, ws := range task.Spec.Workspaces {
		jsonOps = append(jsonOps, json6902.Operation{
			Type:  json6902.Add,
			Path:  "/spec/workspaces/-",
			Value: map[string]string{"name": ws.Name},
		})

		pipelineTask.Workspaces = append(pipelineTask.Workspaces, &tektonV1Beta1.PipelineSpecTaskWorkspace{
			Name:      ws.Name,
			Workspace: ws.Name,
		})
	}

	for _, param := range task.Spec.Parameters {
		if !strings.HasPrefix(param.Name, task.Metadata.Name) {
			log.Fatalf("Parameter '%s' in '%s/%s/%s' should start with '%s'.",
				param.Name, task.GVK.APIVersion, task.GVK.Kind, task.Metadata.Name, task.Metadata.Name,
			)
		}

		jsonOps = append(jsonOps, json6902.Operation{
			Type:  json6902.Add,
			Path:  "/spec/params/-",
			Value: param,
		})

		pipelineSpecTaskParameter := &tektonV1Beta1.PipelineSpecTaskParameter{
			Name: param.Name,
		}

		switch param.Type {
		case "array":
			pipelineSpecTaskParameter.Value = []string{fmt.Sprintf("$(params.%s)", param.Name)}
		case "string":
			pipelineSpecTaskParameter.Value = fmt.Sprintf("$(params.%s)", param.Name)
		default:
			return nil, fmt.Errorf("unsupported parameter type '%s' from parameter '%s' in Task %s",
				param.Type, param.Name, task.Metadata.Name)
		}

		pipelineTask.Parameters = append(pipelineTask.Parameters, pipelineSpecTaskParameter)
	}

	jsonOps = append(jsonOps, json6902.Operation{
		Type:  json6902.Add,
		Path:  "/spec/tasks/-",
		Value: pipelineTask,
	})

	return append(patches, kustomizeV1Alpha1.TargetPatch{
		DescriptiveComment: "Add the Task to the Tekton Pipeline.",
		Target:             &kustomizeV1Alpha1.Target{Kind: "Pipeline"},
		Patch:              mustYAMLString(jsonOps),
	}), nil
}

// Adds a Task in the pipeline that will read the scan id and scan start time of the pipeline
// instance.
func AddScanUUIDAndStartTimeToPipeline(task *tektonV1Beta1.Task) ([]kustomizeV1Alpha1.TargetPatch, error) {
	if !hasLabel(task.Metadata, "v1.dracon.ocurity.com/component", "producer") {
		return nil, nil
	}

	return []kustomizeV1Alpha1.TargetPatch{
		{
			DescriptiveComment: "",
			Target: &kustomizeV1Alpha1.Target{
				Kind: "Pipeline",
			},
			Patch: mustYAMLString([]json6902.Operation{
				{
					Type: json6902.Add,
					Path: "/spec/tasks/-",
					Value: tektonV1Beta1.PipelineSpecTask{
						Name: task.Metadata.Name,
						Parameters: []*tektonV1Beta1.PipelineSpecTaskParameter{
							{
								Name:  "dracon_scan_id",
								Value: "$(tasks.base.results.dracon-scan-id)",
							},
							{
								Name:  "dracon_scan_start_time",
								Value: "$(tasks.base.results.dracon-scan-start-time)",
							},
						},
					},
				},
			}),
		},
	}, nil
}
