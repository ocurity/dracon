package component

import (
	"context"

	ocsf "github.com/ocurity/dracon/sdk/component/gen/com/github/ocsf/ocsf_schema/v1"
)

// Helpers interfaces for common functionalities.
type (
	// Validator allows validating findings by a specified criteria.
	Validator interface {
		// Validate validates the supplied findings and returns an error if invalid.
		Validate(ctx context.Context, findings []*ocsf.VulnerabilityFinding) error
	}

	// Reader allows reading findings from a source.
	Reader interface {
		// Read reads and parses findings from a source.
		Read(ctx context.Context) ([]*ocsf.VulnerabilityFinding, error)
	}

	// Filterer allows filtering out findings by a specified criteria.
	Filterer interface {
		// Filter returns filtered findings from the supplied ones.
		Filter(ctx context.Context, findings []*ocsf.VulnerabilityFinding) ([]*ocsf.VulnerabilityFinding, error)
	}
)

// Components interfaces.
type (
	// Sourcer defines the behaviour of source components.
	Sourcer interface {
		// Source runs an arbitrary step. Useful for interacting with third-party API - i.e. cloning a repository.
		Source(ctx context.Context) error
	}

	// Producer defines the behaviour of producer components.
	// Read -> Process -> Validate -> Store.
	Producer interface {
		Reader
		Validator

		// Process defines the producer behaviour.
		Process(ctx context.Context, findings []*ocsf.VulnerabilityFinding) ([]*ocsf.VulnerabilityFinding, error)
		// Store stores the findings into a destination.
		Store(ctx context.Context, findings []*ocsf.VulnerabilityFinding) error
	}

	// Enricher defines the behaviour of enricher components.
	// Read -> Filter (optional) -> Annotate -> Update.
	Enricher interface {
		Reader
		Filterer

		// Annotate enriches the findings by some criteria.
		Annotate(ctx context.Context, findings []*ocsf.VulnerabilityFinding) ([]*ocsf.VulnerabilityFinding, error)
		// Update updates the existing findings in a destination.
		Update(ctx context.Context, findings []*ocsf.VulnerabilityFinding) error
	}

	// Consumer defines the behaviour of consumer components.
	// Read -> Process.
	Consumer interface {
		Reader

		// Process processes the findings and takes some action. Useful for interacting with third-party API or
		// creating results - i.e. raising tickets.
		Process(ctx context.Context, findings []*ocsf.VulnerabilityFinding) error
	}
)
