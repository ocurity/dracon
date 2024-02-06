package patches

import (
	"fmt"

	"github.com/ocurity/dracon/pkg/types/json6902"
	kustomizeV1Alpha1 "github.com/ocurity/dracon/pkg/types/kustomize.config.k8s.io/v1alpha1"
	tektonV1Beta1 "github.com/ocurity/dracon/pkg/types/tekton.dev/v1beta1"
)

func generateAnchorOp(anchorID string) json6902.Operation {
	return json6902.Operation{
		Type: json6902.Add,
		Path: "/spec/tasks/-/params/-",
		Value: tektonV1Beta1.PipelineSpecTaskParameter{
			Name:  "anchors",
			Value: []string{fmt.Sprintf("$(tasks.%s.results.anchor)", anchorID)},
		},
	}
}

// AddProducerDependencyOnSource generates a kustomize patch that will add an anchor to the
// pipeline tasks to ensure that the producers will wait for the source task to complete its work
// before starting their own work.
func AddProducerDependencyOnSource(task *tektonV1Beta1.Task) ([]kustomizeV1Alpha1.TargetPatch, error) {
	if !hasLabel(task.Metadata, "v1.dracon.ocurity.com/component", "producer") {
		return nil, nil
	}

	return []kustomizeV1Alpha1.TargetPatch{
		{
			DescriptiveComment: "If we have a `source` task in the pipeline (added by a `source` component), depend on the completion of that source by referencing its anchor.",
			Target: &kustomizeV1Alpha1.Target{
				Kind:               "Pipeline",
				AnnotationSelector: `v1.dracon.ocurity.com/has-source=true`,
			},
			Patch: mustYAMLString([]json6902.Operation{generateAnchorOp("source")}),
		},
	}, nil
}

// AddProducerDependencyOnBase generates a kustomize patch that adds an anchor to the pipeline
// tasks to ensure the producers wait for the base task to complete its work before starting their
// own.
func AddProducerDependencyOnBase(task *tektonV1Beta1.Task) ([]kustomizeV1Alpha1.TargetPatch, error) {
	if !hasLabel(task.Metadata, "v1.dracon.ocurity.com/component", "producer") {
		return nil, nil
	}

	return []kustomizeV1Alpha1.TargetPatch{
		{
			DescriptiveComment: "If we have a `source` task in the pipeline (added by a `source` component), depend on the completion of that source by referencing its anchor.",
			Target: &kustomizeV1Alpha1.Target{
				Kind: "Pipeline",
			},
			Patch: mustYAMLString([]json6902.Operation{generateAnchorOp("base")}),
		},
	}, nil
}

// AddEnricherAggregatorAnchor generates a kustomize patch that adds an anchor task to the
// pipeline to ensure that the next step in the pipeline will wait for the enricher to finish
// its job before starting.
func AddEnricherAggregatorAnchor(task *tektonV1Beta1.Task) ([]kustomizeV1Alpha1.TargetPatch, error) {
	if !hasLabel(task.Metadata, "v1.dracon.ocurity.com/component", "enricher") {
		return nil, nil
	}

	return []kustomizeV1Alpha1.TargetPatch{
		{
			DescriptiveComment: "If we have an enricher-aggregator task in the pipeline (added by the enricher-aggregator component), make it depend on the completion of this enricher.",
			Target: &kustomizeV1Alpha1.Target{
				Kind:               "Pipeline",
				AnnotationSelector: `v1.dracon.ocurity.com/has-enricher-aggregator=true`,
			},
			Patch: mustYAMLString([]json6902.Operation{generateAnchorOp(task.Metadata.Name)}),
		},
	}, nil
}

// AddConsumerDependencyOnEnricherAggregator generates a kustomize patch that adds an enchor to the
// results of the enricher aggregator to ensure the consumer starts its work only once the
// enricher aggregator has finished its work.
func AddConsumerDependencyOnEnricherAggregator(task *tektonV1Beta1.Task) ([]kustomizeV1Alpha1.TargetPatch, error) {
	if !hasLabel(task.Metadata, "v1.dracon.ocurity.com/component", "consumer") {
		return nil, nil
	}

	return []kustomizeV1Alpha1.TargetPatch{
		{
			DescriptiveComment: "If we have an enricher-aggregator task in the pipeline (added by the enricher-aggregator component), make the consumer depend on the completion of it.",
			Target: &kustomizeV1Alpha1.Target{
				Kind:               "Pipeline",
				AnnotationSelector: `v1.dracon.ocurity.com/has-enricher-aggregator=true`,
			},
			Patch: mustYAMLString([]json6902.Operation{generateAnchorOp("enricher-aggregator")}),
		},
	}, nil
}

// GeneratePatch implements Patch.GeneratePatch.
func AddProducerAggregatorAnchor(task *tektonV1Beta1.Task) ([]kustomizeV1Alpha1.TargetPatch, error) {
	if !hasLabel(task.Metadata, "v1.dracon.ocurity.com/component", "producer") {
		return nil, nil
	}

	return []kustomizeV1Alpha1.TargetPatch{
		{
			DescriptiveComment: "If we have an enricher-aggregator task in the pipeline (added by the enricher-aggregator component), make the consumer depend on the completion of it.",
			Target: &kustomizeV1Alpha1.Target{
				Kind:               "Pipeline",
				AnnotationSelector: `v1.dracon.ocurity.com/has-producer-aggregator=true`,
			},
			Patch: mustYAMLString([]json6902.Operation{generateAnchorOp("enricher-aggregator")}),
		},
	}, nil
}

func AddEnricherDependencyOnProducerAggregator(task *tektonV1Beta1.Task) ([]kustomizeV1Alpha1.TargetPatch, error) {
	if !hasLabel(task.Metadata, "v1.dracon.ocurity.com/component", "enricher") {
		return nil, nil
	}

	return []kustomizeV1Alpha1.TargetPatch{
		{
			DescriptiveComment: "If we have an enricher-aggregator task in the pipeline (added by the enricher-aggregator component), make the consumer depend on the completion of it.",
			Target: &kustomizeV1Alpha1.Target{
				Kind:               "Pipeline",
				AnnotationSelector: `v1.dracon.ocurity.com/has-producer-aggregator=true`,
			},
			Patch: mustYAMLString([]json6902.Operation{generateAnchorOp("producer-aggregator")}),
		},
	}, nil
}
