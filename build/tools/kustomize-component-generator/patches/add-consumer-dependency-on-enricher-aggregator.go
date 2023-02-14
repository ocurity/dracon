package patches

import (
	"fmt"

	kustomize "github.com/ocurity/dracon/build/tools/kustomize-component-generator/types/kustomize.config.k8s.io/v1alpha1"
	tekton "github.com/ocurity/dracon/build/tools/kustomize-component-generator/types/tekton.dev/v1beta1"
)

// AddConsumerDependencyOnEnricherAggregator implements Patch for making
// consumers depend on the enricher-aggregator if it exists.
type AddConsumerDependencyOnEnricherAggregator struct {
	task *tekton.Task
}

// NewAddConsumerDependencyOnEnricherAggregator returns a new implementation of
// Patch for making consumers depend on the enricher-aggregator if it exists.
func NewAddConsumerDependencyOnEnricherAggregator(task *tekton.Task) *AddConsumerDependencyOnEnricherAggregator {
	return &AddConsumerDependencyOnEnricherAggregator{
		task: task,
	}
}

// GeneratePatch implements Patch.GeneratePatch.
func (p *AddConsumerDependencyOnEnricherAggregator) GeneratePatch() *kustomize.TargetPatch {
	if !hasLabel(p.task.Metadata, "v1.dracon.ocurity.com/component", "consumer") {
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
							"enricher-aggregator",
						),
					},
				},
			},
		},
	)

	return &kustomize.TargetPatch{
		DescriptiveComment: "If we have an enricher-aggregator task in the pipeline (added by the enricher-aggregator component), make the consumer depend on the completion of it.",
		Target: &kustomize.Target{
			Kind:               "Pipeline",
			AnnotationSelector: `v1.dracon.ocurity.com/has-enricher-aggregator=true`,
		},
		Patch: mustYAMLString(pipeline),
	}
}
