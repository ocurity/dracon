package main

import (
	"bytes"
	"log"
	"os"

	wordwrap "github.com/mitchellh/go-wordwrap"
	"github.com/ocurity/dracon/build/tools/kustomize-component-generator/patches"
	kustomize "github.com/ocurity/dracon/build/tools/kustomize-component-generator/types/kustomize.config.k8s.io/v1alpha1"
	tekton "github.com/ocurity/dracon/build/tools/kustomize-component-generator/types/tekton.dev/v1beta1"
	"gopkg.in/yaml.v3"
)

func main() {
	taskYaml := os.Getenv("SRCS_TASK")
	outFile := os.Getenv("OUTS")
	taskBytes, err := os.ReadFile(taskYaml)
	if err != nil {
		panic(err)
	}

	task := &tekton.Task{}
	// task := TektonTask{}
	if err := yaml.Unmarshal(taskBytes, task); err != nil {
		panic(err)
	}

	component, ok := task.Metadata.Labels["v1.dracon.ocurity.com/component"]
	if !ok {
		log.Fatalf("missing .metadata.labels[\"v1.dracon.ocurity.com/component\"]")
	}

	switch component {
	case "base":
	case "source":
	case "producer":
	case "consumer":
	case "enricher":
		break
	default:
		log.Fatalf("unrecognised component: %s", component)
	}

	kustomizeComponent := kustomize.NewComponent()
	kustomizeComponent.Resources = []string{"task.yaml"}

	patchFuncs := []func() *kustomize.TargetPatch{
		patches.NewAddTaskToPipeline(task).GeneratePatch,
		patches.NewAddAnchorsToTask(task).GeneratePatch,
		// patches.NewAddProducerDependencyOnBase(task).GeneratePatch,
		patches.NewAddProducerDependencyOnSource(task).GeneratePatch,
		patches.NewAddProducerAggregatorAnchor(task).GeneratePatch,
		patches.NewAddEnricherDependencyOnProducerAggregator(task).GeneratePatch,
		patches.NewAddEnricherAggregatorAnchor(task).GeneratePatch,
		patches.NewAddConsumerDependencyOnEnricherAggregator(task).GeneratePatch,
		patches.NewAddScanUUIDAndStartTimeToTask(task).GeneratePatch,
		patches.NewAddScanUUIDAndStartTimeToPipeline(task).GeneratePatch,
	}

	for _, pFunc := range patchFuncs {
		if p := pFunc(); p != nil {
			kustomizeComponent.Patches = append(kustomizeComponent.Patches, p)
		}
	}

	yamlBytes, err := yaml.Marshal(kustomizeComponent)
	if err != nil {
		panic(err)
	}

	// add Head Comment
	yamlNode := &yaml.Node{}
	if err := yaml.Unmarshal(yamlBytes, yamlNode); err != nil {
		panic(err)
	}

	yamlNode.HeadComment = wordwrap.WrapString(
		"DO NOT EDIT. Code generated by: github.com/ocurity/dracon//build/tools/kustomize-component-generator.\n",
		80,
	)

	buf := &bytes.Buffer{}
	yamlEncoder := yaml.NewEncoder(buf)
	yamlEncoder.SetIndent(2)
	if err := yamlEncoder.Encode(yamlNode); err != nil {
		panic(err)
	}

	if err := os.WriteFile(outFile, buf.Bytes(), 0o600); err != nil {
		panic(err)
	}

	_ = taskYaml
}
