package document

import (
	"time"
)

// Document represents a Dracon result (issue) object.
type Document struct {
	ScanStartTime  time.Time         `json:"scan_start_time"`
	ScanID         string            `json:"scan_id"`
	ToolName       string            `json:"tool_name"`
	Source         string            `json:"source"`
	Target         string            `json:"target"`
	Type           string            `json:"type"`
	Title          string            `json:"title"`
	SeverityText   string            `json:"severity_text"`
	CVSS           string            `json:"cvss"`
	ConfidenceText string            `json:"confidence_text"`
	Description    string            `json:"description"`
	FirstFound     time.Time         `json:"first_found"`
	Count          string            `json:"count"`
	FalsePositive  string            `json:"false_positive"`
	Hash           string            `json:"hash"`
	CVE            string            `json:"cve"`
	Annotations    map[string]string `json:"annotations"`
	// The fields below are not used in this consumer. We use the text versions instead.
	// Severity   v1.Severity   `json:"severity"`
	// Confidence v1.Confidence `json:"confidence"`
}
