package patches

import (
	"fmt"
	"log"
	"strings"

	kustomize "github.com/ocurity/dracon/pkg/components/generator/types/kustomize.config.k8s.io/v1alpha1"
	tekton "github.com/ocurity/dracon/pkg/components/generator/types/tekton.dev/v1beta1"
)

// AddTaskToPipeline implements Patch for adding the Task to the Tekton
// Pipeline.
type AddTaskToPipeline struct {
	task *tekton.Task
}

// NewAddTaskToPipeline returns a new implementation of Patch for adding the
// Task to the Tekton Pipeline.
func NewAddTaskToPipeline(task *tekton.Task) *AddTaskToPipeline {
	return &AddTaskToPipeline{
		task: task,
	}
}

// GeneratePatch implements Patch.GeneratePatch.
func (p *AddTaskToPipeline) GeneratePatch() *kustomize.TargetPatch {
	pipeline := tekton.NewPipeline()
	pipeline.Metadata.Name = UnusedValue

	p.addTask(pipeline)
	p.addWorkspaces(pipeline)
	p.addParameters(pipeline)

	return &kustomize.TargetPatch{
		DescriptiveComment: "Add the Task to the Tekton Pipeline.",
		Target:             &kustomize.Target{Kind: "Pipeline"},
		Patch:              mustYAMLString(pipeline),
	}
}

func (p *AddTaskToPipeline) addTask(pipeline *tekton.Pipeline) {
	pipeline.Spec.Tasks = append(pipeline.Spec.Tasks, &tekton.PipelineSpecTask{
		Name: p.task.Metadata.Name,
		TaskRef: &tekton.PipelineSpecTaskTaskRef{
			Name: p.task.Metadata.Name,
		},
	})
}

func (p *AddTaskToPipeline) addWorkspaces(pipeline *tekton.Pipeline) {
	if len(p.task.Spec.Workspaces) < 1 {
		return
	}

	for _, ws := range p.task.Spec.Workspaces {
		pipeline.Spec.Workspaces = append(
			pipeline.Spec.Workspaces,
			&tekton.PipelineSpecWorkspace{
				Name: ws.Name, // defines the name of the workspace in the pipeline.
			},
		)

		for _, t := range pipeline.Spec.Tasks {
			t.Workspaces = append(
				t.Workspaces,
				&tekton.PipelineSpecTaskWorkspace{
					Name:      ws.Name, // references the name of the workspace in the task.
					Workspace: ws.Name, // references the name of the workspace in the pipeline.
				},
			)
		}
	}
}

func (p *AddTaskToPipeline) addParameters(pipeline *tekton.Pipeline) {
	if len(p.task.Spec.Parameters) < 1 {
		return
	}

	for _, param := range p.task.Spec.Parameters {
		// ensure that parameter name starts with the same name as the Task
		if !strings.HasPrefix(param.Name, p.task.Metadata.Name) {
			log.Fatalf("Parameter '%s' in '%s/%s/%s' should start with '%s'.",
				param.Name,
				p.task.APIVersion,
				p.task.Kind,
				p.task.Metadata.Name,
				p.task.Metadata.Name,
			)
		}
		pipeline.Spec.Parameters = append(
			pipeline.Spec.Parameters,
			&tekton.PipelineSpecParameter{
				Name:        param.Name,
				Type:        param.Type,
				Description: param.Description,
				Default:     param.Default,
			},
		)

		pipelineSpecTaskParameter := &tekton.PipelineSpecTaskParameter{
			Name: param.Name,
		}
		switch param.Type {
		case "array":
			pipelineSpecTaskParameter.Value = []string{fmt.Sprintf("$(params.%s)", param.Name)}
		case "string":
			pipelineSpecTaskParameter.Value = fmt.Sprintf("$(params.%s)", param.Name)
		default:
			panic(fmt.Sprintf("unsupported parameter type '%s'", param.Type))
		}

		for _, t := range pipeline.Spec.Tasks {
			t.Parameters = append(
				t.Parameters,
				pipelineSpecTaskParameter,
			)
		}
	}
}
