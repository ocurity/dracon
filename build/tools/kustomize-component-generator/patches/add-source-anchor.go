package patches

import (
	kustomize "github.com/ocurity/dracon/build/tools/kustomize-component-generator/types/kustomize.config.k8s.io/v1alpha1"
	tekton "github.com/ocurity/dracon/build/tools/kustomize-component-generator/types/tekton.dev/v1beta1"
)

// AddProducerDependencyOnSource implements Patch for making the producer Task
// depend on the source task if it exists.
type AddProducerDependencyOnSource struct {
	task *tekton.Task
}

// NewAddProducerDependencyOnSource returns a new implementation of Patch for
// making the producer Task depend on the source task if it exists.
func NewAddProducerDependencyOnSource(task *tekton.Task) *AddProducerDependencyOnSource {
	return &AddProducerDependencyOnSource{
		task: task,
	}
}

// GeneratePatch implements Patch.GeneratePatch.
func (p *AddProducerDependencyOnSource) GeneratePatch() *kustomize.TargetPatch {
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
					Value: []string{"$(tasks.source.results.anchor)"},
				},
			},
		},
	)

	return &kustomize.TargetPatch{
		DescriptiveComment: "If we have a `source` task in the pipeline (added by a `source` component), depend on the completion of that source by referencing its anchor.",
		Target: &kustomize.Target{
			Kind:               "Pipeline",
			AnnotationSelector: `v1.dracon.ocurity.com/has-source=true`,
		},
		Patch: mustYAMLString(pipeline),
	}
}
