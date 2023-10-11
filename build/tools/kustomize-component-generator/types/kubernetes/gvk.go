// Package kubernetes holds a representation of Kubenetes resource's apiVersion and kind.
package kubernetes

// GVK represents a Kubenetes resource's apiVersion and kind.
type GVK struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
}
