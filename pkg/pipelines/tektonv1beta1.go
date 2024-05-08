package pipelines

import (
	"context"
	"fmt"
	"slices"

	"github.com/go-errors/errors"
	tektonv1beta1api "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/ocurity/dracon/pkg/components"
	"github.com/ocurity/dracon/pkg/k8s"
)

var (
	// ErrNoComponentsInKustomization is returned when a kustomization has no
	// components listed
	ErrNoComponentsInKustomization = errors.New("no components listed in kustomization")
	// ErrNoTasks is returned when no tasks are provided to the Tekton backend
	ErrNoTasks = errors.New("no tasks provided")
	// ErrNotResolved is returned when a component that has not been resolved
	// is passed to the Orchestrator
	ErrNotResolved = errors.New("component has not been resolved")
)

// addParamsAndEnvVars will add parameters and environment variables to the producer task that will
// allow it to pick the start time, pipeline UUID and any tags that have been given as parameter to
// the pipeline so that the issues discovered can be annotated with these values.
func addParamsAndEnvVars(pipelineTask *tektonv1beta1api.PipelineTask, anchors map[string][]string, task *tektonv1beta1api.Task) {
	pipelineTask.Params = append(pipelineTask.Params, []tektonv1beta1api.Param{
		{
			Name: "dracon_scan_id",
			Value: tektonv1beta1api.ParamValue{
				Type:      tektonv1beta1api.ParamTypeString,
				StringVal: fmt.Sprintf("$(tasks.%s.results.dracon-scan-id)", anchors[components.Base.String()][0]),
			},
		},
		{
			Name: "dracon_scan_start_time",
			Value: tektonv1beta1api.ParamValue{
				Type:      tektonv1beta1api.ParamTypeString,
				StringVal: fmt.Sprintf("$(tasks.%s.results.dracon-scan-start-time)", anchors[components.Base.String()][0]),
			},
		},
		{
			Name: "dracon_scan_tags",
			Value: tektonv1beta1api.ParamValue{
				Type:      tektonv1beta1api.ParamTypeString,
				StringVal: fmt.Sprintf("$(tasks.%s.results.dracon-scan-tags)", anchors[components.Base.String()][0]),
			},
		},
	}...)

	task.Spec.Params = append(task.Spec.Params, tektonv1beta1api.ParamSpecs{
		{
			Name: "dracon_scan_id",
			Type: tektonv1beta1api.ParamTypeString,
		},
		{
			Name: "dracon_scan_start_time",
			Type: tektonv1beta1api.ParamTypeString,
		},
		{
			Name: "dracon_scan_tags",
			Type: tektonv1beta1api.ParamTypeString,
		},
	}...)

	for i, step := range task.Spec.Steps {
		step.Env = append(step.Env, []corev1.EnvVar{
			{
				Name:  "DRACON_SCAN_TIME",
				Value: "$(params.dracon_scan_start_time)",
			},
			{
				Name:  "DRACON_SCAN_ID",
				Value: "$(params.dracon_scan_id)",
			},
			{
				Name:  "DRACON_SCAN_TAGS",
				Value: "$(params.dracon_scan_tags)",
			},
		}...)
		task.Spec.Steps[i] = step
	}
}

// NewTektonV1Beta1Orchestrator returns an Orchestrator implementation for TektonV1Beta1
func NewTektonV1Beta1Orchestrator(clientset k8s.ClientInterface, namespace string) Orchestrator[*tektonv1beta1api.Pipeline] {
	return k8sOrchestrator{
		clientset: clientset,
		namespace: namespace,
	}
}

type k8sOrchestrator struct {
	clientset k8s.ClientInterface
	namespace string
}

// getRepositoryTasks checks all the Tasks of the namespace and creates a
// registry based on their Helm release name.
func (k k8sOrchestrator) getRepositoryTasks(ctx context.Context) (map[string]map[string]tektonv1beta1api.Task, error) {
	taskList, err := k.clientset.
		Tasks(k.namespace).
		List(ctx, metav1.ListOptions{
			LabelSelector: metav1.FormatLabelSelector(
				&metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app.kubernetes.io/managed-by": "Helm",
					},
				},
			),
		})
	if err != nil {
		return nil, errors.Errorf("could not resolve tasks of Helm registry: %w", err)
	}

	helmManagedComponents := map[string]map[string]tektonv1beta1api.Task{}
	for _, task := range taskList.Items {
		// if the task is not annotated then it's not managed by Helm
		helmRegistry, exists := task.Annotations["meta.helm.sh/release-name"]
		if !exists {
			continue
		}

		tasksInHelmRelease, exists := helmManagedComponents[helmRegistry]
		if !exists {
			tasksInHelmRelease = map[string]tektonv1beta1api.Task{}
			helmManagedComponents[helmRegistry] = tasksInHelmRelease
		}
		tasksInHelmRelease[task.Name] = task
	}

	return helmManagedComponents, nil
}

// Prepare will scan the cluster for all Helm managed Tasks and will deploy all
// Tasks managed in a custom way.
func (k k8sOrchestrator) Prepare(ctx context.Context, pipelineComponents []components.Component) error {
	var err error

	helmManagedComponents, err := k.getRepositoryTasks(ctx)
	if err != nil {
		return errors.Errorf("could not scan namespace for Helm managed components: %w", err)
	}

	for i, pipelineComponent := range pipelineComponents {
		if pipelineComponent.OrchestrationType == components.Naive {
			if !pipelineComponent.Resolved || pipelineComponent.Manifest == nil {
				return ErrNotResolved
			}
			k.clientset.Apply(ctx, pipelineComponent.Manifest, k.namespace, false)
		} else if pipelineComponent.OrchestrationType == components.ExternalHelm {
			componentSet, exists := helmManagedComponents[pipelineComponent.Repository]
			if !exists {
				return errors.Errorf("no Helm release with name %s is deployed in the namespace %s (%s)", pipelineComponent.Repository, k.namespace, pipelineComponent.Name)
			}

			component, exists := componentSet[pipelineComponent.Name]
			if !exists {
				return errors.Errorf("component %s/%s could not be found in the cluster", pipelineComponent.Repository, pipelineComponent.Name)
			}

			pipelineComponent.Manifest = &component
			componentTypeLabel, exists := component.Labels[components.LabelKey]
			if !exists {
				return errors.Errorf("%s: task does not have a component type label", component.Name)
			}

			componentType, err := components.ToComponentType(componentTypeLabel)
			if err != nil {
				return errors.Errorf("%s: task has wrong component type: %w", component.Name, err)
			}
			pipelineComponent.Type = componentType
			pipelineComponent.Resolved = true
			pipelineComponents[i] = pipelineComponent
		}
	}

	return nil
}

// Deploy will generate a pipeline based on the components provided
func (k k8sOrchestrator) Deploy(ctx context.Context, basePipeline *tektonv1beta1api.Pipeline, pipelineComponents []components.Component, suffix string, dryRun bool) (*tektonv1beta1api.Pipeline, error) {
	if len(pipelineComponents) == 0 {
		return nil, errors.Errorf("%w", ErrNoTasks)
	}

	taskList := []*tektonv1beta1api.Task{}
	for _, pipelineComponent := range pipelineComponents {
		taskList = append(taskList, pipelineComponent.Manifest.(*tektonv1beta1api.Task))
	}

	for _, task := range taskList {
		// TODO(?): revisit if we need this in the future
		// fixTaskPrefixSuffix(task, prefix, suffix)
		addAnchorParameter(task)
		addAnchorResult(task)
	}

	// Sort tasks based on their component type
	slices.SortFunc(taskList, func(a *tektonv1beta1api.Task, b *tektonv1beta1api.Task) int {
		componentTypeA := components.MustGetComponentType(a.Labels[components.LabelKey])
		componentTypeB := components.MustGetComponentType(b.Labels[components.LabelKey])
		return int(componentTypeA) - int(componentTypeB)
	})

	pipeline, err := k.generatePipeline(basePipeline, taskList)
	if err != nil && dryRun {
		return pipeline, err
	}

	return pipeline, k.clientset.Apply(ctx, pipeline, k.namespace, false)
}

func (k k8sOrchestrator) generatePipeline(pipeline *tektonv1beta1api.Pipeline, taskList []*tektonv1beta1api.Task) (*tektonv1beta1api.Pipeline, error) {
	pipelineWorkspaces := map[string]struct{}{}
	anchors := map[string][]string{}

	for _, task := range taskList {
		componentType := task.Labels[components.LabelKey]
		anchors[componentType] = append(anchors[componentType], task.Name)

		// add task to pipeline tasks
		pipelineTask := tektonv1beta1api.PipelineTask{
			Name: task.Name,
			TaskRef: &tektonv1beta1api.TaskRef{
				Name: task.Name,
			},
		}

		// add task's workspaces to pipeline workspaces
		// make sure to propagate the `optional` field
		for _, ws := range task.Spec.Workspaces {
			if _, inserted := pipelineWorkspaces[ws.Name]; !inserted {
				pipeline.Spec.Workspaces = append(pipeline.Spec.Workspaces, tektonv1beta1api.PipelineWorkspaceDeclaration{
					Name:     ws.Name,
					Optional: ws.Optional,
				})
				pipelineWorkspaces[ws.Name] = struct{}{}
			}
			pipelineTask.Workspaces = append(pipelineTask.Workspaces, tektonv1beta1api.WorkspacePipelineTaskBinding{
				Name:      ws.Name,
				Workspace: ws.Name,
			})
		}

		// add the task's parameters to the pipeline's parameters and
		// reference them in the pipeline task parameters
		pipelineTask.Params = make(tektonv1beta1api.Params, len(task.Spec.Params))

		for i, param := range task.Spec.Params {
			pipelineTask.Params[i] = tektonv1beta1api.Param{
				Name:  param.Name,
				Value: tektonv1beta1api.ParamValue{},
			}

			if param.Name == "anchors" {
				anchorTargetComponentType := components.MustGetComponentType(componentType) - 1
				values := []string{}

				// get all the tasks that should be finished before this one starts
				for _, anchorTarget := range anchors[anchorTargetComponentType.String()] {
					values = append(values, fmt.Sprintf("$(tasks.%s.results.anchor)", anchorTarget))
				}

				pipelineTask.Params[i].Value.ArrayVal = values
				pipelineTask.Params[i].Value.Type = tektonv1beta1api.ParamTypeArray
			} else {
				switch param.Type {
				case tektonv1beta1api.ParamTypeArray:
					pipelineTask.Params[i].Value.Type = param.Type
					pipelineTask.Params[i].Value.ArrayVal = []string{fmt.Sprintf("$(params.%s)", param.Name)}
				case tektonv1beta1api.ParamTypeString:
					pipelineTask.Params[i].Value.Type = param.Type
					pipelineTask.Params[i].Value.StringVal = fmt.Sprintf("$(params.%s)", param.Name)
				case "":
					return nil, errors.Errorf("parameter %s of task %s has no type set", param.Name, task.Name)
				}

				// ensure that the parameter type is always set
				if param.Default != nil && param.Default.Type == "" {
					param.Default.Type = param.Type
				}

				// add parameter to pipeline parameters
				pipeline.Spec.Params = append(pipeline.Spec.Params, tektonv1beta1api.ParamSpec{
					Name:        param.Name,
					Type:        param.Type,
					Description: param.Description,
					Default:     param.Default,
				})
			}
		}

		// add scan ID and scan time to all producers
		if task.Labels[components.LabelKey] == components.Producer.String() {
			addParamsAndEnvVars(&pipelineTask, anchors, task)
		}

		// add task reference to pipeline's tasks
		pipeline.Spec.Tasks = append(pipeline.Spec.Tasks, pipelineTask)
	}

	return pipeline, nil
}
