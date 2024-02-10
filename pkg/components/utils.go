package components

import (
	"fmt"
	"strings"

	tektonV1Beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
)

func ValidateTask(task *tektonV1Beta1.Task) error {
	componentType, exists := task.Metadata.Labels[LabelKey]
	if !exists {
		return fmt.Errorf("%s: task does not have a component type label", task.Metadata.Name)
	}

	if _, err := ToComponentType(componentType); err != nil {
		return fmt.Errorf("%s: task has wrong component type: %w", task.Metadata.Name, err)
	}

	for _, param := range task.Spec.Params {
		if !strings.HasPrefix(param.Name, task.Name) {
			return fmt.Errorf("parameter '%s' in '%s/%s' should be prefixed with '%s'",
				param.Name, task.APIVersion, task.Kind, task.Name,
			)
		}
		if param.Type != "string" && param.Type != "array" {
			return fmt.Errorf("unsupported parameter type '%s' from parameter '%s' in Task %s",
				param.Type, param.Name, task.Metadata.Name)
		}
	}
	return nil
}
