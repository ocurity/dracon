package patches

import (
	"fmt"

	kustomize "github.com/ocurity/dracon/pkg/components/generator/types/kustomize.config.k8s.io/v1alpha1"
	tekton "github.com/ocurity/dracon/pkg/components/generator/types/tekton.dev/v1beta1"
)

// AddEnricherAggregatorAnchor implements Patch for making the
// enricher-aggregator depend on this enricher Task if the enricher-aggregator
// exists.
type AddEnricherAggregatorAnchor struct {
	task *tekton.Task
}

// NewAddEnricherAggregatorAnchor returns a new implementation of Patch for
// making the enricher-aggregator depend on this enricher Task if the
// enricher-aggregator exists.
func NewAddEnricherAggregatorAnchor(task *tekton.Task) *AddEnricherAggregatorAnchor {
	return &AddEnricherAggregatorAnchor{
		task: task,
	}
}

// GeneratePatch implements Patch.GeneratePatch.
func (p *AddEnricherAggregatorAnchor) GeneratePatch() *kustomize.TargetPatch {
	if !hasLabel(p.task.Metadata, "v1.dracon.ocurity.com/component", "enricher") {
		return nil
	}

	pipeline := tekton.NewPipeline()
	pipeline.Metadata.Name = UnusedValue

	pipeline.Spec.Tasks = append(
		pipeline.Spec.Tasks,
		&tekton.PipelineSpecTask{
			Name: "enricher-aggregator",
			Parameters: []*tekton.PipelineSpecTaskParameter{
				{
					Name: "anchors",
					Value: []string{
						fmt.Sprintf(
							"$(tasks.%s.results.anchor)",
							p.task.Metadata.Name,
						),
					},
				},
			},
		},
	)

	return &kustomize.TargetPatch{
		DescriptiveComment: "If we have a enricher-aggregator task in the pipeline (added by the enricher-aggregator component), make it depend on the completion of this enricher.",
		Target: &kustomize.Target{
			Kind:               "Pipeline",
			AnnotationSelector: `v1.dracon.ocurity.com/has-enricher-aggregator=true`,
		},
		Patch: mustYAMLString(pipeline),
	}
}
