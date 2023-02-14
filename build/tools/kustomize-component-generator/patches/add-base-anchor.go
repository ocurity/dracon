package patches

import (
	kustomize "github.com/ocurity/dracon/build/tools/kustomize-component-generator/types/kustomize.config.k8s.io/v1alpha1"
	tekton "github.com/ocurity/dracon/build/tools/kustomize-component-generator/types/tekton.dev/v1beta1"
)

// AddProducerDependencyOnBase implements Patch for making the producer Task
// depend on the source task if it exists.
type AddProducerDependencyOnBase struct {
	task *tekton.Task
}

// NewAddProducerDependencyOnBase returns a new implementation of Patch for
// making the producer Task depend on the source task if it exists.
func NewAddProducerDependencyOnBase(task *tekton.Task) *AddProducerDependencyOnBase {
	return &AddProducerDependencyOnBase{
		task: task,
	}
}

// GeneratePatch implements Patch.GeneratePatch.
func (p *AddProducerDependencyOnBase) GeneratePatch() *kustomize.TargetPatch {
	if !hasLabel(p.task.Metadata, "v1.dracon.ocurity.com/component", "producer") {
		return nil
	}

	pipeline := tekton.NewPipeline()
	pipeline.Metadata.Name = UnusedValue

	pipeline.Spec.Tasks = append(
		pipeline.Spec.Tasks,
		&tekton.PipelineSpecTask{
			Name: p.task.Metadata.Name,
			Parameters: []*tekton.PipelineSpecTaskParameter{
				{
					Name:  "anchors",
					Value: []string{"$(tasks.base.results.anchor)"},
				},
			},
		},
	)

	return &kustomize.TargetPatch{
		Target: &kustomize.Target{
			Kind: "Pipeline",
		},
		Patch: mustYAMLString(pipeline),
	}
}
