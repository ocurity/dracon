package patches

import (
	kustomize "github.com/ocurity/dracon/pkg/components/generator/types/kustomize.config.k8s.io/v1alpha1"
	tekton "github.com/ocurity/dracon/pkg/components/generator/types/tekton.dev/v1beta1"
)

// AddScanUUIDAndStartTimeToPipeline implements Patch for adding DAG anchors to every Task.
type AddScanUUIDAndStartTimeToPipeline struct {
	task *tekton.Task
}

// NewAddScanUUIDAndStartTimeToPipeline returns a new implementation of Patch for adding DAG
// anchors to every Task.
func NewAddScanUUIDAndStartTimeToPipeline(task *tekton.Task) *AddScanUUIDAndStartTimeToPipeline {
	return &AddScanUUIDAndStartTimeToPipeline{
		task: task,
	}
}

// GeneratePatch implements Patch.GeneratePatch.
func (p *AddScanUUIDAndStartTimeToPipeline) GeneratePatch() *kustomize.TargetPatch {
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
					Name:  "dracon_scan_id",
					Value: "$(tasks.base.results.dracon-scan-id)",
				},
				{
					Name:  "dracon_scan_start_time",
					Value: "$(tasks.base.results.dracon-scan-start-time)",
				},
			},
		},
	)

	return &kustomize.TargetPatch{
		DescriptiveComment: "",
		Target: &kustomize.Target{
			Kind: "Pipeline",
		},
		Patch: mustYAMLString(pipeline),
	}
}
