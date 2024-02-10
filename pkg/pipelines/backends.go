package pipelines

// Backend is an interface abstracting away the generator that will be used to generate manifests
// for an execution environment.
type Backend[T any] interface {
	Generate() (T, error)
}
