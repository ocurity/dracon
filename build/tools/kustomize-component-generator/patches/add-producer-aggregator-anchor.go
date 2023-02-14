package patches

import (
	"fmt"

	kustomize "github.com/ocurity/dracon/build/tools/kustomize-component-generator/types/kustomize.config.k8s.io/v1alpha1"
	tekton "github.com/ocurity/dracon/build/tools/kustomize-component-generator/types/tekton.dev/v1beta1"
)

// AddProducerAggregatorAnchor implements Patch for making the
// producer-aggregator depend on this producer Task if the producer-aggregator
// exists.
type AddProducerAggregatorAnchor struct {
	task *tekton.Task
}

// NewAddProducerAggregatorAnchor returns a new implementation of Patch for
// making the producer-aggregator depend on this producer Task if the
// producer-aggregator exists.
func NewAddProducerAggregatorAnchor(task *tekton.Task) *AddProducerAggregatorAnchor {
	return &AddProducerAggregatorAnchor{
		task: task,
	}
}

// GeneratePatch implements Patch.GeneratePatch.
func (p *AddProducerAggregatorAnchor) GeneratePatch() *kustomize.TargetPatch {
	if !hasLabel(p.task.Metadata, "v1.dracon.ocurity.com/component", "producer") {
		return nil
	}

	pipeline := tekton.NewPipeline()
	pipeline.Metadata.Name = UnusedValue

	pipeline.Spec.Tasks = append(
		pipeline.Spec.Tasks,
		&tekton.PipelineSpecTask{
			Name: "producer-aggregator",
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
		DescriptiveComment: "If we have a producer-aggregator task in the pipeline (added by the producer-aggregator component), make it depend on the completion of this producer.",
		Target: &kustomize.Target{
			Kind:               "Pipeline",
			AnnotationSelector: `v1.dracon.ocurity.com/has-producer-aggregator=true`,
		},
		Patch: mustYAMLString(pipeline),
	}
}
