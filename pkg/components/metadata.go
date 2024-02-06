package components

type Type string

const (
	Source   Type = "source"
	Producer      = "producer"
	Enricher      = "enricher"
	Consumer      = "consumer"
)
