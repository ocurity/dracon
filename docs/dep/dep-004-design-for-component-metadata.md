# Programmatically define what a component is, what it does and how it is used

## Introduction

Currently a component is a *Tekton Task that advances data between stages*.
This is Tekton specific and most importantly limiting.
We already had  to write workarounds to make components work in multiple use
cases. e.g. Golang-Nancy contains a step that creates files for an incoming
dependency so that the nancy tool can read them.

```yaml

  - name: go-deps
    image: "$(params.producer-golang-nancy-goImage)"
    script: | # Horrible bash
    volumeMounts:
    - mountPath: /scratch
      name: scratch
  - name: run-nancy
    imagePullPolicy: IfNotPresent
    image: docker.io/sonatypecommunity/nancy:v1.0.42-alpine
    script: | # more bash
    volumeMounts:
    - mountPath: /scratch
      name: scratch
  - name: produce-issues
    imagePullPolicy: IfNotPresent
    image: '{{ default "ghcr.io/ocurity/dracon" .Values.container_registry }}/components/producers/golang-nancy:{{ default "latest" .Values.dracon_os_component_version }}'
    command: ["/app/components/producers/golang-nancy/golang-nancy-parser"]
    args:
    - "-in=/scratch/out.json"
    - "-out=$(workspaces.output.path)/.dracon/producers/golang-nancy.pb"
    volumeMounts:
    - mountPath: /scratch
      name: scratch
```

This model has the following disadvantages.

* Bash is not testable easily
* Needing extra bash to run a tool is non-ideal and inflexible
* Preprocessing should be done in a wrapper

### Suggested solution

What a component is should be refactored to mean
*a pod that produces either json or protobuf output*.
Components are then categorized by the data they produce.
Components communicate depending on their type:

* **producers** output `LaunchToolResponse`s
  through the `putil/load.go` and `putil/write.go` SDK
* **enrichers or consumers** communicate by calling the same SDK
* Enrichers receive dracon `LaunchToolResponse`s and output
  `EnrichedLaunchToolResponse`s
* Consumers receive `EnrichedLaunchToolResponse`s and have no output.
* For **sources** their inputs are a single URL or IP, their outputs are:
  a status alert and as a byproduct a filesystem object or large text

## Design

### Definition of a source

Sources look like the following:

```yaml
---
metadata:
  name: <name of source>
  labels:
      v1.dracon.ocurity.com/component: source
spec:
  description: <>
  params:
   - name: input
     type: string
     default: ""
   - name: some other input
     type: string
  steps:
    - name: ""
      image: ""
      command: []
      args: []
```

They have a single step, that runs 1 command that does something to the
parameter named "input", e.g. if it's a purl, perhaps it resolves it and its
subdependencies. Or if the url is a git url, it clones the repo, if its an ARN
or an IP or some other pointer to infrastructure, perhaps it checks if the
infrastructure is accessible.

### Definition of a producer

A producer is the only component that sometimes needs to have more steps since
it interacts with outside binaries. Moreover, producers can either orchestrate
or just download results. It is suggested that producers are then defined as the
following(with the relevant Tekton-specific additions to make it work):

```yaml
   name: <string>
   metadata:
      labels:
         v1.dracon.ocurity.com/component: producer
   spec:
      description: <string>
      parameters:
         - name: <string>
            type: <a supported type>
            description: <string>
            default: <default type>
      steps:
         - name: fetch-results
         container: <string>
         command: <string>
         args: ["--fetch","<fetch.parameter.default.false>","--from","<some url>","<secrets etc>"]
         - name: orchestrate
         container: <string>
         command: <string>
         args: ["--orchestrate","<orchestrate.parameter.default.false>","--target","<url>","secrets etc"]
         - name: parse
         container: <string>
         command: <string>
         args: <array>
```

Each step is a binary, no more bash.
If we are not fetching results, the fetch parameter is set to false and the step
exits. Similarly if we are not orchestrating.

Since not every producer supports both orchestration and fetching this is up to
the producer.

#### Tool parameters

Tool parameters should be supported as much as possible with component
parameters, an easy way to do this is to provide an array.
Some parameters are necessary for dracon operation
(e.g. output and where the output goes), these should be hardcoded.

### Definition of an Enricher

Enrichers look like the following(Removing the tekton and workspace specific
defs for simplicity):

```yaml
---
---
metadata:
  name: <name>
  labels:
    v1.dracon.ocurity.com/component: enricher
spec:
  params: []
  steps:
  - name: run-enricher
```

They have a single step, that runs 1 command that does enriches the data
provided at runtime. They write to a location (they have no control
over where they read from or where they write to), output helpful messages to
stdout and exit.

### LaunchToolRequest

This is a pipelinerun.yaml, describes the mapping of parameters and values.

### Draconctl

Draconctl at some point in the future should test if a component follows the
definition before it installs it.
