package main

import (
	"testing"

	"github.com/mitchellh/mapstructure"
	v1 "github.com/ocurity/dracon/api/proto/v1"
	"github.com/ocurity/dracon/components/producers"

	"github.com/stretchr/testify/assert"
)

func TestParseIssues(t *testing.T) {

	results, err := producers.ParseMultiJSONMessages([]byte(exampleOutput))
	assert.Nil(t, err)

	truffleResults := make([]TrufflehogOut, len(results))
	for i, r := range results {
		var x TrufflehogOut
		mapstructure.Decode(r, &x)
		truffleResults[i] = x
	}
	issues, err := parseIssues(truffleResults)
	assert.Nil(t, err)
	cs0 := "https://admin:admin@the-internet.herokuapp.com/basic_auth"
	CS1 := "wellnessbirdie-jaworskironni-bennettliz"
	expectedIssues := []v1.Issue{
		{
			Target:         "/code/keys",
			Type:           "trufflehog - filesystem",
			Title:          "PLAIN - URI",
			Severity:       v1.Severity_SEVERITY_UNSPECIFIED,
			Cvss:           0.0,
			Confidence:     v1.Confidence_CONFIDENCE_HIGH,
			Description:    "Raw:https://admin:admin@the-internet.herokuapp.com/basic_auth\nRedacted:https://*****:*****@the-internet.herokuapp.com/basic_auth\n",
			ContextSegment: &cs0,
		},
		{
			Target:         "https://github.com/foo/bar:c27298c30acdf69e611563b51ef2222f6324c916:src/org/zaproxy/zap/extension/directorylistv2_3/files/fuzzers/dirbuster/directory-list-2.3-big.txt:530222",
			Type:           "trufflehog - git",
			Title:          "BASE64 - Blogger",
			Severity:       v1.Severity_SEVERITY_UNSPECIFIED,
			Cvss:           0.0,
			Confidence:     v1.Confidence_CONFIDENCE_MEDIUM,
			Description:    "Raw:wellnessbirdie-jaworskironni-bennettliz\nRedacted:\nTimestamp:2015-04-13 16:07:20 +0000 UTC\nEmail:foo@example.com <foobar@users.noreply.github.com>\n",
			ContextSegment: &CS1,
		},
	}
	assert.Equal(t, expectedIssues[0], *issues[0])
	assert.Equal(t, expectedIssues[1], *issues[1])
}

const exampleOutput = `
{
    "SourceMetadata": {
        "Data": {
            "Filesystem": {
                "file": "/code/keys"
            }
        }
    },
    "SourceID": 15,
    "SourceType": 15,
    "SourceName": "trufflehog - filesystem",
    "DetectorType": 17,
    "DetectorName": "URI",
    "DecoderName": "PLAIN",
    "Verified": true,
    "Raw": "https://admin:admin@the-internet.herokuapp.com/basic_auth",
    "Redacted": "https://*****:*****@the-internet.herokuapp.com/basic_auth",
    "ExtraData": null,
    "StructuredData": null
}
{
    "SourceMetadata": {
      "Data": {
        "Git": {
          "commit": "c27298c30acdf69e611563b51ef2222f6324c916",
          "file": "src/org/zaproxy/zap/extension/directorylistv2_3/files/fuzzers/dirbuster/directory-list-2.3-big.txt",
          "email": "foo@example.com <foobar@users.noreply.github.com>",
          "repository": "https://github.com/foo/bar",
          "timestamp": "2015-04-13 16:07:20 +0000 UTC",
          "line": 530222
        }
      }
    },
    "SourceID": 0,
    "SourceType": 16,
    "SourceName": "trufflehog - git",
    "DetectorType": 302,
    "DetectorName": "Blogger",
    "DecoderName": "BASE64",
    "Verified": false,
    "Raw": "wellnessbirdie-jaworskironni-bennettliz",
    "Redacted": "",
    "ExtraData": null,
    "StructuredData": null
}  
`
