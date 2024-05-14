package pipelines

import (
	tektonv1api "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	tektonv1beta1api "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	// BasePipeline used to build all pipelines
	BasePipeline = &tektonv1beta1api.Pipeline{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pipeline",
			APIVersion: tektonv1beta1api.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{Name: "dracon"},
		Spec: tektonv1beta1api.PipelineSpec{
			Description: "Base pipeline for all Dracon pipelines",
		},
	}

	// BaseTask used to inject tags and timestamps to a pipeline
	BaseTask = &tektonv1beta1api.Task{
		ObjectMeta: metav1.ObjectMeta{
			Name: "base",
			Labels: map[string]string{
				"v1.dracon.ocurity.com/component": "base",
			},
		},
		Spec: tektonv1beta1api.TaskSpec{
			Params: tektonv1beta1api.ParamSpecs{
				tektonv1beta1api.ParamSpec{
					Name: "base-scan-tags",
					Type: "string",
					Default: &tektonv1beta1api.ParamValue{
						Type:      tektonv1beta1api.ParamTypeString,
						StringVal: "",
					},
				},
			},
			Steps: []tektonv1beta1api.Step{
				{
					Name:   "generate-scan-id-start-time",
					Image:  "docker.io/busybox:1.35.0",
					Script: "cat /proc/sys/kernel/random/uuid | tee $(results.dracon-scan-id.path)\ndate +\"%Y-%m-%dT%H:%M:%SZ\" | tee $(results.dracon-scan-start-time.path)\necho \"$(params.base-scan-tags)\" | tee $(results.dracon-scan-tags.path)\n",
					Results: []tektonv1api.StepResult{
						{
							Name:        "dracon-scan-start-time",
							Description: "The scan start time",
						},
						{
							Name:        "dracon-scan-id",
							Description: "The scan unique id",
						},
						{
							Name:        "dracon-scan-tags",
							Description: "serialized map[string]string of tags for this scan",
						},
					},
				},
			},
		},
	}
)
