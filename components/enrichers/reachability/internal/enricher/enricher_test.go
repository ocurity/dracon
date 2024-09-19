package enricher_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ocurity/dracon/components/enrichers/reachability/internal/atom"
	"github.com/ocurity/dracon/components/enrichers/reachability/internal/atom/purl"
	"github.com/ocurity/dracon/components/enrichers/reachability/internal/conf"
	"github.com/ocurity/dracon/components/enrichers/reachability/internal/enricher"
	"github.com/ocurity/dracon/components/enrichers/reachability/internal/fs"
)

func TestNewEnricher(t *testing.T) {
	cfg := &conf.Conf{}

	purlParser, err := purl.NewParser()
	require.NoError(t, err)

	r, err := atom.NewReader("/some/path", purlParser)
	require.NoError(t, err)

	rw, err := fs.NewReadWriter("/some/read/path", "some/write/path")
	require.NoError(t, err)

	for _, tt := range []struct {
		testCase   string
		cfg        *conf.Conf
		atomReader *atom.Reader
		fsRW       *fs.ReadWriter
		expectsErr bool
	}{
		{
			testCase:   "it returns an error because the supplied configuration is nil",
			cfg:        nil,
			atomReader: r,
			fsRW:       rw,
			expectsErr: true,
		},
		{
			testCase:   "it returns an error because the supplied atom reader is nil",
			cfg:        cfg,
			atomReader: nil,
			fsRW:       rw,
			expectsErr: true,
		},
		{
			testCase:   "it returns an error because the supplied read/writer is nil",
			cfg:        cfg,
			atomReader: r,
			fsRW:       nil,
			expectsErr: true,
		},
		{
			testCase:   "it returns a new enricher",
			cfg:        cfg,
			atomReader: r,
			fsRW:       rw,
			expectsErr: false,
		},
	} {
		t.Run(tt.testCase, func(t *testing.T) {
			t.Parallel()

			enr, err := enricher.NewEnricher(tt.cfg, tt.atomReader, tt.fsRW)
			if tt.expectsErr {
				assert.Error(t, err)
				assert.Nil(t, enr)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, enr)
			}
		})
	}
}
