package kustomize

import (
	wordwrap "github.com/mitchellh/go-wordwrap"
	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewComponent returns a new representation of a component with the GVK set.
func NewComponent() Component {
	return Component{
		GroupVersionKind: &metav1.GroupVersionKind{
			Group:   "kustomize.config.k8s.io",
			Version: "v1alpha1",
			Kind:    "Component",
		},
		Resources: []string{},
		Patches:   []TargetPatch{},
	}
}

// Component represents a Kustomize Component configuration.
type Component struct {
	*metav1.GroupVersionKind `yaml:",inline"`
	Resources                []string      `yaml:"resources,omitempty"`
	Patches                  []TargetPatch `yaml:"patches,omitempty"`
}

// TargetPatch represents a patch multiple targets Patch.
type TargetPatch struct {
	DescriptiveComment string  `yaml:"-"`
	Target             *Target `yaml:"target,omitempty"`
	Patch              string  `yaml:"patch,omitempty"`
}

// MarshalYAML implements custom marshalling of TargetPatch which adds the
// DescriptiveComment as a Comment.
func (tp *TargetPatch) MarshalYAML() (interface{}, error) {
	node := &yaml.Node{}
	if err := node.Encode(map[string]interface{}{
		"target": tp.Target,
		"patch":  tp.Patch,
	}); err != nil {
		return nil, err
	}
	node.HeadComment = wordwrap.WrapString(tp.DescriptiveComment, 80)

	return node, nil
}

// Target represents a Target in a patch multiple targets Patch.
type Target struct {
	Group              string `yaml:"group,omitempty"`
	Version            string `yaml:"version,omitempty"`
	Kind               string `yaml:"kind,omitempty"`
	Name               string `yaml:"name,omitempty"`
	Namespace          string `yaml:"namespace,omitempty"`
	LabelSelector      string `yaml:"labelSelector,omitempty"`
	AnnotationSelector string `yaml:"annotationSelector,omitempty"`
}
