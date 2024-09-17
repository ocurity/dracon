package purl_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ocurity/dracon/components/enrichers/reachability/internal/atom/purl"
)

func TestNewParser(t *testing.T) {
	t.Run("should return new parser with valid matchers", func(t *testing.T) {
		p, err := purl.NewParser()
		require.NoError(t, err)
		require.NotNil(t, p)
	})
}
