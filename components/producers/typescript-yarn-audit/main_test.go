package main

import (
	"testing"

	"github.com/ocurity/dracon/components/producers"
	"github.com/stretchr/testify/assert"
)

func TestEndToEndCLIWithJSON(t *testing.T) {
	err := producers.TestEndToEnd(t, "./examples/result.json", "./examples/result.pb")
	assert.NoError(t, err)
}
