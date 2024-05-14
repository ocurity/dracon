package components

import (
	"strings"

	"github.com/go-errors/errors"
	tektonV1Beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
)

func ValidateTask(task *tektonV1Beta1.Task) (ComponentType, error) {
	componentTypeLabel, exists := task.Labels[LabelKey]
	if !exists {
		return Base, errors.Errorf("%s: task does not have a component type label", task.Name)
	}

	componentType, err := ToComponentType(componentTypeLabel)
	if err != nil {
		return Base, errors.Errorf("%s: task has wrong component type: %w", task.Name, err)
	}

	for _, param := range task.Spec.Params {
		if !strings.HasPrefix(param.Name, task.Name) {
			return Base, errors.Errorf("parameter '%s' in '%s/%s' should be prefixed with '%s'",
				param.Name, task.APIVersion, task.Kind, task.Name,
			)
		}
		if param.Type != "string" && param.Type != "array" {
			return Base, errors.Errorf("unsupported parameter type '%s' from parameter '%s' in Task %s",
				param.Type, param.Name, task.Name)
		}
	}
	return componentType, nil
}
