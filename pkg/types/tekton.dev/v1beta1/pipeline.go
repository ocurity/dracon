package tekton

import (
	"github.com/ocurity/dracon/pkg/types/kubernetes"
)

// Pipeline represents a Tekton Pipeline configuration.
type Pipeline struct {
	*kubernetes.GVK `yaml:",inline"`
	Metadata        *kubernetes.Metadata `yaml:"metadata,omitempty"`
	Spec            *PipelineSpec        `yaml:"spec,omitempty"`
}

// NewPipeline returns a Tekton Pipeline configuration with the GVK set.
func NewPipeline() *Pipeline {
	return &Pipeline{
		GVK: &kubernetes.GVK{
			APIVersion: "tekton.dev/v1beta1",
			Kind:       "Pipeline",
		},
		Metadata: &kubernetes.Metadata{},
		Spec: &PipelineSpec{
			Tasks:      []*PipelineSpecTask{},
			Workspaces: []*PipelineSpecWorkspace{},
			Parameters: []*PipelineSpecParameter{},
		},
	}
}

// PipelineSpec represents the spec configuration of a Pipeline.
type PipelineSpec struct {
	Workspaces []*PipelineSpecWorkspace `yaml:"workspaces,omitempty"`
	Tasks      []*PipelineSpecTask      `yaml:"tasks,omitempty"`
	Parameters []*PipelineSpecParameter `yaml:"params,omitempty"`
}

// PipelineSpecTask represents the task configuration of a Pipeline Spec.
type PipelineSpecTask struct {
	Name       string                       `yaml:"name,omitempty"`
	Retries    int                          `yaml:"retries,omitempty"`
	TaskRef    *PipelineSpecTaskTaskRef     `yaml:"taskRef,omitempty"`
	Workspaces []*PipelineSpecTaskWorkspace `yaml:"workspaces,omitempty"`
	Parameters []*PipelineSpecTaskParameter `yaml:"params,omitempty"`
}

// PipelineSpecTaskTaskRef represents the task ref configuration of a
// PipelineSpecTask.
type PipelineSpecTaskTaskRef struct {
	Name string `yaml:"name,omitempty"`
}

// PipelineSpecTaskWorkspace represents the workspace configuration of a
// PipelineSpectTask.
type PipelineSpecTaskWorkspace struct {
	Name      string `yaml:"name,omitempty"`
	Workspace string `yaml:"workspace,omitempty"`
}

// PipelineSpecWorkspace represents the workspace configuration of a
// PipelineSpec.
type PipelineSpecWorkspace struct {
	Name string `yaml:"name,omitempty"`
}

// PipelineSpecTaskParameter represents the parameter configuration of a
// PipelineSpecTask.
type PipelineSpecTaskParameter struct {
	Name  string      `yaml:"name,omitempty"`
	Value interface{} `yaml:"value,omitempty"`
}

// PipelineSpecParameter represents the parameter configuration of a
// PipelineSpec.
type PipelineSpecParameter struct {
	Name        string      `jaml:"name,omitempty"`
	Description string      `yaml:"description,omitempty"`
	Type        string      `yaml:"type,omitempty"`
	Default     interface{} `yaml:"default,omitempty"`
}
