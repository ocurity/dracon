package kubernetes

// GVK represents a Kubenetes resource's apiVersion and kind.
type GVK struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
}
