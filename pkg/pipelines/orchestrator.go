package pipelines

import (
	"context"

	"github.com/go-errors/errors"
	tektonv1beta1api "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/ocurity/dracon/pkg/components"
	"github.com/ocurity/dracon/pkg/k8s"
)

// Orchestrator represents a piece of code
type Orchestrator interface {
	// Prepare checks if the expected components are present in the cluster and
	// performs any operations to ensure that the workflow can be deployed.
	Prepare(context.Context, []components.Component) error

	// Deploy will generate a Pipeline based on the components and return it.
	// If dry run is set to false, the pipeline will also be applied to the
	// cluster.
	Deploy(context.Context, *tektonv1beta1api.Pipeline, []components.Component, string, bool) (*tektonv1beta1api.Pipeline, error)
}

func NewOrchestrator(clientset k8s.ClientInterface, namespace string) Orchestrator {
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

func (k k8sOrchestrator) Prepare(ctx context.Context, pipelineComponents []components.Component) error {
	var err error

	helmManagedComponents, err := k.getRepositoryTasks(ctx)
	if err != nil {
		return errors.Errorf("could not scan namespace for Helm managed components: %w", err)
	}

	for i, pipelineComponent := range pipelineComponents {
		if pipelineComponent.OrchestrationType == components.Naive {
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
			pipelineComponents[i] = pipelineComponent
		}
	}

	return nil
}

func (k k8sOrchestrator) Deploy(ctx context.Context, basePipeline *tektonv1beta1api.Pipeline, pipelineComponents []components.Component, suffix string, dryRun bool) (*tektonv1beta1api.Pipeline, error) {
	taskList := []*tektonv1beta1api.Task{}
	for _, pipelineComponent := range pipelineComponents {
		taskList = append(taskList, pipelineComponent.Manifest.(*tektonv1beta1api.Task))
	}

	tektonBackend, err := NewTektonV1Beta1Backend(basePipeline, taskList, suffix)
	if err != nil {
		return nil, err
	}

	pipeline, err := tektonBackend.Generate()
	if err != nil && dryRun {
		return pipeline, err
	}

	return pipeline, k.clientset.Apply(ctx, pipeline, k.namespace, false)
}
