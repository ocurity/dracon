package components

import (
	"context"
	"strings"

	"github.com/go-errors/errors"
	"github.com/package-url/packageurl-go"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/ocurity/dracon/pkg/manifests"
)

type OrchestrationType int

const (
	// Naive means that a component is deployed on the cluster automatically
	// before a Pipeline is deployed without checking if it's a newer or older
	// version.
	Naive OrchestrationType = iota
	// ExternalHelm means that a component is deployed on the cluster using
	// Helm but the orchestrator itself is not involved in this process.
	ExternalHelm
)

// Component represents a Dracon component. At the moment it can only be a
// Tekton Task, but in the future it might represent other things too.
type Component struct {
	// Name of the component. Should be unique
	Name string

	// Reference is the original reference to the component in the pipeline
	// description
	Reference string

	// Repository is the repository from where this component was fetched
	Repository string

	// Type of the component (base, source, producer, etc...)
	Type ComponentType

	// Resolved shows whether or not the component manifest has been loaded
	// or not. Before this is set to true, the Type and Manifest are not known.
	Resolved bool

	// OrchestrationType shows how the component deployment is managed
	OrchestrationType OrchestrationType

	// Manifest is the K8s manifest of the object
	Manifest runtime.Object
}

// FromReference resolves a reference and returns an instance of a component.
// If the component is naively orchestrated, it will be loaded
func FromReference(ctx context.Context, ref string) (Component, error) {
	zero := Component{}

	if strings.HasPrefix(ref, "pkg:helm/") {
		componentPURL, err := packageurl.FromString(ref)
		if err != nil {
			return zero, errors.Errorf("%s: reference looks like a PURL but is invalid: %w", ref, err)
		}

		return Component{
			Name:              componentPURL.Name,
			Repository:        componentPURL.Namespace,
			Reference:         ref,
			OrchestrationType: ExternalHelm,
		}, nil
	}

	task, err := manifests.LoadTektonV1Beta1Task(ctx, ".", ref)
	if err != nil {
		return zero, errors.Errorf("could not load reference: %w", err)
	}

	componentType, err := ValidateTask(task)
	if err != nil {
		return zero, err
	}

	return Component{
		Name:              task.Name,
		Type:              componentType,
		Reference:         ref,
		Resolved:          true,
		OrchestrationType: Naive,
		Manifest:          task,
	}, nil
}
