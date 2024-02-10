package pipelines

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/ocurity/dracon/pkg/components"
	"github.com/ocurity/dracon/pkg/manifests"
	"github.com/ocurity/dracon/pkg/pipelines"
	"github.com/spf13/cobra"
	tektonV1Beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/kustomize/api/types"
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

func loadPathOrURIOfK8sObj(ctx context.Context, root, pathOrURI, targetFile string) (runtime.Object, *schema.GroupVersionKind, error) {
	loader, err := manifests.NewLoader(root, pathOrURI, targetFile)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: could not resolve path or URI: %w", pathOrURI, err)
	}

	manifestBytes, err := loader.Load(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: could not load path or URI: %w", loader.Path(), err)
	}

	obj, gKV, err := manifests.K8sObjDecoder.Decode(manifestBytes, nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: could not decode file into a K8s object:3 %w", loader.Path(), err)
	}
	return obj, gKV, err
}

func buildPipeline(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("you need to provide the path of exactly one kustomization file")
	}

	kustomizationPath := args[0]
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

	kustomizationDir := kustomizationPath
	if strings.HasSuffix(kustomizationPath, "kustomization.yaml") {
		kustomizationDir = path.Dir(kustomizationPath)
	}

	// load the base pipeline
	if len(kustomization.Resources) != 2 {
		return fmt.Errorf("you need to specify the base pipeline and task in the resources field of the kustomization")
	}

	var baseTaskPath string
	var obj runtime.Object
	var gKV *schema.GroupVersionKind
	if strings.HasSuffix(kustomization.Resources[0], "pipeline.yaml") {
		obj, gKV, err = loadPathOrURIOfK8sObj(cmd.Context(), kustomizationDir, kustomization.Resources[0], "pipeline.yaml")
		baseTaskPath = kustomization.Resources[1]
	} else {
		obj, gKV, err = loadPathOrURIOfK8sObj(cmd.Context(), kustomizationDir, kustomization.Resources[1], "pipeline.yaml")
		baseTaskPath = kustomization.Resources[0]
	}
	if err != nil {
		return err
	}

	basePipeline, isAPipeline := obj.(*tektonV1Beta1.Pipeline)
	if !isAPipeline {
		return fmt.Errorf("object loaded is not a pipeline: %v", gKV)
	}

	if obj, gKV, err = loadPathOrURIOfK8sObj(cmd.Context(), kustomizationDir, baseTaskPath, "task.yaml"); err != nil {
		return err
	}

	baseTask, isATask := obj.(*tektonV1Beta1.Task)
	if !isATask {
		return fmt.Errorf("object loaded is not a task: %v", gKV)
	}

	if len(kustomization.Components) == 0 {
		return fmt.Errorf("%s: no components are listed in the kustomization", kustomizationPath)
	}

	taskList := []*tektonV1Beta1.Task{baseTask}
	for _, pathOrURI := range kustomization.Components {
		if obj, gKV, err = loadPathOrURIOfK8sObj(cmd.Context(), kustomizationDir, pathOrURI, "task.yaml"); err != nil {
			return err
		}

		newTask, isATask := obj.(*tektonV1Beta1.Task)
		if !isATask {
			return fmt.Errorf("object loaded is not a task: %v", gKV)
		}

		if err = components.ValidateTask(newTask); err != nil {
			return fmt.Errorf("%s: invalid task found: %w", newTask.Name, err)
		}

		newTask.Namespace = kustomization.Namespace
		taskList = append(taskList, newTask)
	}

	k8sBackend, err := pipelines.NewTektonV1Beta1Backend(basePipeline, taskList, kustomization.NamePrefix, kustomization.NameSuffix)
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
