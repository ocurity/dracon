package dtrack

import "github.com/google/uuid"

type PolicyCondition struct {
	UUID     uuid.UUID
	Policy   *Policy `json:"policy,omitempty"`
	Operator string  `json:"operator"`
	Subject  string  `json:"subject"`
	Value    string  `json:"value"`
}
