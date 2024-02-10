package components

import "fmt"

const LabelKey string = "v1.dracon.ocurity.com/component"

type ComponentType int

const (
	Base ComponentType = iota
	Source
	Producer
	ProducerAggregator
	Enricher
	EnricherAggregator
	Consumer
)

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

func MustGetComponentType(cts string) ComponentType {
	ct, err := ToComponentType(cts)
	if err != nil {
		panic(err)
	}
	return ct
}
