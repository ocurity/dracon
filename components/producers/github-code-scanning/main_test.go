package main

import (
	"testing"

	"github.com/google/go-github/v63/github"
	"github.com/stretchr/testify/require"

	v1 "github.com/ocurity/dracon/api/proto/v1"
)

func TestParseIssues(t *testing.T) {
	alerts := []*github.Alert{
		{
			Number: github.Int(1),
			Rule: &github.Rule{
				Tags:        []string{"security", "external/cwe/cwe-022"},
				Severity:    github.String("low"),
				Description: github.String("Test description"),
			},
			HTMLURL: github.String("https://example.com"),
			MostRecentInstance: &github.MostRecentInstance{
				Location: &github.Location{
					Path:      github.String("spec-main/api-session-spec.ts"),
					StartLine: github.Int(917),
					EndLine:   github.Int(918),
				},
				Message: &github.Message{
					Text: github.String("Test message"),
				},
			},
		},
	}

	issues := parseIssues(alerts)

	expected := []*v1.Issue{
		{
			Target:      "file://spec-main/api-session-spec.ts:917-918",
			Type:        "1",
			Title:       "Test description",
			Severity:    v1.Severity_SEVERITY_LOW,
			Cvss:        0,
			Confidence:  v1.Confidence_CONFIDENCE_UNSPECIFIED,
			Description: "Test message",
			Source:      "https://example.com",
			Cwe:         []int32{22},
		},
	}

	require.Equal(t, expected, issues)
}

func TestParseGitHubSeverity(t *testing.T) {
	testCases := []struct {
		name     string
		severity string
		expected v1.Severity
	}{
		{
			name:     "low severity",
			severity: "low",
			expected: v1.Severity_SEVERITY_LOW,
		},
		{
			name:     "medium severity",
			severity: "medium",
			expected: v1.Severity_SEVERITY_MEDIUM,
		},
		{
			name:     "high severity",
			severity: "high",
			expected: v1.Severity_SEVERITY_HIGH,
		},
		{
			name:     "critical severity",
			severity: "critical",
			expected: v1.Severity_SEVERITY_CRITICAL,
		},
		{
			name:     "unspecified severity",
			severity: "unknown",
			expected: v1.Severity_SEVERITY_UNSPECIFIED,
		},
		{
			name:     "empty severity",
			severity: "",
			expected: v1.Severity_SEVERITY_UNSPECIFIED,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			severity := parseGitHubSeverity(tc.severity)
			require.Equal(t, tc.expected, severity)
		})
	}

}

func TestParseGitHubCWEsFromTags(t *testing.T) {
	testCases := []struct {
		name     string
		tags     []string
		expected []int32
	}{
		{
			name:     "single CWE",
			tags:     []string{"security", "external/cwe/cwe-022"},
			expected: []int32{22},
		},
		{
			name:     "multiple CWEs",
			tags:     []string{"security", "external/cwe/cwe-022", "external/cwe/cwe-023", "external/cwe/cwe-124"},
			expected: []int32{22, 23, 124},
		},
		{
			name:     "no CWEs",
			tags:     []string{"security"},
			expected: []int32{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cwes := parseGithubCWEsFromTags(tc.tags)
			require.Equal(t, tc.expected, cwes)
		})
	}
}
