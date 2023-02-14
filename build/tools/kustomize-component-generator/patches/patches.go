package patches

import (
	"log"
	"strings"

	"github.com/ocurity/dracon/build/tools/kustomize-component-generator/types/kubernetes"
	kustomize "github.com/ocurity/dracon/build/tools/kustomize-component-generator/types/kustomize.config.k8s.io/v1alpha1"
	"gopkg.in/yaml.v3"
)

// UnusedValue is a representation of an unused value by Kustomize.
const UnusedValue = "unused"

// Patch abstracts the different Kustomize Patches that we want to generat.
type Patch interface {
	// GeneratePatch returns a Kustomize TargetPatch.
	GeneratePatch() *kustomize.TargetPatch
}

func mustYAMLString(in interface{}) string {
	var sb strings.Builder
	enc := yaml.NewEncoder(&sb)
	enc.SetIndent(2)
	if err := enc.Encode(in); err != nil {
		log.Fatalf("could not encode to yaml: %s", err)
	}

	return sb.String()
}

func hasLabel(metadata *kubernetes.Metadata, labelKey string, labelValue string) bool {
	for k, v := range metadata.Labels {
		if k == labelKey && v == labelValue {
			return true
		}
	}

	return false
}
