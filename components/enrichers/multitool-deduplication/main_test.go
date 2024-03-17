package main

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	v1 "github.com/ocurity/dracon/api/proto/v1"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	license = "Foo License v0"
)

func genScanCVE(toolName string, targets []string, cves []string) *v1.LaunchToolResponse { // gen SCA scan
	id := uuid.New()
	startTime := timestamppb.Now()
	scan := v1.LaunchToolResponse{
		ScanInfo: &v1.ScanInfo{ScanUuid: id.String(), ScanStartTime: startTime},
		ToolName: toolName,
		Issues:   []*v1.Issue{},
	}
	for i, t := range targets {
		scan.Issues = append(scan.Issues, &v1.Issue{
			Target:   t,
			Type:     fmt.Sprintf("someID %d", i),
			Title:    fmt.Sprintf("someTitle %d", i),
			Severity: v1.Severity_SEVERITY_UNSPECIFIED,
			Cvss:     float64(i),
			Cve:      cves[i%len(cves)],
		})
	}
	return &scan
}

func genScanCWE(toolName string, targets []string, cwes []int32) *v1.LaunchToolResponse { // gen SAST Scan
	id := uuid.New()
	startTime := timestamppb.Now()
	scan := v1.LaunchToolResponse{
		ScanInfo: &v1.ScanInfo{ScanUuid: id.String(), ScanStartTime: startTime},
		ToolName: toolName,
		Issues:   []*v1.Issue{},
	}
	for i, t := range targets {
		scan.Issues = append(scan.Issues, &v1.Issue{
			Target:   t,
			Type:     fmt.Sprintf("someID %d", i),
			Title:    fmt.Sprintf("someTitle %d", i),
			Severity: v1.Severity_SEVERITY_UNSPECIFIED,
			Cvss:     float64(i),
			Cwe:      cwes,
		})
	}
	return &scan
}

func TestHandlePurl(t *testing.T) {
	cve := []string{}
	targets := []string{}
	for i := 0; i <= 10; i++ {
		cve = append(cve, fmt.Sprintf("CVE-2009-%d", i+500))
		targets = append(targets, fmt.Sprintf("/foo/bar/target%d:%d", i%3, i%3))
	}
	scan1 := genScanCVE("tool1", targets, cve)
	scan2 := genScanCVE("tool2", targets, cve)
	for _, issue := range scan1.GetIssues() {
		assert.Equal(t, "", handlePurl(issue, scan1.ToolName))
	}
	for _, issue := range scan1.GetIssues() {
		assert.Equal(t, "tool1", handlePurl(issue, scan2.ToolName))
	}
}

// func TestConvertTargetToLineRange(t *testing.T) {
// 	type dat struct {
// 		input  string
// 		output lineRange
// 	}
// 	data := []dat{
// 		{
// 			input:  "/positive/max/example.blah:234-235",
// 			output: lineRange{filePath: "/positive/max/example.blah", lineStart: 234, lineEnd: 235},
// 		},
// 		// {
// 		// 	input:  "no/leading/slash.foo:123-124",
// 		// 	output: lineRange{filePath: "no/leading/slash.foo", lineStart: 123, lineEnd: 124},
// 		// },
// 		// {input: "single/line/in/range.go:595", output: lineRange{filePath: "single/line/in/range.go", lineStart: 595, lineEnd: 0}},
// 		// {input: "with/spaces.go : 595", output: lineRange{filePath: "with/spaces", lineStart: 595, lineEnd: 0}},
// 		// {input: "with/spaces/range.go : 595 - 596", output: lineRange{filePath: "with/spaces/range.go", lineStart: 595, lineEnd: 596}},
// 	}
// 	for _, d := range data {
// 		assert.Equal(t, d.output, convertTargetToLineRange(d.input))
// 	}
// }

// func TestHandleFilesystem(t *testing.T) {
// 	cwe := []int32{}
// 	targets := []string{}
// 	for i := 0; i <= 10; i++ {
// 		cwe = append(cwe, int32(i+500))
// 		targets = append(targets, fmt.Sprintf("/foo/bar/target%d:%d", i%3, i%3))
// 	}
// 	scan1 := genScanCWE("tool1", targets, cwe)
// 	scan2 := genScanCWE("tool2", targets, cwe)
// 	for _, issue := range scan1.GetIssues() {
// 		assert.Equal(t, "", handleFilesystem(issue, scan1.ToolName))
// 	}
// 	for _, issue := range scan1.GetIssues() {
// 		assert.Equal(t, "tool2", handleFilesystem(issue, scan2.ToolName))
// 	}
// }
