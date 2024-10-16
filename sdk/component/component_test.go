package component_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ocurity/dracon/sdk/component"
	ocsf "github.com/ocurity/dracon/sdk/component/gen/com/github/ocsf/ocsf_schema/v1"
)

type (
	testSourcer struct{}

	testProducer struct{}

	testEnricher struct{}

	testConsumer struct{}
)

func (t testConsumer) Read(ctx context.Context) ([]*ocsf.VulnerabilityFinding, error) {
	return nil, nil
}

func (t testConsumer) Process(ctx context.Context, findings []*ocsf.VulnerabilityFinding) error {
	return nil
}

func (t testEnricher) Read(ctx context.Context) ([]*ocsf.VulnerabilityFinding, error) {
	return nil, nil
}

func (t testEnricher) Filter(ctx context.Context, findings []*ocsf.VulnerabilityFinding) ([]*ocsf.VulnerabilityFinding, error) {
	return nil, nil
}

func (t testEnricher) Annotate(ctx context.Context, findings []*ocsf.VulnerabilityFinding) ([]*ocsf.VulnerabilityFinding, error) {
	return nil, nil
}

func (t testEnricher) Update(ctx context.Context, findings []*ocsf.VulnerabilityFinding) error {
	return nil
}

func (t testProducer) Read(ctx context.Context) ([]*ocsf.VulnerabilityFinding, error) {
	return nil, nil
}

func (t testProducer) Validate(ctx context.Context, findings []*ocsf.VulnerabilityFinding) error {
	return nil
}

func (t testProducer) Process(ctx context.Context, findings []*ocsf.VulnerabilityFinding) ([]*ocsf.VulnerabilityFinding, error) {
	return nil, nil
}

func (t testProducer) Store(ctx context.Context, findings []*ocsf.VulnerabilityFinding) error {
	return nil
}

func (t testSourcer) Source(ctx context.Context) error {
	return nil
}

func TestImplementations(t *testing.T) {
	assert.Implements(t, (*component.Sourcer)(nil), testSourcer{})
	assert.Implements(t, (*component.Producer)(nil), testProducer{})
	assert.Implements(t, (*component.Enricher)(nil), testEnricher{})
	assert.Implements(t, (*component.Consumer)(nil), testConsumer{})
}
