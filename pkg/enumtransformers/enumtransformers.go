// Package enumtransformers transforms from dracon internal enums to text and back
package enumtransformers

import (
	v1 "github.com/ocurity/dracon/api/proto/v1"
)

// SeverityToText transforms dracon severity into a nicer textual format for use with third party systems
func SeverityToText(severity v1.Severity) string {
	switch severity {
	case v1.Severity_SEVERITY_INFO:
		return "Info"
	case v1.Severity_SEVERITY_LOW:
		return "Low"
	case v1.Severity_SEVERITY_MEDIUM:
		return "Medium"
	case v1.Severity_SEVERITY_HIGH:
		return "High"
	case v1.Severity_SEVERITY_CRITICAL:
		return "Critical"
	default:
		return "N/A"
	}
}

// SeverityToText transforms dracon severity into a nicer textual format for use with third party systems
func TextToSeverity(severity string) v1.Severity {
	switch severity {
	case "Info":
		return v1.Severity_SEVERITY_INFO
	case "Low":
		return v1.Severity_SEVERITY_LOW 
	case "Medium":
		return v1.Severity_SEVERITY_MEDIUM
	case "High":
		return v1.Severity_SEVERITY_HIGH
	case "Critical":
		return v1.Severity_SEVERITY_CRITICAL
	default:
		return v1.Severity_SEVERITY_UNSPECIFIED
	}
}

// ConfidenceToText transforms dracon confidence into a nicer textual format for use with third party systems
func ConfidenceToText(confidence v1.Confidence) string {
	switch confidence {
	case v1.Confidence_CONFIDENCE_INFO:
		return "Info"
	case v1.Confidence_CONFIDENCE_LOW:
		return "Low"
	case v1.Confidence_CONFIDENCE_MEDIUM:
		return "Medium"
	case v1.Confidence_CONFIDENCE_HIGH:
		return "High"
	case v1.Confidence_CONFIDENCE_CRITICAL:
		return "Critical"
	default:
		return "N/A"
	}
}

// ConfidenceToText transforms dracon confidence into a nicer textual format for use with third party systems
func TextToConfidence(confidence string) v1.Confidence {
	switch confidence {
	case "Info":
		return v1.Confidence_CONFIDENCE_INFO
	case "Low":
		return v1.Confidence_CONFIDENCE_LOW
	case "Medium":
		return v1.Confidence_CONFIDENCE_MEDIUM
	case "High":
		return v1.Confidence_CONFIDENCE_HIGH
	case "Critical":
		return v1.Confidence_CONFIDENCE_CRITICAL
	default:
		return v1.Confidence_CONFIDENCE_UNSPECIFIED
	}
}
