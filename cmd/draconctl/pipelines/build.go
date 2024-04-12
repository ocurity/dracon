package pipelines

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
	kustomizetypes "sigs.k8s.io/kustomize/api/types"

	"github.com/ocurity/dracon/pkg/files"
	"github.com/ocurity/dracon/pkg/manifests"
	"github.com/ocurity/dracon/pkg/pipelines"
)

var buildSubCmd = &cobra.Command{
	Use: "build",
	Short: `Build a pipeline out of an arbitrary number of components. The command expects to a
path to a kustomization file where the resources list the base Pipeline and base Task and the
components listed will be applied to the base manifests to generate a pipeline. You can choose
to output the Pipeline in different formats. For the time being we only support Tekton Pipelines.`,
	GroupID: "Pipelines",
	RunE:    buildPipeline,
}

func init() {
	buildSubCmd.Flags().StringP("out", "o", "stdout", "The file to output the generated manifests")
}

func buildPipeline(cmd *cobra.Command, args []string) error {
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

	kustomizationDir := kustomizationPath
	if strings.HasSuffix(kustomizationPath, "kustomization.yaml") {
		kustomizationDir = path.Dir(kustomizationPath)
	}

	// load the base pipeline
	if len(kustomization.Resources) != 2 {
		return fmt.Errorf("you need to specify the base pipeline and task in the resources field of the kustomization")
	}

	pipelineKustomization := &pipelines.Kustomization{Kustomization: kustomization, KustomizationDir: kustomizationDir}
	basePipeline, taskList, err := pipelineKustomization.ResolveKustomizationResources(cmd.Context())
	if err != nil {
		return fmt.Errorf("%s: could not resolve base pipeline and tasks: %w", kustomizationDir, err)
	}

	k8sBackend, err := pipelines.NewTektonV1Beta1Backend(basePipeline, taskList, kustomization.NameSuffix)
	if err != nil {
		return fmt.Errorf("could not initialise backend: %w", err)
	}

	pipeline, err := k8sBackend.Generate()
	if err != nil {
		return fmt.Errorf("could not initialise backend: %w", err)
	}

	output, err := cmd.Flags().GetString("out")
	if err != nil {
		return fmt.Errorf("could not get flag for output file: %w", err)
	}

	if output == "stdout" {
		output = "/dev/stdout"
	}

	out, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0660)
	if err != nil {
		return fmt.Errorf("%s: could not open file for writing manifests to: %w", output, err)
	}

	for _, parameter := range pipeline.Spec.Params {
		if parameter.Type == "" {
			return fmt.Errorf("pipeline parameter %s has no type set", parameter.Name)
		} else if parameter.Default != nil && parameter.Default.Type == "" {
			return fmt.Errorf("pipeline default value for parameter %s has no type set", parameter.Name)
		}
	}
	for _, task := range pipeline.Spec.Tasks {
		for _, param := range task.Params {
			if param.Value.Type == "" {
				return fmt.Errorf("pipeline's task %s param %s has no value type set", task.Name, param.Name)
			}
		}
	}

	if err = manifests.TektonV1Beta1ObjEncoder.Encode(pipeline, out); err != nil {
		return fmt.Errorf("could not marshal pipeline manifest: %w", err)
	}

	for _, task := range taskList {
		if _, err = fmt.Fprint(out, "---\n"); err != nil {
			return fmt.Errorf("could not write to output: %w", err)
		}
		if err = manifests.TektonV1Beta1ObjEncoder.Encode(task, out); err != nil {
			return fmt.Errorf("could not marshal task manifest: %w", err)
		}
	}

	return nil
}
