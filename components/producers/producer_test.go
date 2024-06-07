package producers

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/package-url/packageurl-go"

	v1 "github.com/ocurity/dracon/api/proto/v1"
	"github.com/ocurity/dracon/components"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

type testJ struct {
	Foo string
}

func TestWriteDraconOut(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "dracon-test")
	require.NoError(t, err)
	defer require.NoError(t, os.Remove(tmpFile.Name()))

	baseTime := time.Now().UTC()
	timestamp := baseTime.Format(time.RFC3339)
	require.NoError(t, os.Setenv(components.EnvDraconStartTime, timestamp))
	require.NoError(t, os.Setenv(components.EnvDraconScanID, "ab3d3290-cd9f-482c-97dc-ec48bdfcc4de"))

	OutFile = tmpFile.Name()
	Append = false

	err = WriteDraconOut(
		"dracon-test",
		[]*v1.Issue{
			{
				Target:      "/workspace/output/foobar",
				Title:       "/workspace/output/barfoo",
				Description: "/workspace/output/example.yaml",
				Cve:         "123-321",
			},
		},
	)
	require.NoError(t, err)

	pBytes, err := os.ReadFile(tmpFile.Name())
	require.NoError(t, err)

	res := v1.LaunchToolResponse{}
	require.NoError(t, proto.Unmarshal(pBytes, &res))

	assert.Equal(t, "dracon-test", res.GetToolName())
	assert.Equal(t, "./foobar", res.GetIssues()[0].GetTarget())
	assert.Equal(t, "./barfoo", res.GetIssues()[0].GetTitle())
	assert.Equal(t, "./example.yaml", res.GetIssues()[0].GetDescription())
	assert.Equal(t, baseTime.Unix(), res.GetScanInfo().GetScanStartTime().GetSeconds())
	assert.Equal(t, "ab3d3290-cd9f-482c-97dc-ec48bdfcc4de", res.GetScanInfo().GetScanUuid())
	assert.Equal(t, "123-321", res.GetIssues()[0].GetCve())
}

func TestWriteDraconOutAppend(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "dracon-test")
	require.NoError(t, err)
	defer require.NoError(t, os.Remove(tmpFile.Name()))

	baseTime := time.Now().UTC()
	timestamp := baseTime.Format(time.RFC3339)
	require.NoError(t, os.Setenv(components.EnvDraconStartTime, timestamp))
	require.NoError(t, os.Setenv(components.EnvDraconScanID, "ab3d3290-cd9f-482c-97dc-ec48bdfcc4de"))

	OutFile = tmpFile.Name()
	Append = true

	for _, i := range []int{0, 1, 2} {
		err = WriteDraconOut(
			"dracon-test",
			[]*v1.Issue{
				{
					Target:      fmt.Sprintf("target%d", i),
					Title:       fmt.Sprintf("title%d", i),
					Description: fmt.Sprintf("desc%d", i),
					Cve:         fmt.Sprintf("cve%d", i),
				},
			},
		)
		require.NoError(t, err)
	}

	pBytes, err := os.ReadFile(tmpFile.Name())
	require.NoError(t, err)

	res := v1.LaunchToolResponse{}
	require.NoError(t, proto.Unmarshal(pBytes, &res))

	assert.Equal(t, "dracon-test", res.GetToolName())
	assert.Equal(t, baseTime.Unix(), res.GetScanInfo().GetScanStartTime().GetSeconds())
	assert.Equal(t, "ab3d3290-cd9f-482c-97dc-ec48bdfcc4de", res.GetScanInfo().GetScanUuid())
	assert.Equal(t, 3, len(res.GetIssues()))

	for _, i := range []int{0, 1, 2} {
		assert.Equal(t, fmt.Sprintf("target%d", i), res.GetIssues()[i].GetTarget())
		assert.Equal(t, fmt.Sprintf("title%d", i), res.GetIssues()[i].GetTitle())
		assert.Equal(t, fmt.Sprintf("desc%d", i), res.GetIssues()[i].GetDescription())
		assert.Equal(t, fmt.Sprintf("cve%d", i), res.GetIssues()[i].GetCve())
	}
}

func TestParseMultiJSONMessages(t *testing.T) {
	testJSON := `{"Foo":"bar"}{"Foo":"barbar"}{"Foo":"barbarbar"}`

	inJSON, err := ParseMultiJSONMessages([]byte(testJSON))
	require.NoError(t, err)
	want := make([]testJ, len(inJSON))

	for i, v := range inJSON {
		var x testJ
		require.NoError(t, mapstructure.Decode(v, &x))
		want[i] = x
	}
	assert.Equal(t, want[0].Foo, "bar")
}

func TestGetPURLTarget(t *testing.T) {
	target := GetPURLTarget("deb", "debian", "curl", "7.68.0", nil, "")
	require.Equal(t, "pkg:deb/debian/curl@7.68.0", target)

	target = GetPURLTarget("bitbucket", "birkenfeld", "pygments-main", "244fd47e07d1014f0aed9c", nil, "")
	require.Equal(t, "pkg:bitbucket/birkenfeld/pygments-main@244fd47e07d1014f0aed9c", target)

	target = GetPURLTarget("docker", "customer", "dockerimage", "sha256:244fd47e07d1004f0aed9c", packageurl.Qualifiers{
		{Key: "repository_url", Value: "gcr.io"},
	}, "")
	require.Equal(t, "pkg:docker/customer/dockerimage@sha256:244fd47e07d1004f0aed9c?repository_url=gcr.io", target)

	target = GetPURLTarget("npm", "", "foobar", "12.3.1", nil, "")
	require.Equal(t, "pkg:npm/foobar@12.3.1", target)

	target = GetPURLTarget("pypi", "", "django", "1.11.1", nil, "")
	require.Equal(t, "pkg:pypi/django@1.11.1", target)

	target = GetPURLTarget("deb", "debian", "curl", "7.50.3-1", packageurl.Qualifiers{
		{Key: "arch", Value: "i386"},
		{Key: "distro", Value: "jessie"},
	}, "")

	require.Equal(t, "pkg:deb/debian/curl@7.50.3-1?arch=i386&distro=jessie", target)
}
