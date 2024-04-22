package pipelines

import corev1 "k8s.io/api/core/v1"

type GenerationOpts struct {
	ImagePullPolicy *corev1.PullPolicy
}

// Backend is an interface abstracting away the generator that will be used to generate manifests
// for an execution environment.
type Backend[T any] interface {
	Generate(GenerationOpts) (T, error)
}
