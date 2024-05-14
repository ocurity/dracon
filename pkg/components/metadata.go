package components

import (
	"bytes"
	"fmt"

	"github.com/go-errors/errors"
)

// LabelKey is the key of the label where the value must be a string of the
// ComponentType enum
const LabelKey string = "v1.dracon.ocurity.com/component"

// ComponentType represents all the types of components that Dracon supports
type ComponentType int

const (
	// Unknown is the default component type value
	Unknown ComponentType = iota
	// Base represents the base component of a pipeline
	Base
	// Source represents the source component of a pipeline
	Source
	// Producer represents the producer component of a pipeline
	Producer
	// ProducerAggregator represents the producer aggregator component of a
	// pipeline
	ProducerAggregator
	// Enricher represents the enricher component of a pipeline
	Enricher
	// EnricherAggregator represents the enricher aggregator component of a
	// pipeline
	EnricherAggregator
	// Consumer represents the consumer component of a pipeline
	Consumer
)

// String converts the ComponentType enum value to a string
func (ct ComponentType) String() string {
	switch ct {
	case Unknown:
		return "unknown"
	case Base:
		return "base"
	case Source:
		return "source"
	case Producer:
		return "producer"
	case ProducerAggregator:
		return "producer-aggregator"
	case Enricher:
		return "enricher"
	case EnricherAggregator:
		return "enricher-aggregator"
	case Consumer:
		return "consumer"
	default:
		panic(fmt.Sprintf("unknown component type: %d", ct))
	}
}

// ToComponentType converts a string into a ComponentType enum value or
// returns an error
func ToComponentType(cts string) (ComponentType, error) {
	switch cts {
	case "unknown":
		return Unknown, nil
	case "base":
		return Base, nil
	case "source":
		return Source, nil
	case "producer":
		return Producer, nil
	case "producer-aggregator":
		return ProducerAggregator, nil
	case "enricher":
		return Enricher, nil
	case "enricher-aggregator":
		return EnricherAggregator, nil
	case "consumer":
		return Consumer, nil
	default:
		return Base, fmt.Errorf("%s: unknown component type", cts)
	}
}

// MustGetComponentType converts a string into a ComponentType enum value or
// panics
func MustGetComponentType(cts string) ComponentType {
	ct, err := ToComponentType(cts)
	if err != nil {
		panic(err)
	}
	return ct
}

// LabelValueOneOf checks if the labels map has the expected key set and if
// that key has any one of the expected values
func LabelValueOneOf(labels map[string]string, componentTypes ...ComponentType) (bool, error) {
	labelValue, hasComponentType := labels[LabelKey]
	if !hasComponentType {
		return false, errors.Errorf("no %s key in labels", LabelKey)
	}

	labelCt, err := ToComponentType(labelValue)
	if err != nil {
		return false, err
	}

	for _, componentType := range componentTypes {
		if labelCt == componentType {
			return true, nil
		}
	}

	return false, nil
}

// Equal checks if the other is a ComponentType or some type that can be parsed
// and then checks for equality.
func (ct ComponentType) Equal(other any) (bool, error) {
	switch o := other.(type) {
	case string:
		parsed, err := ToComponentType(o)
		if err != nil {
			return false, err
		}
		return parsed == ct, nil
	case ComponentType:
		return ct == o, nil
	case fmt.Stringer:
		return ct.Equal(o.String())
	default:
		return false, errors.Errorf("unknown component type")
	}
}

// MarshalJSON marshals a `ComponentType` into JSON bytes
func (ct ComponentType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + ct.String() + `"`), nil
}

// MarshalText marshals the `ComponentType` into text bytes
func (ct ComponentType) MarshalText() ([]byte, error) {
	return []byte(`"` + ct.String() + `"`), nil
}

// UnmarshalJSON unmarshalls bytes into a `ComponentType`
func (ct *ComponentType) UnmarshalJSON(b []byte) error {
	b = bytes.Trim(bytes.Trim(b, `"`), ` `)
	parsedComponent, err := ToComponentType(string(b))
	if err == nil {
		*ct = parsedComponent
	}
	return err
}

// UnmarshalText unmarshalls bytes into a `ComponentType`
func (ct *ComponentType) UnmarshalText(text []byte) error {
	text = bytes.Trim(bytes.Trim(text, `"`), ` `)
	parsedComponent, err := ToComponentType(string(text))
	if err == nil {
		*ct = parsedComponent
	}
	return err
}
