package components

import "fmt"

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
