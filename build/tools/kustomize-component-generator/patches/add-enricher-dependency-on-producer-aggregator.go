package patches

import (
	"fmt"

	kustomize "github.com/ocurity/dracon/build/tools/kustomize-component-generator/types/kustomize.config.k8s.io/v1alpha1"
	tekton "github.com/ocurity/dracon/build/tools/kustomize-component-generator/types/tekton.dev/v1beta1"
)

// AddEnricherDependencyOnProducerAggregator implements Patch for making
// enrichers depend on the producer-aggregator if it exists.
type AddEnricherDependencyOnProducerAggregator struct {
	task *tekton.Task
}

// NewAddEnricherDependencyOnProducerAggregator returns a new implementation of
// Patch for making enrichers depend on the producer-aggregator if it exists.
func NewAddEnricherDependencyOnProducerAggregator(task *tekton.Task) *AddEnricherDependencyOnProducerAggregator {
	return &AddEnricherDependencyOnProducerAggregator{
		task: task,
	}
}

// GeneratePatch implements Patch.GeneratePatch.
func (p *AddEnricherDependencyOnProducerAggregator) GeneratePatch() *kustomize.TargetPatch {
	if !hasLabel(p.task.Metadata, "v1.dracon.ocurity.com/component", "enricher") {
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
					Name: "anchors",
					Value: []string{
						fmt.Sprintf(
							"$(tasks.%s.results.anchor)",
							"producer-aggregator",
						),
					},
				},
			},
		},
	)

	return &kustomize.TargetPatch{
		DescriptiveComment: "If we have an producer-aggregator task in the pipeline (added by the producer-aggregator component), make the enricher depend on the completion of it.",
		Target: &kustomize.Target{
			Kind:               "Pipeline",
			AnnotationSelector: `v1.dracon.ocurity.com/has-producer-aggregator=true`,
		},
		Patch: mustYAMLString(pipeline),
	}
}
