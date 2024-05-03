package pipelines

import (
	"context"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/ocurity/dracon/pkg/components"
	"github.com/ocurity/dracon/pkg/k8s"
)

func TestStuff(t *testing.T) {
	kubeConfig := path.Join(os.Getenv("HOME"), ".kube/config")
	restCfg, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	require.NoError(t, err)

	client, err := k8s.NewTypedClientForConfig(restCfg, "bla")
	require.NoError(t, err)

	orch := NewOrchestrator(client, "dracon")
	orch.Prepare(context.Background(), []components.Component{
		{
			Name:              "git-clone",
			OrchestrationType: components.ExternalHelm,
			Repository:        "dracon-oss-components",
		},
		{
			Name:              "producer-aggregator",
			OrchestrationType: components.ExternalHelm,
			Repository:        "dracon-oss-components",
		},
		{
			Name:              "producer-bla",
			OrchestrationType: components.ExternalHelm,
			Repository:        "dracon-oss-components",
		},
	})
}
