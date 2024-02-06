package json6902

type OperationType string

const (
	Add     OperationType = "add"
	Replace               = "replace"
)

type Operation struct {
	Path  string
	Type  OperationType `json:"op" yaml:"op"`
	Value interface{}
}
