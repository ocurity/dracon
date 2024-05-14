package components

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestComponentTypeMarshaling(t *testing.T) {
	var testVal struct {
		Type ComponentType
	}

	testVal.Type = Base
	marshaled, err := json.Marshal(testVal)
	require.NoError(t, err)
	require.Equal(t, `{"Type":"base"}`, string(marshaled))

	var unmarshaledTestVal struct {
		Type ComponentType
	}

	require.NoError(t, json.Unmarshal(marshaled, &unmarshaledTestVal))
	require.Equal(t, Base, unmarshaledTestVal.Type)
}
