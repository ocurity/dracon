package pipelines

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/ocurity/dracon/pkg/components"
	"github.com/ocurity/dracon/pkg/manifests"
	"github.com/ocurity/dracon/pkg/pipelines"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"sigs.k8s.io/kustomize/api/types"

	tektonV1Beta1 "github.com/ocurity/dracon/pkg/types/tekton.dev/v1beta1"
)

func newBuildSubCmd() *cobra.Command {
	buildSubCmd := &cobra.Command{
		Use: "build",
		Short: `Build a pipeline out of an arbitrary number of components. The command expects to a
	path to a kustomization file where the resources list the base Pipeline and base Task and the
	components listed will be applied to the base manifests to generate a pipeline. You can choose
	to output the Pipeline in different formats. For the time being we only support Tekton Pipelines.`,
		GroupID: "Pipelines",
		RunE:    buildPipeline,
	}

	buildSubCmd.Flags().StringP("out", "o", "stdout", "The file to output the generated manifests")

	return buildSubCmd
}

func loadPathOrURI(ctx context.Context, root, pathOrURI, targetFile string, manifestStruct interface{}) error {
	loader, err := manifests.NewLoader(root, pathOrURI, targetFile)
	if err != nil {
		return fmt.Errorf("%s: could not resolve path or URI: %w", pathOrURI, err)
	}

	manifestBytes, err := loader.Load(ctx)
	if err != nil {
		return fmt.Errorf("%s: could not load path or URI: %w", pathOrURI, err)
	}

	if err = yaml.Unmarshal(manifestBytes, manifestStruct); err != nil {
		return fmt.Errorf("%s: could not unmarshal path or URI: %w", pathOrURI, err)
	}

	return nil
}

func loadBasePipelineAndTask(ctx context.Context, basePipeline *tektonV1Beta1.Pipeline, baseTask *tektonV1Beta1.Task, kustomizationPath string, resources []string) error {
	err := loadPathOrURI(ctx, path.Dir(kustomizationPath), path.Dir(resources[0]), "pipeline.yaml", &basePipeline)
	if err == nil {
		return loadPathOrURI(ctx, path.Dir(kustomizationPath), path.Dir(resources[1]), "task.yaml", &baseTask)
	}

	err = loadPathOrURI(ctx, path.Dir(kustomizationPath), path.Dir(resources[0]), "task.yaml", &baseTask)
	if err != nil {
		return err
	}

	return loadPathOrURI(ctx, path.Dir(kustomizationPath), resources[1], "pipeline.yaml", &basePipeline)
}

func buildPipeline(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("you need to provide the path of exactly one kustomization file")
	}

	kustomizationPath := args[0]
	// kustomizationPath := positionalArgs[0]
	kustomizationLoader, err := manifests.NewLoader(".", kustomizationPath, "kustomization.yaml")
	if err != nil {
		return fmt.Errorf("%s: could not read contents of file: %w", kustomizationPath, err)
	}

	// Load Pipeline kustomization file
	fileContents, err := kustomizationLoader.Load(cmd.Context())
	if err != nil {
		return fmt.Errorf("%s: could not read contents of file: %w", kustomizationPath, err)
	}

	// Parse Pipeline kustomization
	kustomization := types.Kustomization{}
	if err = yaml.Unmarshal(fileContents, &kustomization); err != nil {
		return fmt.Errorf("%s: could not unmarshal YAML file: %w", kustomizationPath, err)
	}

	// load the base pipeline
	if len(kustomization.Resources) != 2 {
		return fmt.Errorf("you need to specify the base pipeline and task in the resources field of the kustomization")
	}

	basePipeline := tektonV1Beta1.Pipeline{}
	baseTask := tektonV1Beta1.Task{}
	if err = loadBasePipelineAndTask(cmd.Context(), &basePipeline, &baseTask, kustomizationPath, kustomization.Resources); err != nil {
		return err
	}

	if len(kustomization.Components) == 0 {
		return fmt.Errorf("%s: no components are listed in the kustomization", kustomizationPath)
	}

	taskList := []tektonV1Beta1.Task{baseTask}
	for _, pathOrURI := range kustomization.Components {
		// Load Task file
		newTask := tektonV1Beta1.Task{}
		err = loadPathOrURI(cmd.Context(), path.Dir(kustomizationPath), pathOrURI, "task.yaml", &newTask)
		if err != nil {
			return err
		}

		if err = components.ValidateTask(&newTask); err != nil {
			return fmt.Errorf("%s: invalid task found: %w", newTask.Name, err)
		}

		newTask.Namespace = kustomization.Namespace
		taskList = append(taskList, newTask)
	}

	k8sBackend, err := pipelines.NewTektonBackend(&basePipeline, taskList, kustomization.NamePrefix, kustomization.NameSuffix)
	if err != nil {
		return fmt.Errorf("could not initialise backend: %w", err)
	}

	pipeline, err := k8sBackend.Generate()
	if err != nil {
		return fmt.Errorf("could not initialise backend: %w", err)
	}

	manifestYAMLBytes, err := yaml.Marshal(pipeline)
	if err != nil {
		return fmt.Errorf("could not marshal pipeline manifest: %w", err)
	}

	manifestYAMLBytes = append(manifestYAMLBytes, []byte("\n---\n")...)

	for _, task := range taskList {
		taskYAMLBytes, err := yaml.Marshal(task)
		if err != nil {
			return fmt.Errorf("could not marshal task manifest: %w", err)
		}

		manifestYAMLBytes = append(manifestYAMLBytes, taskYAMLBytes...)
		manifestYAMLBytes = append(manifestYAMLBytes, []byte("---\n")...)
	}

	output, err := cmd.Flags().GetString("out")
	if err != nil {
		return fmt.Errorf("could not get flag for output file: %w", err)
	}

	if output == "stdout" {
		output = "/dev/stdout"
	}

	if err = os.WriteFile(output, manifestYAMLBytes, 0600); err != nil {
		return fmt.Errorf("%s: could not write manifests to: %w", output, err)
	}

	return nil
}
