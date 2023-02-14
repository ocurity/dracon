package tekton

import (
	"github.com/ocurity/dracon/build/tools/kustomize-component-generator/types/kubernetes"
)

// Task represents a Tekton Task configuration.
type Task struct {
	*kubernetes.GVK `yaml:",inline"`

	Metadata *kubernetes.Metadata `yaml:"metadata,omitempty"`

	Spec *TaskSpec `yaml:"spec,omitempty"`
}

// TaskSpec represents the spec configuration of a Task.
type TaskSpec struct {
	Workspaces []*TaskSpecWorkspace `yaml:"workspaces,omitempty"`
	Parameters []*TaskSpecParameter `yaml:"params,omitempty"`
	Results    []*TaskSpecResult    `yaml:"results,omitempty"`
	Steps      []*TaskSpecStep      `yaml:"steps,omitempty"`
}

// TaskSpecWorkspace represents the workspace configuration of a TaskSpec.
type TaskSpecWorkspace struct {
	Name        string `yaml:"name,omitempty"`
	Description string `yaml:"description,omitempty"`
}

// TaskSpecParameter represents the parameter configuration of a TaskSpec.
type TaskSpecParameter struct {
	Name        string      `yaml:"name,omitempty"`
	Type        string      `yaml:"type,omitempty"`
	Description string      `yaml:"description,omitempty"`
	Default     interface{} `yaml:"default,omitempty"`
	Value       interface{} `yaml:"value,omitempty"`
}

// TaskSpecResult represents the result configuration of a TaskSpec.
type TaskSpecResult struct {
	Name        string `yaml:"name,omitempty"`
	Description string `yaml:"description,omitempty"`
}

// TaskSpecStep represents the step configuration of a TaskSpec.
type TaskSpecStep struct {
	Name   string             `yaml:"name,omitempty"`
	Image  string             `yaml:"image,omitempty"`
	Script string             `yaml:"script,omitempty"`
	Env    []*TaskSpecStepEnv `yaml:"env,omitempty"`
}

// TaskSpecStepEnv represents the "env" stanza of a TaskSpecStep.
type TaskSpecStepEnv struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}
