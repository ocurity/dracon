package patches

import (
	kustomizeV1Alpha1 "github.com/ocurity/dracon/pkg/types/kustomize.config.k8s.io/v1alpha1"
	tektonV1Beta1 "github.com/ocurity/dracon/pkg/types/tekton.dev/v1beta1"
)

// UnusedValue is a representation of an unused value by Kustomize.
const UnusedValue = "unused"

// Generator is a function that will accept a Task as a parameter and will return
// a number of kustomize patches that will modify the final YAML manifests in order for the
// Task to be correctly invoked.
type Generator func(*tektonV1Beta1.Task) ([]kustomizeV1Alpha1.TargetPatch, error)
