package npmquickaudit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	gock "gopkg.in/h2non/gock.v1"
)

var advisoryAST = &AdvisoryData{
	CVEs:               []string{"CVE-2020-15168"},
	CWE:                "CWE-400",
	Overview:           "Node Fetch did not honor the size option after following a redirect, which means that when a content size was over the limit, a FetchError would never get thrown and the process would end without failure.\n\nFor most people, this fix will have a little or no impact. However, if you are relying on node-fetch to gate files above a size, the impact could be significant, for example: If you don't double-check the size of the data after fetch() has completed, your JS thread could get tied up doing work on a large file (DoS) and/or cost you money in computing.",
	PatchedVersions:    ">=2.6.1 <3.0.0-beta.1|| >= 3.0.0-beta.9",
	Recommendation:     "Upgrade to version 2.6.1 or 3.0.0-beta.9",
	References:         "",
	VulnerableVersions: "< 2.6.1 || >= 3.0.0-beta.1 < 3.0.0-beta.9",
}

func TestNewAdvisoryDataNotFound(t *testing.T) {
	defer gock.Off()
	gock.New("https://npmjs.com").
		Get("/advisories/404").
		MatchHeader("X-Spiferack", "1").
		Reply(404)

	advisory, err := NewAdvisoryData("https://npmjs.com/advisories/404")
	assert.Nil(t, advisory)
	assert.Error(t, err)

	assert.True(t, gock.IsDone())
}

func TestNewAdvisoryDataNotJSON(t *testing.T) {
	defer gock.Off()
	gock.New("https://npmjs.com").
		Get("/advisories/666").
		MatchHeader("X-Spiferack", "1").
		Reply(200).
		AddHeader("Content-Type", "application/json").
		File("components/producers/typescript-npm-audit/types/npmquickaudit/npm_advisory_not_json")

	advisory, err := NewAdvisoryData("https://npmjs.com/advisories/666")
	assert.Nil(t, advisory)
	assert.Errorf(t, err, "npm Registry did not respond with JSON content")

	assert.True(t, gock.IsDone())
}

func TestNewAdvisoryDataNoAdvisoryData(t *testing.T) {
	defer gock.Off()
	gock.New("https://npmjs.com").
		Get("/advisories/999").
		MatchHeader("X-Spiferack", "1").
		Reply(200).
		AddHeader("Content-Type", "application/json").
		File("components/producers/typescript-npm-audit/types/npmquickaudit/npm_advisory_no_advisorydata")

	advisory, err := NewAdvisoryData("https://npmjs.com/advisories/999")
	assert.Nil(t, advisory)
	assert.Errorf(t, err, "npm Registry response did not contain an advisoryData key")

	assert.True(t, gock.IsDone())
}

func TestNewAdvisoryDataValid(t *testing.T) {
	defer gock.Off()
	gock.New("https://npmjs.com").
		Get("/advisories/1556").
		MatchHeader("X-Spiferack", "1").
		Reply(200).
		AddHeader("Content-Type", "application/json").
		File("components/producers/typescript-npm-audit/types/npmquickaudit/npm_advisory_1556")

	advisory, err := NewAdvisoryData("https://npmjs.com/advisories/1556")
	assert.NoError(t, err)
	assert.True(t, assert.ObjectsAreEqual(advisoryAST, advisory))

	assert.True(t, gock.IsDone())
}
