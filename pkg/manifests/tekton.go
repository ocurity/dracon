package manifests

import (
	"context"
	"fmt"

	tektonV1Beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
)

func LoadTektonV1Beta1Task(ctx context.Context, root, pathOrURI string) (*tektonV1Beta1.Task, error) {
	obj, gKV, err := LoadK8sManifest(ctx, root, pathOrURI, "task.yaml")
	if err != nil {
		return nil, err
	}

	task, isATask := obj.(*tektonV1Beta1.Task)
	if !isATask {
		return nil, fmt.Errorf("object loaded is not a task: %v", gKV)
	}
	return task, nil
}

func LoadTektonV1Beta1Pipeline(ctx context.Context, root, pathOrURI string) (*tektonV1Beta1.Pipeline, error) {
	obj, gKV, err := LoadK8sManifest(ctx, root, pathOrURI, "pipeline.yaml")
	if err != nil {
		return nil, err
	}

	pipeline, isAPipeline := obj.(*tektonV1Beta1.Pipeline)
	if !isAPipeline {
		return nil, fmt.Errorf("object loaded is not a pipeline: %v", gKV)
	}
	return pipeline, nil
}
