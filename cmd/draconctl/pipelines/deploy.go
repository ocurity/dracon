package pipelines

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	kustomizetypes "sigs.k8s.io/kustomize/api/types"

	"github.com/ocurity/dracon/pkg/components"
	"github.com/ocurity/dracon/pkg/files"
	"github.com/ocurity/dracon/pkg/k8s"
	"github.com/ocurity/dracon/pkg/manifests"
	"github.com/ocurity/dracon/pkg/pipelines"
)

var deploySubCmd = &cobra.Command{
	Use: "deploy",
	Short: `Deploy a pipeline out of an arbitrary number of components. The command expects to a
path to a kustomization file where the resources list the base Pipeline and base Task and the
components listed will be applied to the base manifests to generate a pipeline. You can choose
to output the Pipeline in different formats. For the time being we only support Tekton Pipelines.`,
	GroupID: "Pipelines",
	RunE:    deployPipeline,
}

var deploySubCmdFlags struct {
	ssa            string
	dryRun         bool
	k8sConfigFlags *genericclioptions.ConfigFlags
}

func init() {
	deploySubCmd.Flags().BoolVar(&deploySubCmdFlags.dryRun, "dry-run", false, "If set, print the generated pipeline in the stdout")
	deploySubCmd.Flags().StringVar(&deploySubCmdFlags.ssa, "ssa", "draconctl", "Server-side apply ID")

	deploySubCmdFlags.k8sConfigFlags = genericclioptions.NewConfigFlags(false)
	deploySubCmdFlags.k8sConfigFlags.AddFlags(deploySubCmd.Flags())
}

func deployPipeline(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("you need to provide the path of exactly one kustomization file")
	}

	kustomizationPath := args[0]
	kustomizationLoader, err := files.NewLoader(".", kustomizationPath, "kustomization.yaml")
	if err != nil {
		return fmt.Errorf("%s: could not read contents of file: %w", kustomizationPath, err)
	}

	// Load Pipeline kustomization file
	fileContents, err := kustomizationLoader.Load(cmd.Context())
	if err != nil {
		return fmt.Errorf("%s: could not read contents of file: %w", kustomizationPath, err)
	}

	// Parse Pipeline kustomization
	kustomization := &kustomizetypes.Kustomization{}
	if err = kustomization.Unmarshal(fileContents); err != nil {
		return fmt.Errorf("%s: could not unmarshal YAML file: %w", kustomizationPath, err)
	}

	basePipeline := pipelines.BasePipeline.DeepCopy()
	pipelineComponents := []components.Component{}
	for _, reference := range kustomization.Components {
		pipelineComponent, err := components.FromReference(cmd.Context(), reference)
		if err != nil {
			return err
		}

		pipelineComponents = append(pipelineComponents, pipelineComponent)
	}

	restCfg, err := deploySubCmdFlags.k8sConfigFlags.ToRESTConfig()
	if err != nil {
		return fmt.Errorf("could not initialise K8s client config with: %w", err)
	}

	client, err := k8s.NewTypedClientForConfig(restCfg, deploySubCmdFlags.ssa)
	if err != nil {
		return err
	}

	if *deploySubCmdFlags.k8sConfigFlags.Namespace == "" {
		*deploySubCmdFlags.k8sConfigFlags.Namespace = "dracon"
	}

	deploymentOrchestrator := pipelines.NewOrchestrator(client, *deploySubCmdFlags.k8sConfigFlags.Namespace)
	if err = deploymentOrchestrator.Prepare(cmd.Context(), pipelineComponents); err != nil {
		return err
	}

	pipeline, err := deploymentOrchestrator.Deploy(cmd.Context(), basePipeline, pipelineComponents, kustomization.NameSuffix, deploySubCmdFlags.dryRun)
	if err != nil {
		return err
	}

	if deploySubCmdFlags.dryRun {
		return manifests.TektonV1Beta1ObjEncoder.Encode(pipeline, os.Stdout)
	}

	return nil
}
