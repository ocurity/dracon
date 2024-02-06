package patches

import (
	"log"
	"strings"

	"github.com/ocurity/dracon/pkg/types/kubernetes"
	"gopkg.in/yaml.v3"
)

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
